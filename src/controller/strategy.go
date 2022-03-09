package main

import (
	"owl/common/types"
	"sync"
	"time"
)

//持续定时加载报警策略
func (c *Controller) loadStrategiesForever() {
	var wg sync.WaitGroup
	for {
		// 获取所有产品线
		products := mydb.GetProducts()
		// 更新产品线告警队列
		c.refreshQueue(products)
		for {
			if len(c.nodePool.Nodes) != 0 {
				break
			}
			lg.Warn("no inspector connected, do not generate task, retry after 1 seconds")
			time.Sleep(time.Second)
		}

		for _, product := range products {
			// 根据产品线 id 获取策略
			strategies := mydb.GetStrategies(product.ID)
			wg.Add(1)
			go func(strategies []*types.Strategy) {
				defer wg.Done()
				for _, strategy := range strategies {
					if strategy.Enable == false {
						lg.Info("strategy %s is not enabled, skipped it.", strategy.Name)
						continue
					}
					// 根据策略 id 获取 trigger
					triggers := mydb.GetTriggersByStrategyID(strategy.ID)
					//如果没有 trigger 则忽略
					if len(triggers) == 0 {
						lg.Warn("strategy %s has no trigger, skipped it.", strategy.Name)
						continue
					}
					// 生成 AlarmTask
					c.processSingleStrategy(strategy, triggers)
				}
			}(strategies)
		}
		wg.Wait()
		time.Sleep(time.Second * time.Duration(GlobalConfig.LOAD_STRATEGIES_INTERVAL))
	}
}

func (c *Controller) processSingleStrategy(strategy *types.Strategy, triggers map[string]*types.Trigger) {
	globalHosts := make([]*types.Host, 0)
	for _, group := range mydb.GetGroupsByStrategyID(strategy.ID) {
		groupHosts := mydb.GetHostsByGroupID(group.ID)
		globalHosts = append(globalHosts, groupHosts...)
	}
	exHosts := mydb.GetHostsExByStrategyID(strategy.ID)
	for _, host := range globalHosts {
		// 过滤静音主机
		if host.IsMute() {
			lg.Info("strategy %d:%v host is mute %v:%v:%v, mute_time:%s",
				strategy.ID, strategy.Name, host.ID, host.IP, host.Hostname, host.MuteTime)
			continue
		}
		// 过滤排除主机
		if _, ok := exHosts[host.ID]; ok {
			lg.Info("strategy %d:%v exclude host %v:%v:%v",
				strategy.ID, strategy.Name, host.ID, host.IP, host.Hostname)
			continue
		}
		// 向 taskCache 添加任务
		task := types.NewAlarmTask(host, strategy, triggers)
		if err := c.taskPool.putTask(task); err != nil {
			lg.Error("put new task into task pool failed %v, maybe you need to increase the task_pool_size", err)
			continue
		}
		c.taskCache.Set(task.ID, task, 10*time.Minute)
	}
}

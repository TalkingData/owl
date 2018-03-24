package main

import (
	"owl/common/types"
	"sync"
	"time"
)

func (c *Controller) Add(task *types.AlarmTask) {
	if len(task.Triggers) == 0 {
		lg.Warn("task %v has no triggers, skipped it.", task.ID)
		return
	}
	if task.Host.Status == "2" {
		lg.Debug("host %v is forbidden, skipped it.", task.Host.ID)
		return
	}
	c.taskCache.Set(task.ID, task, 10*time.Minute)
}

//持续定时加载报警策略
func (c *Controller) loadStrategiesForever() {
	for {
		var wait_group sync.WaitGroup
		products := mydb.GetProducts()
		c.refreshQueue(products)
		for _, product := range products {
			strategies := mydb.GetStrategies(product.ID)
			wait_group.Add(1)
			go func(strategies []*types.Strategy) {
				defer wait_group.Done()
				for _, strategy := range strategies {
					if strategy.Enable == false {
						lg.Info("strategy %v is not enabled, skipped it.", strategy.Name)
						continue
					}
					triggers := mydb.GetTriggersByStrategyID(strategy.ID)
					c.processSingleStrategy(strategy, triggers)
				}
			}(strategies)
		}
		wait_group.Wait()
		c.taskPool.PutTasks(c.taskCache.GetItems())
		lg.Info("loaded tasks %v for all products", c.taskCache.ItemCount())
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
		if _, ok := exHosts[host.ID]; !ok {
			c.Add(types.NewAlarmTask(host, strategy, triggers))
		}
	}
}

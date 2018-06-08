package main

import (
	"fmt"
	"owl/common/types"
	"time"
)

const (
	NEW_ALARM = iota + 1
	OLD_ALARM
	RESTORE_ALARM
)

type QueueEvent struct {
	status         int
	strategy_event *types.StrategyEvent
	trigger_events map[string]*types.TriggerEvent
}

//持续处理报警结果
func (c *Controller) processStrategyResultForever() {
	for i := 0; i < GlobalConfig.WORKER_COUNT; i++ {
		go func() {
			for {
				select {
				case result := <-c.resultPool.results:
					lg.Debug("get task result from result pool, taskid:%s", result.TaskID)
					c.processResult(result)
				default:
					time.Sleep(time.Millisecond * 100)
				}
			}
		}()
	}
}

//处理来自Inspector计算后的结果对象
func (c *Controller) processResult(strategyResult *types.StrategyResult) {
	item, ok := c.taskCache.Get(strategyResult.TaskID)
	if !ok {
		lg.Error(fmt.Sprintf("task %v not in cached task pool", strategyResult.TaskID))
		return
	}
	task, _ := item.(*types.AlarmTask)

	if strategyResult.ErrorMessage != "" {
		eventStatus := types.EVENT_UNKNOW
		if strategyResult.ErrorMessage == "no data" {
			eventStatus = types.EVENT_NODATA
		}
		lg.Warn(fmt.Sprintf("task %v %v strategy:%s hostname:%s ip:%s",
			strategyResult.TaskID, strategyResult.ErrorMessage, task.Strategy.Name, task.Host.Hostname, task.Host.IP))
		mydb.CreateStrategyEventFailed(task.Strategy.ID, task.Host.ID, eventStatus, strategyResult.ErrorMessage)
		return
	}

	go syncStrategyEventFailed(task.Strategy.ID, task.Host.ID)

	var awareStrategyEvent, oldStrategyEvent *types.StrategyEvent

	//获取知悉
	if awareStrategyEvent = mydb.GetStrategyEvent(task.Strategy.ID, types.EVENT_AWARE, task.Host.ID); awareStrategyEvent != nil {
		strategyEvent, triggerEvents := generateEvent(awareStrategyEvent, strategyResult, task)
		// 告警恢复
		if strategyResult.Triggered == false {
			triggeredTriggerEvents := mydb.GetTriggeredTriggerEvents(strategyEvent.ID)
			triggerEvents = tagChangedTrigger(triggeredTriggerEvents, triggerEvents)
			strategyEvent.Status = types.EVENT_CLOSED
			event := &QueueEvent{RESTORE_ALARM, strategyEvent, triggerEvents}
			c.eventQueuesMutex.RLock()
			if err := c.eventQueues[strategyEvent.ProductID].putQueueEvent(event); err != nil {
				lg.Error("put strategy event to event queue failed %s", err)
			} else {
				lg.Info("put strategy event to queue %d strategy:%s hostname:%s ip:%s",
					strategyEvent.ProductID, strategyEvent.StrategyName, strategyEvent.HostName, strategyEvent.IP)
			}
			c.eventQueuesMutex.RUnlock()
			lg.Debug(fmt.Sprintf("put restore event by strategy %s into event queue.", strategyEvent.StrategyName))
			return
		}
		// 知悉过期
		if time.Now().Sub(strategyEvent.AwareEndTime) > 0 && strategyEvent.AwareEndTime.Sub(types.DEFAULT_TIME) != 0 {
			strategyEvent.Status = types.EVENT_NEW
			strategyEvent.AwareEndTime = types.DEFAULT_TIME
			// 未达到最大告警次数
			if strategyEvent.Count < task.Strategy.AlarmCount || task.Strategy.AlarmCount == 0 {
				strategyEvent.Count++
				event := &QueueEvent{OLD_ALARM, strategyEvent, triggerEvents}
				c.eventQueuesMutex.RLock()
				if err := c.eventQueues[strategyEvent.ProductID].putQueueEvent(event); err != nil {
					lg.Error("put strategy event to event queue failed %s", err)
				} else {
					lg.Info("put strategy event to queue %d strategy:%s hostname:%s ip:%s",
						strategyEvent.ProductID, strategyEvent.StrategyName, strategyEvent.HostName, strategyEvent.IP)
				}
				c.eventQueuesMutex.RUnlock()
				return
			}
		}
		mydb.UpdateStrategyEvent(strategyEvent, triggerEvents)
		return
	}

	//获取已产生报警
	if oldStrategyEvent = mydb.GetStrategyEvent(task.Strategy.ID, types.EVENT_NEW, task.Host.ID); oldStrategyEvent != nil {
		strategyEvent, triggerEvents := generateEvent(oldStrategyEvent, strategyResult, task)
		//告警恢复
		if strategyResult.Triggered == false {
			triggeredTriggerEvents := mydb.GetTriggeredTriggerEvents(strategyEvent.ID)
			triggerEvents = tagChangedTrigger(triggeredTriggerEvents, triggerEvents)
			strategyEvent.Status = types.EVENT_CLOSED
			event := &QueueEvent{RESTORE_ALARM, strategyEvent, triggerEvents}
			c.eventQueuesMutex.RLock()
			if err := c.eventQueues[strategyEvent.ProductID].putQueueEvent(event); err != nil {
				lg.Error("put strategy event to event queue failed %s", err)
			} else {
				lg.Info("put strategy event to queue %d strategy:%s hostname:%s ip:%s",
					strategyEvent.ProductID, strategyEvent.StrategyName, strategyEvent.HostName, strategyEvent.IP)
			}
			c.eventQueuesMutex.RUnlock()
			return
		}
		if strategyEvent.Count >= task.Strategy.AlarmCount && task.Strategy.AlarmCount != 0 {
			mydb.UpdateStrategyEvent(strategyEvent, triggerEvents)
			return
		}
		strategyEvent.Count++
		event := &QueueEvent{OLD_ALARM, strategyEvent, triggerEvents}
		c.eventQueuesMutex.RLock()
		if err := c.eventQueues[strategyEvent.ProductID].putQueueEvent(event); err != nil {
			lg.Error("put strategy event to event queue failed %s", err)
		} else {
			lg.Info("put strategy event to queue %d strategy:%s hostname:%s ip:%s",
				strategyEvent.ProductID, strategyEvent.StrategyName, strategyEvent.HostName, strategyEvent.IP)
		}
		c.eventQueuesMutex.RUnlock()
		return
	}

	// 新产生的告警
	if oldStrategyEvent == nil && awareStrategyEvent == nil {
		strategyEvent, triggerEvents := generateEvent(nil, strategyResult, task)
		if strategyResult.Triggered == true {
			strategyEvent.Status = types.EVENT_NEW
			event := &QueueEvent{NEW_ALARM, strategyEvent, triggerEvents}
			c.eventQueuesMutex.RLock()
			if err := c.eventQueues[strategyEvent.ProductID].putQueueEvent(event); err != nil {
				lg.Error("put strategy event to event queue failed %s", err)
			} else {
				lg.Info("put strategy event to queue %d strategy:%s hostname:%s ip:%s",
					strategyEvent.ProductID, strategyEvent.StrategyName, strategyEvent.HostName, strategyEvent.IP)
			}
			c.eventQueuesMutex.RUnlock()
		}
	}
}

func tagChangedTrigger(triggeredTriggerEvents []*types.TriggerEvent, triggerEvents map[string]*types.TriggerEvent) map[string]*types.TriggerEvent {
	for _, e := range triggeredTriggerEvents {
		if value, ok := triggerEvents[e.Index+e.Metric+e.Tags]; ok {
			value.TriggerChanged = true
		}
	}
	return triggerEvents
}

func syncStrategyEventFailed(strategy_id int, host_id string) error {
	return mydb.DeleteStrategyFailed(strategy_id, host_id)
}

//根据计算后的结果生成新的报警事件,将策略更改后的信息同步到历史报警事件
func generateEvent(strategy_event *types.StrategyEvent, strategy_result *types.StrategyResult, task *types.AlarmTask) (merged_strategy_event *types.StrategyEvent, merged_trigger_events map[string]*types.TriggerEvent) {
	merged_strategy_event = types.NewStrategyEvent(
		task.Strategy.ProductID,
		task.Strategy.ID,
		task.Strategy.Name,
		task.Strategy.Priority,
		task.Strategy.Cycle,
		task.Strategy.AlarmCount,
		task.Strategy.Expression,
		strategy_result.CreateTime,
		task.Host.ID,
		task.Host.Hostname,
		task.Host.IP,
		strategy_result.ErrorMessage)

	if strategy_event != nil {
		merged_strategy_event.ID = strategy_event.ID
		merged_strategy_event.Count = strategy_event.Count
		merged_strategy_event.Status = strategy_event.Status
		merged_strategy_event.CreateTime = strategy_event.CreateTime
		merged_strategy_event.AwareEndTime = strategy_event.AwareEndTime
	}

	merged_trigger_events = make(map[string]*types.TriggerEvent)
	for index, trigger_result_set := range strategy_result.TriggerResultSets {
		trigger := task.Triggers[index]
		for _, trigger_result := range trigger_result_set.TriggerResults {
			merged_trigger_events[index+trigger.Metric+trigger_result.Tags] = types.NewTriggerEvent(merged_strategy_event.ID, index, trigger.Metric, trigger_result.Tags, trigger_result.AggregateTags, trigger.Symbol, trigger.Method, trigger.Number, trigger.Threshold, trigger_result.CurrentThreshold, trigger_result.Triggered)
		}
	}
	return
}

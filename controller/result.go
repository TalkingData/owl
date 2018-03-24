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
	worker_count := GlobalConfig.WORKER_COUNT
	for worker_count > 0 {
		go func() {
			for {
				select {
				case result := <-c.resultPool.results:
					lg.Debug("get result %v from result pool", result.TaskID)
					c.processResult(result)
				default:
					time.Sleep(time.Millisecond * 100)
				}
			}
		}()
		worker_count -= 1
	}
}

//处理来自Inspector计算后的结果对象
func (c *Controller) processResult(strategy_result *types.StrategyResult) {
	item, ok := c.taskCache.Get(strategy_result.TaskID)
	if !ok {
		lg.Error(fmt.Sprintf("task %v not in cached task pool", strategy_result.TaskID))
		return
	}
	task, _ := item.(*types.AlarmTask)

	if strategy_result.ErrorMessage != "" {
		event_status := 0
		if strategy_result.ErrorMessage == "no data" {
			event_status = types.EVENT_NODATA
		} else {
			event_status = types.EVENT_UNKNOW
		}
		mydb.CreateStrategyEventFailed(task.Strategy.ID, task.Host.ID, event_status, strategy_result.ErrorMessage)
		lg.Error(fmt.Sprintf("task %v has problem %v", strategy_result.TaskID, strategy_result.ErrorMessage))
		return
	} else {
		go syncStrategyEventFailed(task.Strategy.ID, task.Host.ID)
	}

	var aware_strategy_event, old_strategy_event *types.StrategyEvent

	if aware_strategy_event = mydb.GetStrategyEvent(task.Strategy.ID, types.EVENT_AWARE, task.Host.ID); aware_strategy_event != nil {
		strategy_event, trigger_events := generateEvent(aware_strategy_event, strategy_result, task)
		if strategy_result.Triggered == false {
			triggered_trigger_events := mydb.GetTriggeredTriggerEvents(strategy_event.ID)
			trigger_events = tagChangedTrigger(triggered_trigger_events, trigger_events)
			strategy_event.Status = types.EVENT_CLOSED
			c.eventQueuesMutex.RLock()
			c.eventQueues[strategy_event.ProductID].PutNoWait(&QueueEvent{RESTORE_ALARM, strategy_event, trigger_events})
			c.eventQueuesMutex.RUnlock()
			lg.Debug(fmt.Sprintf("put restore event by strategy %s into event queue.", strategy_event.StrategyName))
			return
		}
		if time.Now().Sub(strategy_event.AwareEndTime) > 0 && strategy_event.AwareEndTime.Sub(types.DEFAULT_TIME) != 0 {
			strategy_event.Status = types.EVENT_NEW
			strategy_event.AwareEndTime = types.DEFAULT_TIME
			if strategy_event.Count < task.Strategy.AlarmCount || task.Strategy.AlarmCount == 0 {
				strategy_event.Count += 1
				c.eventQueuesMutex.RLock()
				c.eventQueues[strategy_event.ProductID].PutNoWait(&QueueEvent{OLD_ALARM, strategy_event, trigger_events})
				c.eventQueuesMutex.RUnlock()
				lg.Debug(fmt.Sprintf("put old alarm event by strategy %s into event queue.", strategy_event.StrategyName))
				return
			}
		}
		mydb.UpdateStrategyEvent(strategy_event, trigger_events)
		return
	}

	if old_strategy_event = mydb.GetStrategyEvent(task.Strategy.ID, types.EVENT_NEW, task.Host.ID); old_strategy_event != nil {
		strategy_event, trigger_events := generateEvent(old_strategy_event, strategy_result, task)
		if strategy_result.Triggered == false {
			triggered_trigger_events := mydb.GetTriggeredTriggerEvents(strategy_event.ID)
			trigger_events = tagChangedTrigger(triggered_trigger_events, trigger_events)
			strategy_event.Status = types.EVENT_CLOSED
			c.eventQueuesMutex.RLock()
			c.eventQueues[strategy_event.ProductID].PutNoWait(&QueueEvent{RESTORE_ALARM, strategy_event, trigger_events})
			c.eventQueuesMutex.RUnlock()
			lg.Debug(fmt.Sprintf("put restore event by strategy %s into event queue.", strategy_event.StrategyName))
			return
		}
		if strategy_event.Count >= task.Strategy.AlarmCount && task.Strategy.AlarmCount != 0 {
			mydb.UpdateStrategyEvent(strategy_event, trigger_events)
			return
		}
		strategy_event.Count += 1
		c.eventQueuesMutex.RLock()
		c.eventQueues[strategy_event.ProductID].PutNoWait(&QueueEvent{OLD_ALARM, strategy_event, trigger_events})
		c.eventQueuesMutex.RUnlock()
		lg.Debug(fmt.Sprintf("put old alarm event by strategy %s into event queue.", strategy_event.StrategyName))
		return
	}

	if old_strategy_event == nil && aware_strategy_event == nil {
		strategy_event, trigger_events := generateEvent(nil, strategy_result, task)
		if strategy_result.Triggered == true {
			strategy_event.Status = types.EVENT_NEW
			c.eventQueuesMutex.RLock()
			c.eventQueues[strategy_event.ProductID].PutNoWait(&QueueEvent{NEW_ALARM, strategy_event, trigger_events})
			c.eventQueuesMutex.RUnlock()
			lg.Debug(fmt.Sprintf("put new alarm event by strategy %s into event queue.", strategy_event.StrategyName))
		}
	}
}

func tagChangedTrigger(triggered_trigger_events []*types.TriggerEvent, trigger_events map[string]*types.TriggerEvent) map[string]*types.TriggerEvent {
	for _, e := range triggered_trigger_events {
		if value, ok := trigger_events[e.Index+e.Metric+e.Tags]; ok {
			value.TriggerChanged = true
		}
	}
	return trigger_events
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

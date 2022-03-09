package main

import (
	"errors"
	"sync"
	"time"

	"owl/common/types"
)

var (
	ErrTaskPoolFull   = errors.New("task pool is full")
	ErrResultPoolFull = errors.New("result pool is full")
	ErrEventPoolFull  = errors.New("event pool is full")
)

type TaskPool struct {
	tasks chan *types.AlarmTask
}

func NewTaskPool(size int) *TaskPool {
	return &TaskPool{make(chan *types.AlarmTask, size)}
}

// func (this *TaskPool) PutTasks(items map[string]cache.Item) error {
// 	if len(this.tasks) != 0 {
// 		this.clean()
// 	}
// 	for _, item := range items {
// 		task, _ := item.Object.(*types.AlarmTask)

// 		for {
// 			err := this.putTask(task)
// 			if err == nil {
// 				break
// 			}
// 			if err == ErrTaskPoolFull {
// 				expireTask := <-this.tasks
// 				lg.Warn("task pool is full, drop %s hostid:%s hostname:%s ip:%s",
// 					expireTask.Strategy.Name, expireTask.Host.ID, expireTask.Host.Hostname, expireTask.Host.IP)
// 			}
// 		}
// 	}
// 	return nil
// }

// func (this *TaskPool) clean() {
// 	for {
// 		select {
// 		case task := <-this.tasks:
// 			lg.Warn("drop expires task %s hostid:%s hostname:%s ip:%s",
// 				task.Strategy.Name, task.Host.ID, task.Host.Hostname, task.Host.IP)
// 		default:
// 			return
// 		}
// 	}
// }

func (tp *TaskPool) putTask(task *types.AlarmTask) error {
	select {
	case tp.tasks <- task:
		lg.Info("put new task into task pool, taskid:%s strategy:%s hostname:%s ip:%s",
			task.ID, task.Strategy.Name, task.Host.Hostname, task.Host.IP)
		return nil
	default:
		return ErrTaskPoolFull
	}
}

func (tp *TaskPool) getTasks(batchSize int) []*types.AlarmTask {
	tasks := make([]*types.AlarmTask, 0)
	for ; batchSize > 0; batchSize-- {
		if task := tp.getTask(); task != nil {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (tp *TaskPool) getTask() *types.AlarmTask {
	var task *types.AlarmTask
	select {
	case task = <-tp.tasks:
	default:
	}
	return task
}

type ResultPool struct {
	results chan *types.StrategyResult
}

func NewResultPool(size int) *ResultPool {
	return &ResultPool{make(chan *types.StrategyResult, size)}
}

func (this *ResultPool) PutResults(ar *types.AlarmResults) {
	for _, result := range ar.Results {
		select {
		case this.results <- result:
			lg.Debug("put task result into result pool, taskid:%s", result.TaskID)
		default:
			lg.Warn("result pool is full")
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}

func (this *ResultPool) putResult(result *types.StrategyResult) error {
	select {
	case this.results <- result:
		lg.Debug("put task result into result pool, taskid:%s", result.TaskID)
		return nil
	default:
		return ErrResultPoolFull
	}
}

type NodePool struct {
	Nodes map[string]*types.Node
	Lock  *sync.Mutex
}

func NewNodePool() *NodePool {
	return &NodePool{make(map[string]*types.Node), &sync.Mutex{}}
}

type EventPool struct {
	name        string
	events      chan *QueueEvent
	update_time time.Time
	mute        bool
}

func NewEventPool(name string, size int) *EventPool {
	return &EventPool{name, make(chan *QueueEvent, size), time.Now(), false}
}

func (ep *EventPool) putQueueEvent(event *QueueEvent) error {
	select {
	case ep.events <- event:
		return nil
	default:
		return ErrEventPoolFull
	}
}

func (ep *EventPool) getQueueEvent() *QueueEvent {
	select {
	case qEvent := <-ep.events:
		return qEvent
	}
}

func (ep *EventPool) cap() int {
	return cap(ep.events)
}

func (ep *EventPool) len() int {
	return len(ep.events)
}

func (ep *EventPool) clean() {
	for {
		select {
		case <-ep.events:
		default:
			return
		}
	}
}

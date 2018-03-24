package main

import (
	"owl/controller/cache"
	"sync"
	"time"

	"owl/common/types"
)

type TaskPool struct {
	tasks chan *types.AlarmTask
}

func NewTaskPool(size int) *TaskPool {
	return &TaskPool{make(chan *types.AlarmTask, size)}
}

func (this *TaskPool) PutTasks(items map[string]cache.Item) error {
	for {
		if len(this.tasks) != 0 {
			lg.Warn("task pool still have task: %v, wait it to be finished.", len(this.tasks))
			time.Sleep(time.Second * time.Duration(1))
			continue
		} else {
			break
		}
	}
	for _, item := range items {
		task, _ := item.Object.(*types.AlarmTask)
	LOOP:
		for {
			select {
			case this.tasks <- task:
				lg.Debug("load task %v into task pool", task.ID)
				break LOOP
			default:
				lg.Warn("task pool is full")
				time.Sleep(time.Second * time.Duration(1))
			}
		}
	}
	return nil
}

func (this *TaskPool) GetTasks(task_count int) []*types.AlarmTask {
	tasks := make([]*types.AlarmTask, 0)
OUTTER_LOOP:
	for task_count > 0 {
	INNER_LOOP:
		for {
			select {
			case task := <-this.tasks:
				tasks = append(tasks, task)
				task_count -= 1
				lg.Debug("get task %v from task pool", task.ID)
				break INNER_LOOP
			default:
				break OUTTER_LOOP
			}
		}
	}
	return tasks
}

type ResultPool struct {
	results chan *types.StrategyResult
}

func NewResultPool(size int) *ResultPool {
	return &ResultPool{make(chan *types.StrategyResult, size)}
}

func (this *ResultPool) PutResults(ar *types.AlarmResults) {
	for _, result := range ar.Results {
	LOOP:
		for {
			select {
			case this.results <- result:
				lg.Debug("load result %v into result pool", result.TaskID)
				break LOOP
			default:
				lg.Warn("result pool is full")
				time.Sleep(time.Second * time.Duration(1))
			}
		}
	}
}

type NodePool struct {
	Nodes map[string]*types.Node
	Lock  *sync.Mutex
}

func NewNodePool() *NodePool {
	return &NodePool{make(map[string]*types.Node), &sync.Mutex{}}
}

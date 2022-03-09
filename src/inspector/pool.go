package main

import (
	"time"

	"owl/common/types"
)

type TaskPool struct {
	tasks chan *types.AlarmTask
}

func NewTaskPool(size int) *TaskPool {
	return &TaskPool{make(chan *types.AlarmTask, size)}
}

func (this *TaskPool) PutTasks(tasks []*types.AlarmTask) {
	for _, task := range tasks {
	LOOP:
		for {
			select {
			case this.tasks <- task:
				lg.Debug("put task %v into task pool", task.ID)
				break LOOP
			default:
				lg.Warn("task pool is full")
				time.Sleep(time.Second * time.Duration(1))
			}
		}
	}
}

type ResultPool struct {
	results chan *types.StrategyResult
}

func NewResultPool(size int) *ResultPool {
	return &ResultPool{make(chan *types.StrategyResult, size)}
}

func (this *ResultPool) PutResult(result *types.StrategyResult) {
LOOP:
	for {
		select {
		case this.results <- result:
			lg.Debug("put result %v into result pool", result.TaskID)
			break LOOP
		default:
			lg.Warn("result pool is full")
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}

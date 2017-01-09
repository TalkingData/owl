package main

import (
	"time"

	. "owl/common/types"
)

type TaskPool struct {
	tasks chan *AlarmTask
}

func NewTaskPool(size int) *TaskPool {
	return &TaskPool{make(chan *AlarmTask, size)}
}

func (this *TaskPool) PutTasks(tasks map[string]*AlarmTask) error {
	for {
		if len(this.tasks) != 0 {
			lg.Warn("Task pool still have task: %v, wait it to be finished.", len(this.tasks))
			time.Sleep(time.Second * time.Duration(1))
			continue
		} else {
			break
		}
	}
	for _, task := range tasks {
	LOOP:
		for {
			select {
			case this.tasks <- task:
				lg.Debug("Load task %v into task pool", task.ID)
				break LOOP
			default:
				lg.Warn("Task pool is full")
				time.Sleep(time.Second * time.Duration(1))
			}
		}
	}
	return nil
}

func (this *TaskPool) GetTasks(task_count int) []*AlarmTask {
	tasks := make([]*AlarmTask, 0)
OUTTER_LOOP:
	for task_count > 0 {
	INNER_LOOP:
		for {
			select {
			case task := <-this.tasks:
				tasks = append(tasks, task)
				task_count -= 1
				lg.Debug("Get task %v from task pool", task.ID)
				break INNER_LOOP
			default:
				break OUTTER_LOOP
			}
		}
	}
	return tasks
}

type ResultPool struct {
	results chan *StrategyResult
}

func NewResultPool(size int) *ResultPool {
	return &ResultPool{make(chan *StrategyResult, size)}
}

func (this *ResultPool) PutResult(result *StrategyResult) {
LOOP:
	for {
		select {
		case this.results <- result:
			lg.Debug("Load result %v into result pool", result.TaskID)
			break LOOP
		default:
			lg.Warn("Result pool is full")
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}

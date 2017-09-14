package main

import (
	"encoding/json"
	"fmt"
	"owl/common/types"
	"owl/common/utils"
	"sync"
	"time"
)

var (
	tasklist = NewTaskList()
)

func init() {
}

type Task struct {
	types.Plugin
	timer *time.Ticker
	cycle int
	exit  chan struct{}
	flag  int
	args  []string
}

type TaskList struct {
	lock *sync.RWMutex
	M    map[int]*Task
}

func NewTaskList() *TaskList {
	tsl := &TaskList{}
	tsl.lock = new(sync.RWMutex)
	tsl.M = make(map[int]*Task)
	return tsl
}

func NewTask(p types.Plugin) *Task {
	t := new(Task)
	t.Plugin = p

	t.cycle = p.Interval
	t.exit = make(chan struct{})
	if len(p.Args) != 0 {
		t.args = parseCommandArgs(p.Args)
	}
	return t
}

func (this *Task) Do() {
	go func() {
		now := time.Now().Unix()
		diff := 60 - (now % 60)
		time.Sleep(time.Second * time.Duration(diff))
		this.timer = time.NewTicker(time.Duration(this.cycle) * time.Second)
		this.Run()
		for {
			select {
			case <-this.timer.C:
				this.Run()
			case <-this.exit:
				this.timer.Stop()
				return
			}
		}
	}()
}

func (this *Task) Stop() {
	close(this.exit)
}

func (this *Task) Run() {
	var (
		fpath string
		err   error
	)
	lg.Info("run %s %s", this.Name, this.args)
	fpath = fmt.Sprintf("./plugins/%s", this.Name)
	output, err := utils.RunCmdWithTimeout(fpath, this.args, this.Timeout)
	if err != nil {
		lg.Error("run %s/%s %s", fpath, this.args, err.Error())
		return
	}
	result := []types.TimeSeriesData{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		lg.Error("unmarshal task result error %s %s", output, err.Error())
		return
	}
	ts := time.Now().Unix()
	for _, tsd := range result {
		if tsd.Metric == "" || tsd.DataType == "" {
			continue
		}
		tsd.Cycle = this.Interval
		tsd.Timestamp = ts
		agent.SendChan <- tsd
	}

}

func (this *TaskList) AddTask(t *Task) {
	lg.Info("add task %v", t)
	this.lock.Lock()
	this.M[t.ID] = t
	this.lock.Unlock()
	t.Do()
}

func (this *TaskList) DelTask(pk int) {
	t, ok := this.M[pk]
	if ok {
		lg.Info("delete %s", t)
		t.Stop()
		this.lock.Lock()
		delete(this.M, pk)
		this.lock.Unlock()
	}
}

func (this *TaskList) All() map[int]*Task {
	this.lock.RLock()
	l := make(map[int]*Task, len(this.M))
	for k, v := range this.M {
		l[k] = v
	}
	this.lock.RUnlock()
	return l
}

func DelNotUsePlugin(pls []types.Plugin) {
	for _, t := range tasklist.All() {
		del := true
		for _, p := range pls {
			if t.ID == p.ID {
				del = false
				break
			}
		}
		if del {
			tasklist.DelTask(t.ID)
		}
	}
}

func MergePlugin(pls []types.Plugin) {
	for _, p := range pls {
		add := false
		if t, ok := tasklist.M[p.ID]; ok {
			if p.Name != t.Name || p.Args != t.Args || p.Interval != t.Interval {
				tasklist.DelTask(t.ID)
				add = true
			}
		} else {
			add = true
		}
		if add {
			t := NewTask(p)
			tasklist.AddTask(t)
		}
	}
}

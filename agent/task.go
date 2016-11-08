package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"owl/common/types"
	"strings"
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
	//t.timer = time.NewTicker(time.Duration(p.Interval) * time.Second)
	t.cycle = p.Interval
	t.exit = make(chan struct{})
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
		fpath  string
		args   []string
		err    error
		stderr bytes.Buffer
		stdout bytes.Buffer
		done   chan error = make(chan error, 1)
	)
	lg.Info("run %s %s", this.Name, this.Args)
	fpath = fmt.Sprintf("./plugins/%s", this.Name)
	if len(this.Args) != 0 {
		args = strings.Split(this.Args, " ")
	}
	cmd := exec.Command(fpath, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err = cmd.Start(); err != nil {
		lg.Error("start task error %s %s %s", this.Name, this.Args, err.Error())
		return
	}
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Second * time.Duration(this.Timeout)):
		if err = cmd.Process.Kill(); err != nil {
			lg.Error("failed to kill:", err.Error())
			return
		}
		lg.Warn("process killed as timeout reached, %s %s", fpath, this.Args)
	case err = <-done:
		if err != nil {
			lg.Error("task run error %s %s %s %s", fpath, this.Args, err.Error(), stderr.String())
		} else {
			result := []types.TimeSeriesData{}
			err = json.Unmarshal(stdout.Bytes(), &result)
			if err != nil {
				lg.Error("unmarshal task result error %s %s", stdout.String(), err.Error())
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
	}
}

func (this *TaskList) AddTask(t *Task) {
	lg.Info("add task %v", t)
	this.lock.RLock()
	this.M[t.ID] = t
	this.lock.RUnlock()
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

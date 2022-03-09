package main

import (
	"bytes"
	"encoding/json"
	"owl/common/types"
	"path/filepath"
	"sync"
	"time"

	"github.com/wuyingsong/utils"
)

var (
	tasklist = newTaskList()
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

func newTaskList() *TaskList {
	tsl := &TaskList{}
	tsl.lock = new(sync.RWMutex)
	tsl.M = make(map[int]*Task)
	return tsl
}

func newTask(p types.Plugin) *Task {
	t := new(Task)
	t.Plugin = p

	t.cycle = p.Interval
	t.exit = make(chan struct{})
	if len(p.Args) != 0 {
		t.args = utils.ParseCommandArgs(p.Args)
	}
	return t
}

func (task *Task) do() {
	go func() {
		now := time.Now().Unix()
		diff := 60 - (now % 60)
		time.Sleep(time.Second * time.Duration(diff))
		task.timer = time.NewTicker(time.Duration(task.cycle) * time.Second)
		task.run()
		for {
			select {
			case <-task.timer.C:
				task.run()
			case <-task.exit:
				task.timer.Stop()
				return
			}
		}
	}()
}

func (task *Task) stop() {
	close(task.exit)
}

func (task *Task) run() {
	var (
		err   error
		fpath string
	)
	lg.Info("run %s %s %s", task.Name, task.Path, task.args)
	fpath = filepath.Join(GlobalConfig.PluginDir, task.Path)
	if checksum, err := utils.GetFileMD5(fpath); err == nil {
		if checksum != task.Checksum {
			lg.Error("%s checksum verification failed, want(%s), have(%s)", fpath, task.Checksum, checksum)
			return
		}
	}
	ts := time.Now().Unix()
	ts = ts - (ts % int64(task.Interval))
	output, err := utils.RunCmdWithTimeout(fpath, task.args, task.Timeout)
	if err != nil {
		lg.Error("run %s %s %s %s", fpath, task.args, bytes.TrimSpace(output), err.Error())
		return
	}
	result := []types.TimeSeriesData{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		lg.Error("unmarshal task result error %s %s", output, err.Error())
		return
	}
	for _, tsd := range result {
		if err := tsd.Validate(); err != nil {
			lg.Warn("time series data validation failed %v, error:%s", tsd, err)
			continue
		}
		tsd.Cycle = task.Interval
		tsd.Timestamp = ts
		agent.SendChan <- tsd
	}

}

func (tl *TaskList) addTask(t *Task) {
	lg.Info("add task %v", t)
	tl.lock.Lock()
	tl.M[t.ID] = t
	tl.lock.Unlock()
	t.do()
}

func (tl *TaskList) removeTask(pk int) {
	if t, ok := tl.M[pk]; ok {
		lg.Info("remove task %s", t)
		t.stop()
		tl.lock.Lock()
		delete(tl.M, pk)
		tl.lock.Unlock()
	}
}

func (tl *TaskList) tasks() map[int]*Task {
	tl.lock.RLock()
	l := make(map[int]*Task, len(tl.M))
	for k, v := range tl.M {
		l[k] = v
	}
	tl.lock.RUnlock()
	return l
}

func removeNoUsePlugin(pls []types.Plugin) {
	for _, t := range tasklist.tasks() {
		del := true
		for _, p := range pls {
			if t.ID == p.ID {
				del = false
				break
			}
		}
		if del {
			lg.Info("remove no used plugin:%v", t.Plugin)
			tasklist.removeTask(t.ID)
		}
	}
}

func mergePlugin(pls []types.Plugin) {
	syncMap := make(map[string]struct{})
	for _, p := range pls {
		// 插件已经存在
		if t, ok := tasklist.M[p.ID]; ok {
			if p.Equal(t.Plugin) {
				goto sync
			}
			lg.Info("plugin change, old:%v, new:%v, removed", t.Plugin, p)
			tasklist.removeTask(t.ID)
		}
		tasklist.addTask(newTask(p))
	sync:
		canSync := true
		if checksum, err := utils.GetFileMD5(filepath.Join(GlobalConfig.PluginDir, p.Path)); err == nil {
			// 插件存在并且获取到checksum
			if checksum == p.Checksum {
				canSync = false
			}
		}
		if canSync {
			if _, ok := syncMap[p.Path]; ok {
				lg.Info("plugin %s already in sync queue, skiped", p.Path)
				continue
			}
			syncMap[p.Path] = struct{}{}
			agent.sendSyncPluginRequest(p)
		}
	}
}

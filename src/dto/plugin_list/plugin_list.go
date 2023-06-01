package builtin

import (
	"sync"
)

type PluginList struct {
	plugins map[string]*Plugin

	mu sync.RWMutex
}

// NewPluginList 新建插件列表
func NewPluginList() *PluginList {
	return &PluginList{plugins: make(map[string]*Plugin)}
}

func (pl *PluginList) StartAllPluginTask() {
	for _, v := range pl.List() {
		v.StartTask()
	}
}

func (pl *PluginList) StopAllPluginTask() {
	for _, v := range pl.List() {
		v.StopTask()
	}
}

// Put 新增
func (pl *PluginList) Put(pk string, p *Plugin) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.plugins[pk] = p
}

// StopTaskAndRemove 停止任务并移除
func (pl *PluginList) StopTaskAndRemove(pk string) {
	if p, ok := pl.plugins[pk]; ok {
		pl.mu.Lock()
		defer pl.mu.Unlock()

		p.StopTask()
		delete(pl.plugins, pk)
	}
}

// List 获取全部数据
func (pl *PluginList) List() map[string]*Plugin {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	return pl.plugins
}

// Len 长度
func (pl *PluginList) Len() int {
	return len(pl.plugins)
}

// Exists Task是否存在
func (pl *PluginList) Exists(pk string) bool {
	_, ok := pl.plugins[pk]
	return ok
}

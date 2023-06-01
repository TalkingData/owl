package builtin

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type PluginList struct {
	plugins map[string]*Plugin

	promMetric prometheus.Gauge
	mu         sync.RWMutex
}

// NewPluginList 新建插件列表
func NewPluginList() *PluginList {
	pm := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "plugin_list_length",
		Help: "The length of plugin list.",
	})
	prometheus.MustRegister(pm)
	pm.Set(0)

	return &PluginList{
		plugins:    make(map[string]*Plugin),
		promMetric: pm,
	}
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
	pl.promMetric.Set(float64(pl.Len()))
}

// StopAndRemoveTask 停止任务并移除
func (pl *PluginList) StopAndRemoveTask(pk string) {
	if p, ok := pl.plugins[pk]; ok {
		pl.mu.Lock()
		defer pl.mu.Unlock()

		p.StopTask()
		delete(pl.plugins, pk)
	}
	pl.promMetric.Set(float64(pl.Len()))
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

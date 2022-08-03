package metric_list

import (
	"sync"
)

type MetricList struct {
	mu      *sync.RWMutex
	metrics map[string]*metric
}

// NewMetricList 新建Metric列表
func NewMetricList() *MetricList {
	return &MetricList{
		mu:      new(sync.RWMutex),
		metrics: make(map[string]*metric),
	}
}

// Put 新增，存在相同Key的对象时会覆盖
func (ml *MetricList) Put(pk string, m *metric) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	ml.metrics[pk] = m
}

// List 获取全部数据
func (ml *MetricList) List() map[string]*metric {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	return ml.metrics
}

func (ml *MetricList) Get(pk string) (m *metric, exist bool) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	m, exist = ml.metrics[pk]
	return
}

// Len 长度
func (ml *MetricList) Len() int {
	return len(ml.metrics)
}

// Exists Task是否存在
func (ml *MetricList) Exists(k string) bool {
	_, ok := ml.metrics[k]
	return ok
}

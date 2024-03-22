package dto

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type TsDataMap struct {
	tsData map[string]*TsData

	promMetric prometheus.Gauge
	mu         sync.RWMutex
}

// NewTsDataMap 新建TsDataMap
func NewTsDataMap() *TsDataMap {
	pm := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ts_data_map_length",
		Help: "The length of time series data map.",
	})
	prometheus.MustRegister(pm)

	pm.Set(0)

	return &TsDataMap{
		promMetric: pm,
		tsData:     make(map[string]*TsData),
	}
}

// Put 新增，存在相同Key的对象时会覆盖
func (tsDataMap *TsDataMap) Put(pk string, m *TsData) {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	tsDataMap.tsData[pk] = m
	tsDataMap.promMetric.Set(float64(tsDataMap.Len()))
}

// Remove 删除
func (tsDataMap *TsDataMap) Remove(pk string) {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	delete(tsDataMap.tsData, pk)
	tsDataMap.promMetric.Set(float64(tsDataMap.Len()))
}

// List 获取全部数据
func (tsDataMap *TsDataMap) List() map[string]*TsData {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	res := make(map[string]*TsData)
	for k, v := range tsDataMap.tsData {
		res[k] = v.DeepCopyTsData()
	}
	return res
}

func (tsDataMap *TsDataMap) Get(pk string) (m *TsData, exist bool) {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	m, exist = tsDataMap.tsData[pk]
	return
}

// Len 长度
func (tsDataMap *TsDataMap) Len() int {
	return len(tsDataMap.tsData)
}

// Exists Task是否存在
func (tsDataMap *TsDataMap) Exists(k string) bool {
	_, ok := tsDataMap.tsData[k]
	return ok
}

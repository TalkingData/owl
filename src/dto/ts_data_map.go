package dto

import (
	"sync"
)

type TsDataMap struct {
	tsData map[string]*TsData

	mu sync.RWMutex
}

// NewTsDataMap 新建TsDataMap
func NewTsDataMap() *TsDataMap {
	return &TsDataMap{
		tsData: make(map[string]*TsData),
	}
}

// Put 新增，存在相同Key的对象时会覆盖
func (tsDataMap *TsDataMap) Put(pk string, m *TsData) {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	tsDataMap.tsData[pk] = m
}

// List 获取全部数据
func (tsDataMap *TsDataMap) List() map[string]*TsData {
	tsDataMap.mu.Lock()
	defer tsDataMap.mu.Unlock()

	return tsDataMap.tsData
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

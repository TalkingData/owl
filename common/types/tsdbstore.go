package types

import "sync"

type TimeSeriesDataStore struct {
	m    map[string]TimeSeriesData
	lock *sync.RWMutex
}

func NewTimeSeriesDataStore() *TimeSeriesDataStore {
	return &TimeSeriesDataStore{
		m:    make(map[string]TimeSeriesData),
		lock: new(sync.RWMutex),
	}
}

func (ts *TimeSeriesDataStore) Add(t TimeSeriesData) {
	ts.lock.Lock()
	ts.m[t.PK()] = t
	ts.lock.Unlock()
}

func (ts *TimeSeriesDataStore) Get(pk string) (tsd TimeSeriesData, exist bool) {
	ts.lock.RLock()
	tsd, ok := ts.m[pk]
	ts.lock.RUnlock()
	return tsd, ok
}

func (ts *TimeSeriesDataStore) Remove(pk string) {
	ts.lock.Lock()
	delete(ts.m, pk)
	ts.lock.Unlock()
}

func (ts *TimeSeriesDataStore) GetAll() []TimeSeriesData {
	tsdArr := []TimeSeriesData{}
	for _, t := range ts.clone().m {
		tsdArr = append(tsdArr, t)
	}
	return tsdArr
}

func (ts *TimeSeriesDataStore) clone() TimeSeriesDataStore {
	ts.lock.RLock()
	defer ts.lock.RUnlock()
	newts := NewTimeSeriesDataStore()
	for _, t := range ts.m {
		newts.Add(t)
	}
	return *newts
}

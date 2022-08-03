package metric_list

import (
	"fmt"
	cfcProto "owl/cfc/proto"
	"sort"
	"strings"
	"sync"
)

type metric struct {
	HostId    string
	Metric    string
	DataType  string
	Value     float64
	Timestamp int64
	Cycle     int32
	Tags      map[string]string

	mu *sync.RWMutex
}

func NewMetric(hostId, _metric, dataType string, val float64, ts int64, cycle int32, tags map[string]string) *metric {
	return &metric{
		HostId:    hostId,
		Metric:    _metric,
		DataType:  dataType,
		Value:     val,
		Timestamp: ts,
		Cycle:     cycle,
		Tags:      tags,

		mu: new(sync.RWMutex),
	}
}

func (m *metric) ToCfcMetric() *cfcProto.Metric {
	m.mu.Lock()
	defer m.mu.Unlock()

	return &cfcProto.Metric{
		HostId:   m.HostId,
		Metric:   m.Metric,
		DataType: m.DataType,
		Cycle:    m.Cycle,
		Tags:     m.Tags,
	}
}

func (m *metric) GetPk() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return fmt.Sprintf("%s.%s", m.Metric, m.tags2Str())
}

func (m *metric) tags2Str() (res string) {
	if len(m.Tags) == 0 {
		return
	}

	keyArr := []string{}
	for k := range m.Tags {
		tagStr := fmt.Sprintf("%s=%s", strings.TrimSpace(k), strings.TrimSpace(m.Tags[k]))
		keyArr = append(keyArr, tagStr)
	}

	sort.Strings(keyArr)
	return strings.Join(keyArr, ",")
}

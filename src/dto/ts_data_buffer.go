package dto

import (
	"github.com/prometheus/client_golang/prometheus"
	commonpb "owl/common/proto"
	"sync"
)

type TsDataBuffer struct {
	content []*commonpb.TsData

	pmLen         prometheus.Gauge
	pmBuffWriteCt prometheus.Counter
	mu            sync.Mutex
}

func NewTsDataBuffer() *TsDataBuffer {
	pmLen := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ts_data_buffer_length",
		Help: "The length of time series data buffer.",
	})
	prometheus.MustRegister(pmLen)
	pmLen.Set(0)

	pmWriteCt := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ts_data_buffer_write_times_total",
		Help: "The total times of time series data buffer.",
	})
	prometheus.MustRegister(pmWriteCt)
	pmWriteCt.Add(0)

	return &TsDataBuffer{
		content:       make([]*commonpb.TsData, 0),
		pmLen:         pmLen,
		pmBuffWriteCt: pmWriteCt,
	}
}

func (buf *TsDataBuffer) Put(data ...*commonpb.TsData) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.content = append(buf.content, data...)
	buf.pmLen.Set(float64(len(buf.content)))
	buf.pmBuffWriteCt.Add(float64(len(buf.content)))
}

func (buf *TsDataBuffer) Get(size int) []*commonpb.TsData {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	if size > len(buf.content) {
		size = len(buf.content)
	}

	batch := buf.content[:size]
	buf.content = buf.content[size:]
	buf.pmLen.Set(float64(len(buf.content)))
	return batch
}

func (buf *TsDataBuffer) Len() int {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	return len(buf.content)
}

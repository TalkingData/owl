package builtin

import (
	"github.com/shirou/gopsutil/load"
	"owl/common/types"
	"time"
)

func LoadMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range loadMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}

func loadMetrics(cycle int) []*types.TimeSeriesData {
	cnt, err := load.Avg()
	if err != nil {
		return nil
	}
	ts := time.Now().Unix()
	metrics := make([]*types.TimeSeriesData, 3)

	metrics[0] = &types.TimeSeriesData{
		Metric:    "sys.load1",
		Value:     cnt.Load1,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[1] = &types.TimeSeriesData{
		Metric:    "sys.load5",
		Value:     cnt.Load5,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[2] = &types.TimeSeriesData{
		Metric:    "sys.load15",
		Value:     cnt.Load15,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	return metrics
}

package builtin

import (
	"owl/common/types"
	"time"

	"github.com/shirou/gopsutil/mem"
)

func MemoryMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range memoryMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}

func memoryMetrics(cycle int) []*types.TimeSeriesData {
	metrics := make([]*types.TimeSeriesData, 9)
	cnt, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}
	ts := time.Now().Unix()
	metrics[0] = &types.TimeSeriesData{
		Metric:    "system.mem.active",
		Value:     float64(cnt.Active),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[1] = &types.TimeSeriesData{
		Metric:    "system.mem.available",
		Value:     float64(cnt.Available),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[2] = &types.TimeSeriesData{
		Metric:    "system.mem.buffers",
		Value:     float64(cnt.Buffers),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[3] = &types.TimeSeriesData{
		Metric:    "system.mem.cached",
		Value:     float64(cnt.Cached),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[4] = &types.TimeSeriesData{
		Metric:    "system.mem.free",
		Value:     float64(cnt.Free),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[5] = &types.TimeSeriesData{
		Metric:    "system.mem.inactive",
		Value:     float64(cnt.Inactive),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[6] = &types.TimeSeriesData{
		Metric:    "system.mem.total",
		Value:     float64(cnt.Total),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[7] = &types.TimeSeriesData{
		Metric:    "system.mem.used",
		Value:     float64(cnt.Used),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[8] = &types.TimeSeriesData{
		Metric:    "system.mem.used_pct",
		Value:     float64(cnt.UsedPercent),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	return metrics
}

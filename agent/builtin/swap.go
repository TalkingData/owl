package builtin

import (
	"github.com/shirou/gopsutil/mem"
	"owl/common/types"
	"time"
)

func SwapMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range swapMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}

func swapMetrics(cycle int) []*types.TimeSeriesData {
	metrics := make([]*types.TimeSeriesData, 6)
	cnt, err := mem.SwapMemory()
	if err != nil {
		return nil
	}
	ts := time.Now().Unix()
	metrics[0] = &types.TimeSeriesData{
		Metric:    "swap.total",
		Value:     float64(cnt.Total),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[1] = &types.TimeSeriesData{
		Metric:    "swap.usedprecent",
		Value:     cnt.UsedPercent,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[2] = &types.TimeSeriesData{
		Metric:    "swap.free",
		Value:     float64(cnt.Free),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[3] = &types.TimeSeriesData{
		Metric:    "swap.used",
		Value:     float64(cnt.Used),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[4] = &types.TimeSeriesData{
		Metric:    "swap.sin",
		Value:     float64(cnt.Sin),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[5] = &types.TimeSeriesData{
		Metric:    "swap.sout",
		Value:     float64(cnt.Sout),
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	return metrics
}

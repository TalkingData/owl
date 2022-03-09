package builtin

import (
	"owl/common/types"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func getAllTimes() []cpu.TimesStat {
	cnt := []cpu.TimesStat{}
	cnt1, _ := cpu.Times(true)
	cnt2, _ := cpu.Times(false)
	cnt = append(cnt, cnt1...)
	cnt = append(cnt, cnt2...)
	return cnt
}

func CpuMetrics(cycle int, channel chan types.TimeSeriesData) {
	for _, metric := range cpuMetrics(cycle) {
		if metric == nil {
			continue
		}
		channel <- *metric
	}
}

func cpuMetrics(cycle int) []*types.TimeSeriesData {
	metrics := []*types.TimeSeriesData{}
	ts := time.Now().Unix()
	t1 := getAllTimes()
	time.Sleep(time.Second * 1)
	t2 := getAllTimes()
	for idx, v2 := range t2 {
		v1 := t1[idx]
		total := v2.Total() - v1.Total()
		metrics = append(metrics,
			&types.TimeSeriesData{
				Metric:    "system.cpu.user",
				Value:     (v2.User - v1.User) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.system",
				Value:     (v2.System - v1.System) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.guest",
				Value:     (v2.Guest - v1.Guest) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.guestnice",
				Value:     (v2.GuestNice - v1.GuestNice) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.idle",
				Value:     (v2.Idle - v1.Idle) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.iowait",
				Value:     (v2.Iowait - v1.Iowait) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.irq",
				Value:     (v2.Irq - v1.Irq) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.nice",
				Value:     (v2.Nice - v1.Nice) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.softirq",
				Value:     (v2.Softirq - v1.Stolen) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.steal",
				Value:     (v2.Steal - v1.Steal) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
			&types.TimeSeriesData{
				Metric:    "system.cpu.stolen",
				Value:     (v2.Stolen - v1.Stolen) / total * 100,
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"cpu": v2.CPU},
			},
		)
	}

	return metrics
}

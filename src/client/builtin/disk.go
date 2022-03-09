package builtin

import (
	"owl/common/types"
	"time"

	"github.com/shirou/gopsutil/disk"
)

var (
	Cycle int
)

func DiskMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range diskMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}

func diskMetrics(cycle int) []*types.TimeSeriesData {
	metrics := []*types.TimeSeriesData{}
	Cycle = cycle
	metrics = append(metrics, diskPartitionMetrics()...)
	metrics = append(metrics, diskIOMetrics()...)
	return metrics
}

func diskIOMetrics() []*types.TimeSeriesData {
	cnt, err := disk.IOCounters()
	if err != nil {
		return nil
	}
	metrics := []*types.TimeSeriesData{}
	ts := time.Now().Unix()
	for _, v := range cnt {
		metrics = append(metrics,
			&types.TimeSeriesData{
				Metric:    "system.disk.bytes",
				Value:     float64(v.ReadBytes),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.count",
				Value:     float64(v.ReadCount),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.time",
				Value:     float64(v.ReadTime),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.bytes",
				Value:     float64(v.WriteBytes),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.count",
				Value:     float64(v.WriteCount),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.time",
				Value:     float64(v.WriteTime),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"device": v.Name, "direction": "out"},
			},
		)
	}
	return metrics
}
func diskPartitionMetrics() []*types.TimeSeriesData {
	pts, err := disk.Partitions(false)
	if err != nil {
		return nil
	}
	metrics := []*types.TimeSeriesData{}
	ts := time.Now().Unix()
	ts = ts - (ts % int64(Cycle))
	for _, p := range pts {
		cnt, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		metrics = append(metrics,
			&types.TimeSeriesData{
				Metric:    "system.disk.free",
				Value:     float64(cnt.Free),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.total",
				Value:     float64(cnt.Total),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.used",
				Value:     float64(cnt.Used),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.used_pct",
				Value:     cnt.UsedPercent,
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.inodes.free",
				Value:     float64(cnt.InodesFree),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.inodes.total",
				Value:     float64(cnt.InodesTotal),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.inodes.used",
				Value:     float64(cnt.InodesUsed),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
			&types.TimeSeriesData{
				Metric:    "system.disk.inodes.used_pct",
				Value:     float64(cnt.InodesUsedPercent),
				Cycle:     Cycle,
				Timestamp: ts,
				DataType:  "GAUGE",
				Tags:      map[string]string{"path": cnt.Path, "fstype": cnt.Fstype},
			},
		)

	}
	return metrics
}

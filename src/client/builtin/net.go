package builtin

import (
	"owl/common/types"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
)

var (
	IfaceNamePrefix = []string{"en", "em", "eth", "br"}
)

func NetMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range netMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}

func netMetrics(cycle int) []*types.TimeSeriesData {
	cnt, err := net.IOCounters(true)
	if err != nil {
		return nil
	}
	ts := time.Now().Unix()
	metrics := []*types.TimeSeriesData{}

	for _, v := range cnt {
		if NameNotAvalid(v.Name) {
			continue
		}
		metrics = append(metrics,
			&types.TimeSeriesData{
				Metric:    "system.net.bytes",
				Value:     float64(v.BytesSent),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.bytes",
				Value:     float64(v.BytesRecv),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.packets",
				Value:     float64(v.PacketsSent),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.packets",
				Value:     float64(v.PacketsRecv),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.err",
				Value:     float64(v.Errin),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.err",
				Value:     float64(v.Errout),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.drop",
				Value:     float64(v.Dropin),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.drop",
				Value:     float64(v.Dropout),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.fifo",
				Value:     float64(v.Fifoin),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "in"},
			},
			&types.TimeSeriesData{
				Metric:    "system.net.fifo",
				Value:     float64(v.Fifoout),
				Cycle:     cycle,
				Timestamp: ts,
				DataType:  "COUNTER",
				Tags:      map[string]string{"iface": v.Name, "direction": "out"},
			},
		)
	}

	return metrics
}

func NameNotAvalid(name string) bool {
	for _, n := range IfaceNamePrefix {
		if strings.HasPrefix(name, n) {
			return false
		}
	}
	return true
}

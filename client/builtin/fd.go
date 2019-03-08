package builtin

import (
	"bufio"
	"os"
	"owl/common/types"
	"strconv"
	"strings"
	"time"
)

const (
	filename = "/proc/sys/fs/file-nr"
)

func FdMetrics(cycle int, ch chan types.TimeSeriesData) {
	for _, metric := range fdMetrics(cycle) {
		if metric == nil {
			continue
		}
		ch <- *metric
	}
}
func fdMetrics(cycle int) []*types.TimeSeriesData {
	var (
		allocated, max, unused float64
		err                    error
	)
	fd, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer fd.Close()
	ts := time.Now().Unix()
	r := bufio.NewReader(fd)
	line, err := r.ReadString('\n')
	if err != nil {
		return nil
	}
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil
	}
	allocated, _ = strconv.ParseFloat(fields[0], 64)
	unused, _ = strconv.ParseFloat(fields[1], 64)
	max, _ = strconv.ParseFloat(fields[2], 64)

	metrics := make([]*types.TimeSeriesData, 4)
	metrics[0] = &types.TimeSeriesData{
		Metric:    "system.fd.allocated",
		Value:     allocated,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[1] = &types.TimeSeriesData{
		Metric:    "system.fd.unused",
		Value:     unused,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[2] = &types.TimeSeriesData{
		Metric:    "system.fd.max",
		Value:     max,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	metrics[3] = &types.TimeSeriesData{
		Metric:    "system.fd.used_pct",
		Value:     (allocated / max) * 100,
		Cycle:     cycle,
		Timestamp: ts,
		DataType:  "GAUGE",
	}
	return metrics
}

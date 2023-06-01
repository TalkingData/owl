package executor

import (
	"bufio"
	"os"
	"owl/common/logger"
	"owl/dto"
	"strconv"
	"strings"
)

const (
	fdFile = "/proc/sys/fs/file-nr"
)

func (e *Executor) ExecCollectFd(ts int64, cycle int32) dto.TsDataArray {
	e.logger.Info("Executor.ExecCollectFd called.")
	defer e.logger.Info("Executor.ExecCollectFd end.")

	var allocated, max, unused float64

	fd, err := os.Open(fdFile)
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"fd_file": fdFile,
			"cycle":   cycle,
			"error":   err,
		}, "An error occurred while calling os.Open.")
		return nil
	}
	defer func() {
		_ = fd.Close()
	}()

	r := bufio.NewReader(fd)
	line, err := r.ReadString('\n')
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while calling r.ReadString.")
		return nil
	}
	fields := strings.Fields(line)
	if len(fields) < 3 {
		e.logger.ErrorWithFields(logger.Fields{
			"fields_length": len(fields),
			"cycle":         cycle,
			"error":         err,
		}, "An error occurred while calling strings.Fields, len(fields) < 3.")
		return nil
	}

	allocated, _ = strconv.ParseFloat(fields[0], 64)
	unused, _ = strconv.ParseFloat(fields[1], 64)
	max, _ = strconv.ParseFloat(fields[2], 64)

	return dto.TsDataArray{
		{
			Metric:    "system.fd.allocated",
			DataType:  dto.TsDataTypeGauge,
			Value:     allocated,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.fd.unused",
			DataType:  dto.TsDataTypeGauge,
			Value:     unused,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.fd.max",
			DataType:  dto.TsDataTypeGauge,
			Value:     max,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.fd.used_pct",
			DataType:  dto.TsDataTypeGauge,
			Value:     (allocated / max) * 100,
			Timestamp: ts,
			Cycle:     cycle,
		},
	}
}

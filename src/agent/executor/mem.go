package executor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectMem(cycle int32) dto.TsDataArray {
	e.logger.Info("Executor.ExecCollectMem called.")
	defer e.logger.Info("Executor.ExecCollectMem end.")

	memStat, err := mem.VirtualMemory()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while Executor.ExecCollectMem.")
		return nil
	}

	currTs := time.Now().Unix()
	return dto.TsDataArray{
		{
			Metric:    "system.mem.active",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Active),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.available",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Available),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.buffers",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Buffers),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.cached",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Cached),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.free",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Free),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.inactive",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Inactive),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.total",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Total),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.used",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Used),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.used_pct",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.UsedPercent),
			Timestamp: currTs,
			Cycle:     cycle,
		},
	}
}

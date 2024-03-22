package executor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"owl/common/logger"
	"owl/dto"
)

func (e *Executor) ExecCollectMem(ts int64, cycle int32) dto.TsDataArray {
	e.logger.Info("Executor.ExecCollectMem called.")
	defer e.logger.Info("Executor.ExecCollectMem end.")

	memStat, err := mem.VirtualMemory()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while calling Executor.ExecCollectMem.")
		return nil
	}

	return dto.TsDataArray{
		{
			Metric:    "system.mem.active",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Active),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.available",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Available),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.buffers",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Buffers),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.cached",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Cached),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.free",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Free),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.inactive",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Inactive),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.total",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Total),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.used",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.Used),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.mem.used_pct",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(memStat.UsedPercent),
			Timestamp: ts,
			Cycle:     cycle,
		},
	}
}

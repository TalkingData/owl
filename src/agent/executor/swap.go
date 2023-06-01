package executor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"owl/common/logger"
	"owl/dto"
)

func (e *Executor) ExecCollectSwap(ts int64, cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.ExecCollectSwap called.")
	defer e.logger.Info("Executor.ExecCollectSwap end.")

	swapStat, err := mem.SwapMemory()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while calling Executor.ExecCollectSwap.")
		return nil
	}

	return dto.TsDataArray{
		{
			Metric:    "system.swap.total",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Total),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.used_pct",
			DataType:  dto.TsDataTypeGauge,
			Value:     swapStat.UsedPercent,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.free",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Free),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.used",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Used),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.sin",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Sin),
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.sout",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Sout),
			Timestamp: ts,
			Cycle:     cycle,
		},
	}
}

package executor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectSwap(cycle int32) (res dto.TsDataArray) {
	e.logger.Info("Executor.ExecCollectSwap called.")
	defer e.logger.Info("Executor.ExecCollectSwap end.")

	swapStat, err := mem.SwapMemory()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while Executor.ExecCollectSwap.")
		return nil
	}
	currTs := time.Now().Unix()
	return dto.TsDataArray{
		{
			Metric:    "system.swap.total",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Total),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.used_pct",
			DataType:  dto.TsDataTypeGauge,
			Value:     swapStat.UsedPercent,
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.free",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Free),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.used",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Used),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.sin",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Sin),
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.swap.sout",
			DataType:  dto.TsDataTypeGauge,
			Value:     float64(swapStat.Sout),
			Timestamp: currTs,
			Cycle:     cycle,
		},
	}
}

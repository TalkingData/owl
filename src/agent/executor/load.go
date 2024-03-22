package executor

import (
	"github.com/shirou/gopsutil/v3/load"
	"owl/common/logger"
	"owl/dto"
)

func (e *Executor) ExecCollectLoad(ts int64, cycle int32) dto.TsDataArray {
	e.logger.Info("Executor.ExecCollectLoad called.")
	defer e.logger.Info("Executor.ExecCollectLoad end.")

	loadAvgStat, err := load.Avg()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while calling Executor.ExecCollectLoad.")
		return nil
	}

	return dto.TsDataArray{
		{
			Metric:    "system.load.1min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load1,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.load.5min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load5,
			Timestamp: ts,
			Cycle:     cycle,
		},
		{
			Metric:    "system.load.15min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load15,
			Timestamp: ts,
			Cycle:     cycle,
		},
	}
}

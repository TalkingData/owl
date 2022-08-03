package executor

import (
	"github.com/shirou/gopsutil/v3/load"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectLoad(cycle int32) []*dto.TsData {
	e.logger.Info("Executor.ExecCollectLoad called.")
	defer e.logger.Info("Executor.ExecCollectLoad end.")

	loadAvgStat, err := load.Avg()
	if err != nil {
		e.logger.ErrorWithFields(logger.Fields{
			"cycle": cycle,
			"error": err,
		}, "An error occurred while Executor.ExecCollectLoad.")
		return nil
	}

	currTs := time.Now().Unix()
	return []*dto.TsData{
		{
			Metric:    "system.load.1min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load1,
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.load.5min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load5,
			Timestamp: currTs,
			Cycle:     cycle,
		},
		{
			Metric:    "system.load.15min",
			DataType:  dto.TsDataTypeGauge,
			Value:     loadAvgStat.Load15,
			Timestamp: currTs,
			Cycle:     cycle,
		},
	}
}

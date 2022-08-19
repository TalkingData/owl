package executor

import (
	"owl/dto"
	"time"
)

func (e *Executor) ExecCollectAlive(cycle int32) dto.TsDataArray {
	return dto.TsDataArray{{
		Metric:    "agent.alive",
		DataType:  dto.TsDataTypeGauge,
		Value:     1,
		Timestamp: time.Now().Unix(),
		Cycle:     cycle,
	}}
}

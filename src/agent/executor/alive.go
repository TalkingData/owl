package executor

import (
	"owl/dto"
)

func (e *Executor) ExecCollectAlive(ts int64, cycle int32) dto.TsDataArray {
	return dto.TsDataArray{{
		Metric:    "agent.alive",
		DataType:  dto.TsDataTypeGauge,
		Value:     1,
		Timestamp: ts,
		Cycle:     cycle,
	}}
}

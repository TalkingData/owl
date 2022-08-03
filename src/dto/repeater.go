package dto

import (
	repProto "owl/repeater/proto"
)

// TransRepeaterTsData2Dto 将RepeaterTsData转为dto.TsData
func TransRepeaterTsData2Dto(in *repProto.TsData) *TsData {
	return &TsData{
		Metric:    in.Metric,
		DataType:  in.DataType,
		Value:     in.Value,
		Cycle:     in.Cycle,
		Timestamp: in.Timestamp,
		Tags:      in.Tags,
	}
}

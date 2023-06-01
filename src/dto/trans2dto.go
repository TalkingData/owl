package dto

import (
	commonpb "owl/common/proto"
)

// TransRepeaterTsData2Dto 将RepeaterTsData转为dto.TsData
func TransRepeaterTsData2Dto(in *commonpb.TsData) *TsData {
	return &TsData{
		Metric:    in.Metric,
		DataType:  in.DataType,
		Value:     in.Value,
		Cycle:     in.Cycle,
		Timestamp: in.Timestamp,
		Tags:      deepCopyTags(in.Tags),
	}
}

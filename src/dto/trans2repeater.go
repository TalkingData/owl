package dto

import (
	proxyProto "owl/proxy/proto"
	repProto "owl/repeater/proto"
)

func TransProxyTsData2Repeater(in *proxyProto.TsData) *repProto.TsData {
	return &repProto.TsData{
		Metric:    in.Metric,
		DataType:  in.DataType,
		Value:     in.Value,
		Timestamp: in.Timestamp,
		Cycle:     in.Cycle,
		Tags:      in.Tags,
	}
}

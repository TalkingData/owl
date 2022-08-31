package dto

import (
	cfcProto "owl/cfc/proto"
	proxyProto "owl/proxy/proto"
)

func TransProxyAgentInfo2Cfc(in *proxyProto.AgentInfo) *cfcProto.AgentInfo {
	return &cfcProto.AgentInfo{
		HostId:       in.HostId,
		Ip:           in.Ip,
		Hostname:     in.Hostname,
		AgentVersion: in.AgentVersion,
		AgentOs:      in.AgentOs,
		AgentArch:    in.AgentArch,
		Uptime:       in.Uptime,
		IdlePct:      in.IdlePct,
		Metadata:     in.Metadata,
	}
}

func TransProxyMetric2Cfc(in *proxyProto.Metric) *cfcProto.Metric {
	return &cfcProto.Metric{
		HostId:   in.HostId,
		Metric:   in.Metric,
		DataType: in.DataType,
		Cycle:    in.Cycle,
		Tags:     in.Tags,
	}
}

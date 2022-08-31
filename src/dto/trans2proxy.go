package dto

import (
	cfcProto "owl/cfc/proto"
	proxyProto "owl/proxy/proto"
)

func TransCfcPlugin2Proxy(in *cfcProto.Plugin) *proxyProto.Plugin {
	return &proxyProto.Plugin{
		Id:       in.Id,
		Name:     in.Name,
		Path:     in.Path,
		Checksum: in.Checksum,
		Args:     in.Args,
		Interval: in.Interval,
		Timeout:  in.Timeout,
		Comment:  in.Comment,
	}
}

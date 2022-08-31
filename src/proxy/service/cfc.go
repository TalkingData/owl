package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	cfcProto "owl/cfc/proto"
	"owl/common/logger"
	"owl/dto"
	proxyProto "owl/proxy/proto"
)

func (proxySrv *OwlProxyService) ReceiveAgentHeartbeat(ctx context.Context, req *proxyProto.AgentInfo) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveAgentHeartbeat called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveAgentHeartbeat end.")

	empty := new(emptypb.Empty)

	_, err := proxySrv.cfcCli.ReceiveAgentHeartbeat(ctx, dto.TransProxyAgentInfo2Cfc(req))
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while proxySrv.cfcCli.ReceiveAgentHeartbeat, Skipped.")
		return empty, err
	}

	return empty, nil
}

func (proxySrv *OwlProxyService) ReceiveAgentMetric(ctx context.Context, req *proxyProto.Metric) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveAgentMetric called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveAgentMetric end.")

	empty := new(emptypb.Empty)

	_, err := proxySrv.cfcCli.ReceiveAgentMetric(ctx, dto.TransProxyMetric2Cfc(req))
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while proxySrv.cfcCli.ReceiveAgentMetric, Skipped.")
		return empty, err
	}

	return empty, nil
}

func (proxySrv *OwlProxyService) RegisterAgent(ctx context.Context, req *proxyProto.AgentInfo) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.RegisterAgent called.")
	defer proxySrv.logger.Debug("proxySrv.RegisterAgent end.")

	empty := new(emptypb.Empty)

	_, err := proxySrv.cfcCli.RegisterAgent(ctx, dto.TransProxyAgentInfo2Cfc(req))
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while proxySrv.cfcCli.RegisterAgent, Skipped.")
		return empty, err
	}

	return empty, nil
}

func (proxySrv *OwlProxyService) ListAgentPlugins(ctx context.Context, req *proxyProto.HostIdReq) (*proxyProto.Plugins, error) {
	proxySrv.logger.Debug("proxySrv.ListAgentPlugins called.")
	defer proxySrv.logger.Debug("proxySrv.ListAgentPlugins end.")

	// 构造返回值
	ret := &proxyProto.Plugins{Plugins: make([]*proxyProto.Plugin, 0)}

	plugins, err := proxySrv.cfcCli.ListAgentPlugins(ctx, &cfcProto.HostIdReq{HostId: req.HostId})
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while proxySrv.cfcCli.ListAgentPlugins.")
		return ret, err
	}

	for _, p := range plugins.Plugins {
		ret.Plugins = append(ret.Plugins, dto.TransCfcPlugin2Proxy(p))
	}

	return ret, nil
}

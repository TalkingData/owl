package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/common/logger"
	commonpb "owl/common/proto"
)

func (proxySrv *OwlProxyService) ReceiveAgentHeartbeat(
	ctx context.Context, req *commonpb.AgentInfo,
) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveAgentHeartbeat called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveAgentHeartbeat end.")

	ret, err := proxySrv.cfcCli.ReceiveAgentHeartbeat(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.cfcCli.ReceiveAgentHeartbeat, Skipped.")
		return ret, err
	}

	return ret, nil
}

func (proxySrv *OwlProxyService) ReceiveAgentMetric(ctx context.Context, req *commonpb.Metric) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveAgentMetric called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveAgentMetric end.")

	ret, err := proxySrv.cfcCli.ReceiveAgentMetric(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.cfcCli.ReceiveAgentMetric, Skipped.")
		return ret, err
	}

	return ret, nil
}

func (proxySrv *OwlProxyService) ReceiveAgentMetrics(
	ctx context.Context, req *commonpb.Metrics,
) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveAgentMetrics called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveAgentMetrics end.")

	ret, err := proxySrv.cfcCli.ReceiveAgentMetrics(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.cfcCli.ReceiveAgentMetrics, Skipped.")
		return ret, err
	}

	return ret, nil
}

func (proxySrv *OwlProxyService) RegisterAgent(ctx context.Context, req *commonpb.AgentInfo) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.RegisterAgent called.")
	defer proxySrv.logger.Debug("proxySrv.RegisterAgent end.")

	ret, err := proxySrv.cfcCli.RegisterAgent(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.cfcCli.RegisterAgent, Skipped.")
		return ret, err
	}

	return ret, nil
}

func (proxySrv *OwlProxyService) ListAgentPlugins(
	ctx context.Context, req *commonpb.HostIdReq,
) (*commonpb.Plugins, error) {
	proxySrv.logger.Debug("proxySrv.ListAgentPlugins called.")
	defer proxySrv.logger.Debug("proxySrv.ListAgentPlugins end.")

	plugins, err := proxySrv.cfcCli.ListAgentPlugins(ctx, &commonpb.HostIdReq{HostId: req.HostId})
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.cfcCli.ListAgentPlugins, Skipped.")
		return plugins, err
	}

	return plugins, nil
}

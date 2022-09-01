package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/cfc/biz"
	"owl/cfc/conf"
	cfcProto "owl/cfc/proto"
	"owl/common/logger"
	"owl/common/utils"
	"owl/dao"
)

// OwlCfcService struct
type OwlCfcService struct {
	dao *dao.Dao

	biz            *biz.Biz
	grpcDownloader *utils.Downloader

	conf   *conf.Conf
	logger *logger.Logger
}

// NewOwlCfcService 新建Cfc服务
func NewOwlCfcService(dao *dao.Dao, conf *conf.Conf, logger *logger.Logger) *OwlCfcService {
	return &OwlCfcService{
		dao: dao,

		biz:            biz.NewBiz(dao, conf, logger),
		grpcDownloader: new(utils.Downloader),

		conf:   conf,
		logger: logger,
	}
}

// RegisterAgent 客户端注册
func (cfcSrv *OwlCfcService) RegisterAgent(_ context.Context, agent *cfcProto.AgentInfo, _ *emptypb.Empty) error {
	cfcSrv.logger.Debug("cfcSrv.RegisterAgent called.")
	defer cfcSrv.logger.Debug("cfcSrv.RegisterAgent end.")

	// 使用业务层分离处理客户端注册
	err := cfcSrv.biz.RegisterAgent(
		agent.HostId, agent.Ip, agent.Hostname, agent.AgentVersion, agent.Uptime, agent.IdlePct, agent.Metadata,
	)
	if err != nil {
		// 汇报类操作，错误不需要返给Agent，记录日志即可
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  agent.HostId,
			"agent_ip":       agent.Ip,
			"agent_hostname": agent.Hostname,
			"error":          err,
		}, "An error occurred while biz.RegisterAgent in cfcSrv.RegisterAgent.")
	}

	return nil
}

// ListAgentPlugins 列出客户端需要执行的插件
func (cfcSrv *OwlCfcService) ListAgentPlugins(_ context.Context, req *cfcProto.HostIdReq, rsp *cfcProto.Plugins) error {
	cfcSrv.logger.Debug("cfcSrv.ListAgentPlugins called.")
	defer cfcSrv.logger.Debug("cfcSrv.ListAgentPlugins end.")

	// 构造返回值
	rsp.Plugins = make([]*cfcProto.Plugin, 0)

	// 使用业务层分离处理
	plugins, err := cfcSrv.biz.ListAgentPlugins(req.HostId)
	if err != nil {
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"agent_host_id": req.HostId,
			"error":         err,
		}, "An error occurred while biz.ListAgentPlugins in cfcSrv.ListAgentPlugins.")
		return err
	}
	// 填充返回内容
	for _, p := range plugins {
		rsp.Plugins = append(rsp.Plugins, &cfcProto.Plugin{
			Id:       uint32(p.Id),
			Name:     p.Name,
			Path:     p.Path,
			Checksum: p.Checksum,
			Args:     p.Args,
			Interval: int32(p.Interval),
			Timeout:  int32(p.Timeout),
			Comment:  p.Comment,
		})
	}

	return nil
}

// ReceiveAgentHeartbeat 接收客户端上报的心跳
func (cfcSrv *OwlCfcService) ReceiveAgentHeartbeat(_ context.Context, agent *cfcProto.AgentInfo, _ *emptypb.Empty) error {
	cfcSrv.logger.Debug("cfcSrv.ReceiveAgentHeartbeat called.")
	defer cfcSrv.logger.Debug("cfcSrv.ReceiveAgentHeartbeat end.")

	// 使用业务层分离处理客户端注册
	err := cfcSrv.biz.ReceiveAgentHeartbeat(
		agent.HostId, agent.Ip, agent.Hostname, agent.AgentVersion, agent.Uptime, agent.IdlePct,
	)
	if err != nil {
		// 汇报类操作，错误不需要返给Agent，记录日志即可
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  agent.HostId,
			"agent_ip":       agent.Ip,
			"agent_hostname": agent.Hostname,
			"error":          err,
		}, "An error occurred while biz.ReceiveAgentHeartbeat in cfcSrv.ReceiveAgentHeartbeat.")
	}

	return nil
}

// ReceiveAgentMetric 接收客户端上报的Metric
func (cfcSrv *OwlCfcService) ReceiveAgentMetric(_ context.Context, metric *cfcProto.Metric, _ *emptypb.Empty) error {
	cfcSrv.logger.Debug("cfcSrv.ReceiveAgentMetric called.")
	defer cfcSrv.logger.Debug("cfcSrv.ReceiveAgentMetric end.")

	// 使用业务层分离处理客户端注册
	cfcSrv.biz.ReceiveAgentMetric(
		metric.HostId,
		metric.Metric,
		metric.DataType,
		metric.Cycle,
		metric.Tags,
	)

	return nil
}
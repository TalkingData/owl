package service

import (
	"context"
	"errors"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"owl/cfc/biz"
	"owl/cfc/conf"
	cfcProto "owl/cfc/proto"
	"owl/common/logger"
	"owl/common/orm"
	"owl/common/utils"
	"owl/dao"
	"path"
	"path/filepath"
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
func NewOwlCfcService(biz *biz.Biz, dao *dao.Dao, conf *conf.Conf, logger *logger.Logger) *OwlCfcService {
	return &OwlCfcService{
		dao: dao,

		biz:            biz,
		grpcDownloader: new(utils.Downloader),

		conf:   conf,
		logger: logger,
	}
}

// RegisterAgent 客户端注册
func (cfcSrv *OwlCfcService) RegisterAgent(_ context.Context, agent *cfcProto.AgentInfo) (*emptypb.Empty, error) {
	cfcSrv.logger.Debug("cfcSrv.RegisterAgent called.")
	defer cfcSrv.logger.Debug("cfcSrv.RegisterAgent end.")

	empty := new(emptypb.Empty)

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

	return empty, nil
}

// ListAgentPlugins 列出客户端需要执行的插件
func (cfcSrv *OwlCfcService) ListAgentPlugins(_ context.Context, req *cfcProto.HostIdReq) (*cfcProto.Plugins, error) {
	cfcSrv.logger.Debug("cfcSrv.ListAgentPlugins called.")
	defer cfcSrv.logger.Debug("cfcSrv.ListAgentPlugins end.")

	// 构造返回值
	ret := &cfcProto.Plugins{Plugins: make([]*cfcProto.Plugin, 0)}

	// 使用业务层分离处理
	plugins, err := cfcSrv.biz.ListAgentPlugins(req.HostId)
	if err != nil {
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"agent_host_id": req.HostId,
			"error":         err,
		}, "An error occurred while biz.ListAgentPlugins in cfcSrv.ListAgentPlugins.")
		return ret, err
	}
	// 填充返回内容
	for _, p := range plugins {
		ret.Plugins = append(ret.Plugins, &cfcProto.Plugin{
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

	return ret, nil
}

// DownloadPluginFile 下载插件文件
func (cfcSrv *OwlCfcService) DownloadPluginFile(req *cfcProto.PluginIdReq, stream cfcProto.OwlCfcService_DownloadPluginFileServer) error {
	cfcSrv.logger.Debug("cfcSrv.DownloadPluginFile called.")
	defer cfcSrv.logger.Debug("cfcSrv.DownloadPluginFile end.")

	plugin, err := cfcSrv.dao.GetPlugin(orm.Query{"id": req.PluginId})
	if err != nil {
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while dao.GetPlugin in cfcSrv.DownloadPluginFile.")
		return err
	}

	if plugin == nil || plugin.Id < 1 {
		err = errors.New("plugin not found")
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while dao.GetPlugin in cfcSrv.DownloadPluginFile.")
		return err
	}

	pluginPathname, err := filepath.Abs(path.Join(cfcSrv.conf.PluginDir, plugin.Path))
	if err != nil {
		cfcSrv.logger.ErrorWithFields(logger.Fields{
			"plugin_dir":  cfcSrv.conf.PluginDir,
			"plugin_path": plugin.Path,
			"error":       err,
		}, "An error occurred while filepath.Abs in cfcSrv.DownloadPluginFile.")
	}

	cfcSrv.logger.ErrorWithFields(logger.Fields{
		"plugin_pathname": pluginPathname,
	}, "cfcSrv.DownloadPluginFile prepare send plugin file.")
	err = cfcSrv.grpcDownloader.Download(pluginPathname, func(buffer []byte) error {
		return stream.Send(&cfcProto.PluginFile{Buffer: buffer})
	})
	if err != nil {
		switch err {
		case utils.ErrEndOfFileExit:
			return status.Error(utils.DefaultDownloaderEndOfFileExitCode, io.EOF.Error())
		case utils.ErrNormallyExit:
			return status.Error(utils.DefaultDownloaderNormallyExitCode, "normally exit")
		}
	}

	return err
}

// ReceiveAgentHeartbeat 接收客户端上报的心跳
func (cfcSrv *OwlCfcService) ReceiveAgentHeartbeat(_ context.Context, agent *cfcProto.AgentInfo) (*emptypb.Empty, error) {
	cfcSrv.logger.Debug("cfcSrv.ReceiveAgentHeartbeat called.")
	defer cfcSrv.logger.Debug("cfcSrv.ReceiveAgentHeartbeat end.")

	empty := new(emptypb.Empty)

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

	return empty, nil
}

// ReceiveAgentMetric 接收客户端上报的Metric
func (cfcSrv *OwlCfcService) ReceiveAgentMetric(_ context.Context, metric *cfcProto.Metric) (*emptypb.Empty, error) {
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

	return new(emptypb.Empty), nil
}

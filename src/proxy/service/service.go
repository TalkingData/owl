package service

import (
	cfcProto "owl/cfc/proto"
	"owl/cli"
	"owl/common/logger"
	"owl/common/utils"
	"owl/proxy/conf"
	repProto "owl/repeater/proto"
)

// OwlProxyService struct
type OwlProxyService struct {
	cfcCli cfcProto.OwlCfcService
	repCli repProto.OwlRepeaterService

	grpcDownloader *utils.Downloader

	conf   *conf.Conf
	logger *logger.Logger
}

// NewOwlProxyService 新建Proxy服务
func NewOwlProxyService(conf *conf.Conf, logger *logger.Logger) *OwlProxyService {
	return &OwlProxyService{
		cfcCli: cli.NewCfcClient(conf.EtcdUsername, conf.EtcdPassword, conf.EtcdAddresses),
		repCli: cli.NewRepeaterClient(conf.EtcdUsername, conf.EtcdPassword, conf.EtcdAddresses),

		grpcDownloader: new(utils.Downloader),

		conf:   conf,
		logger: logger,
	}
}

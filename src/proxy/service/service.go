package service

import (
	"github.com/micro/go-micro/v2/client"
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
		cfcCli: cli.NewCfcClient(
			conf.EtcdUsername,
			conf.EtcdPassword,
			conf.EtcdAddresses,
			client.Retries(conf.CallCfcRetries),
		),
		repCli: cli.NewRepeaterClient(
			conf.EtcdUsername,
			conf.EtcdPassword,
			conf.EtcdAddresses,
			client.Retries(conf.CallRepeaterRetries),
		),

		grpcDownloader: new(utils.Downloader),

		conf:   conf,
		logger: logger,
	}
}

package cli

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	cfcpb "owl/cfc/proto"
	"owl/common/global"
)

// NewCfcClient 创建Cfc客户端
func NewCfcClient(etcdUname, etcdPasswd string, etcdAddrs []string, cliOpt ...client.Option) cfcpb.OwlCfcService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version(global.SrvVersion),
	)

	_ = cli.Client().Init(cliOpt...)

	return cfcpb.NewOwlCfcService(global.OwlCfcRpcRegisterKey, cli.Client())
}

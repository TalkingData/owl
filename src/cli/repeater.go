package cli

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/common/global"
	reppb "owl/repeater/proto"
)

// NewRepeaterClient 创建Repeater客户端
func NewRepeaterClient(etcdUname, etcdPasswd string, etcdAddrs []string, cliOpt ...client.Option) reppb.OwlRepeaterService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version(global.SrvVersion),
	)

	_ = cli.Client().Init(cliOpt...)

	return reppb.NewOwlRepeaterService(global.OwlRepeaterRpcRegisterKey, cli.Client())
}

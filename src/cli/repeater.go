package cli

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/common/global"
	repProto "owl/repeater/proto"
)

// NewRepeaterClient 创建Repeater客户端
func NewRepeaterClient(etcdUname, etcdPasswd string, etcdAddrs []string, cliOpt ...client.Option) repProto.OwlRepeaterService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	_ = cli.Client().Init(cliOpt...)

	return repProto.NewOwlRepeaterService(global.OwlRepeaterServiceName, cli.Client())
}

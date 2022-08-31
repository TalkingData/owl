package cli

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	cfcProto "owl/cfc/proto"
	"owl/common/global"
)

// NewCfcClient 创建Cfc客户端
func NewCfcClient(etcdUname, etcdPasswd string, etcdAddrs []string) cfcProto.OwlCfcService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	return cfcProto.NewOwlCfcService(global.OwlCfcServiceName, cli.Client())
}

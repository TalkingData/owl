package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/common/global"
	"owl/common/logger"
	"owl/repeater/conf"
	reppb "owl/repeater/proto"
	"owl/repeater/service"
)

// Repeater interface
type Repeater interface {
	// Start 启动Repeater
	Start() error
	// Stop 关闭Repeater
	Stop()
}

type defaultRepeater struct {
	srv micro.Service

	conf   *conf.Conf
	logger *logger.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewRepeater(conf *conf.Conf, lg *logger.Logger) Repeater {
	return newDefaultRepeater(conf, lg)
}

func newDefaultRepeater(conf *conf.Conf, lg *logger.Logger) *defaultRepeater {
	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(repConf.EtcdAddresses...),
		etcdv3.Auth(repConf.EtcdUsername, repConf.EtcdPassword),
	)

	srv := micro.NewService(
		micro.Name(repConf.Const.RpcRegisterKey),
		micro.Address(repConf.Listen),
		micro.Version(global.SrvVersion),
		micro.Registry(etcdReg),
		micro.RegisterTTL(repConf.MicroRegisterTtl),
		micro.RegisterInterval(repConf.MicroRegisterInterval),
		micro.Context(ctx),
	)

	_ = reppb.RegisterOwlRepeaterHandler(srv.Server(), service.NewOwlRepeaterService(conf, lg))

	return &defaultRepeater{
		srv: srv,

		conf:   conf,
		logger: lg,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (rep *defaultRepeater) Start() error {
	rep.logger.InfoWithFields(logger.Fields{
		"branch":  global.Branch,
		"commit":  global.Commit,
		"version": global.Version,
	}, "Starting owl repeater...")

	return rep.srv.Run()
}

func (rep *defaultRepeater) Stop() {
	if rep.cancelFunc != nil {
		rep.cancelFunc()
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/cfc/biz"
	"owl/cfc/conf"
	cfcProto "owl/cfc/proto"
	"owl/cfc/service"
	"owl/common/global"
	"owl/common/logger"
	"owl/dao"
)

// Cfc interface
type Cfc interface {
	// Start 启动Cfc
	Start() error
	// Stop 关闭Cfc
	Stop()
}

type defaultCfc struct {
	srv micro.Service
	biz *biz.Biz

	conf   *conf.Conf
	logger *logger.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewCfc(dao *dao.Dao, conf *conf.Conf, lg *logger.Logger) Cfc {
	return newDefaultCfc(dao, conf, lg)
}

func newDefaultCfc(dao *dao.Dao, conf *conf.Conf, lg *logger.Logger) *defaultCfc {
	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(cfcConf.EtcdAddresses...),
		etcdv3.Auth(cfcConf.EtcdUsername, cfcConf.EtcdPassword),
	)

	srv := micro.NewService(
		micro.Name(cfcConf.Const.ServiceName),
		micro.Address(cfcConf.Listen),
		micro.Version("v1"),
		micro.Registry(etcdReg),
		micro.RegisterTTL(cfcConf.MicroRegisterTtl),
		micro.RegisterInterval(cfcConf.MicroRegisterInterval),
		micro.Context(ctx),
	)

	_ = cfcProto.RegisterOwlCfcServiceHandler(srv.Server(), service.NewOwlCfcService(dao, conf, lg))

	return &defaultCfc{
		srv: srv,
		biz: biz.NewBiz(dao, conf, lg),

		conf:   conf,
		logger: lg,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (cfc *defaultCfc) Start() error {
	cfcLg.Info(fmt.Sprintf("Starting owl cfc %s...", global.Version))

	go cfc.biz.RefreshHostStatus(cfc.ctx)
	go cfc.biz.CleanExpiredMetric(cfc.ctx)

	return cfc.srv.Run()
}

func (cfc *defaultCfc) Stop() {
	if cfc.cancelFunc != nil {
		cfc.cancelFunc()
	}
}

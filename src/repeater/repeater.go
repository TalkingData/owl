package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/common/global"
	"owl/common/logger"
	"owl/common/prom"
	"owl/repeater/conf"
	reppb "owl/repeater/proto"
	"owl/repeater/service"
	"sync"
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

	prom prom.Prom

	conf   *conf.Conf
	logger *logger.Logger

	wg         sync.WaitGroup
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

	_prom := prom.NewProm(conf.MetricListen)

	srv := micro.NewService(
		micro.Name(repConf.Const.RpcRegisterKey),
		micro.Address(repConf.Listen),
		micro.Version(global.SrvVersion),
		micro.Registry(etcdReg),
		micro.RegisterTTL(repConf.MicroRegisterTtl),
		micro.RegisterInterval(repConf.MicroRegisterInterval),
		micro.Context(ctx),
		micro.WrapHandler(_prom.GoMicroHandlerWrapper()),
	)

	_ = reppb.RegisterOwlRepeaterHandler(srv.Server(), service.NewOwlRepeaterService(conf, lg))

	return &defaultRepeater{
		srv: srv,

		prom: _prom,

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

	// 启动Prometheus的metrics http server
	go func() {
		rep.wg.Add(1)
		defer rep.Stop()
		defer rep.wg.Done()

		rep.logger.Info(fmt.Sprintf("Owl repeater's metrics http server listening on: %s", rep.conf.MetricListen))
		if err := rep.prom.MetricServerStart(); err != nil {
			rep.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling repeater.prom.ListenAndServe.")
			return
		}
		rep.logger.Info("Owl repeater's metrics server closed.")
	}()

	return rep.srv.Run()
}

func (rep *defaultRepeater) Stop() {
	defer rep.wg.Wait()

	if rep.cancelFunc != nil {
		rep.cancelFunc()
	}

	// 关闭metrics http server
	if rep.prom != nil {
		ctx, cancel := context.WithTimeout(context.Background(), rep.conf.Const.MetricServerShutdownTimeoutSecs)
		defer cancel()
		rep.prom.MetricServerStop(ctx)
	}
}

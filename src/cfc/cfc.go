package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"owl/cfc/biz"
	"owl/cfc/conf"
	cfcpb "owl/cfc/proto"
	"owl/cfc/service"
	"owl/common/global"
	"owl/common/logger"
	"owl/common/prom"
	"owl/dao"
	"sync"
	"time"
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

	prom prom.Prom

	conf   *conf.Conf
	logger *logger.Logger

	wg         sync.WaitGroup
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

	_prom := prom.NewProm(conf.MetricListen)

	srv := micro.NewService(
		micro.Name(cfcConf.Const.RpcRegisterKey),
		micro.Address(cfcConf.Listen),
		micro.Version(global.SrvVersion),
		micro.Registry(etcdReg),
		micro.RegisterTTL(cfcConf.MicroRegisterTtl),
		micro.RegisterInterval(cfcConf.MicroRegisterInterval),
		micro.Context(ctx),
		micro.WrapHandler(_prom.GoMicroHandlerWrapper()),
	)

	_ = cfcpb.RegisterOwlCfcServiceHandler(srv.Server(), service.NewOwlCfcService(dao, conf, lg))

	return &defaultCfc{
		srv: srv,
		biz: biz.NewBiz(dao, conf, lg),

		prom: _prom,

		conf:   conf,
		logger: lg,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (cfc *defaultCfc) Start() error {
	cfc.logger.InfoWithFields(logger.Fields{
		"branch":  global.Branch,
		"commit":  global.Commit,
		"version": global.Version,
	}, "Starting owl cfc...")

	// 启动Prometheus的metrics http server
	go func() {
		cfc.wg.Add(1)
		defer cfc.Stop()
		defer cfc.wg.Done()

		cfc.logger.Info(fmt.Sprintf("Owl cfc's metrics http server listening on: %s", cfc.conf.MetricListen))
		if err := cfc.prom.MetricServerStart(); err != nil {
			cfc.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling cfc.prom.ListenAndServe.")
			return
		}
		cfc.logger.Info("Owl cfc's metrics server closed.")
	}()

	// 启动定时任务
	refreshHostStatusTk := time.Tick(cfc.conf.RefreshHostStatusIntervalSecs)
	if refreshHostStatusTk == nil {
		cfc.logger.Info("biz.RefreshHostStatus not enabled or conf.RefreshHostStatusIntervalSecs is 0.")
	}
	cleanExpiredMetricTk := time.Tick(cfc.conf.CleanExpiredMetricIntervalSecs)
	if cleanExpiredMetricTk == nil {
		cfc.logger.Info("biz.CleanExpiredMetric not enabled or conf.CleanExpiredMetricIntervalSecs is 0.")
	}
	go func() {
		for {
			select {
			case <-refreshHostStatusTk:
				cfc.biz.RefreshHostStatus(cfc.ctx)

			case <-cleanExpiredMetricTk:
				cfc.biz.CleanExpiredMetric(cfc.ctx)

			case <-cfc.ctx.Done():
				cfc.logger.InfoWithFields(logger.Fields{
					"context_error": cfc.ctx.Err(),
				}, "Owl agent exited by context done.")
				return
			}
		}
	}()

	return cfc.srv.Run()
}

func (cfc *defaultCfc) Stop() {
	defer cfc.wg.Wait()

	if cfc.cancelFunc != nil {
		cfc.cancelFunc()
	}

	// 关闭metrics http server
	if cfc.prom != nil {
		ctx, cancel := context.WithTimeout(context.Background(), cfc.conf.Const.MetricServerShutdownTimeoutSecs)
		defer cancel()
		cfc.prom.MetricServerStop(ctx)
	}
}

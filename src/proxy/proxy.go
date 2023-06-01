package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"owl/common/global"
	"owl/common/logger"
	"owl/common/prom"
	"owl/proxy/conf"
	proxypb "owl/proxy/proto"
	"owl/proxy/service"
	"sync"
)

// Proxy interface
type Proxy interface {
	// Start 启动Proxy
	Start() error
	// Stop 关闭Proxy
	Stop()
}

type defaultProxy struct {
	listener   net.Listener
	grpcServer *grpc.Server

	prom prom.Prom

	conf   *conf.Conf
	logger *logger.Logger

	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewProxy(ctx context.Context, conf *conf.Conf, logger *logger.Logger) Proxy {
	pCtx, pCancel := context.WithCancel(ctx)

	return &defaultProxy{
		prom: prom.NewProm(conf.MetricListen),

		conf:   conf,
		logger: logger,

		ctx:        pCtx,
		cancelFunc: pCancel,
	}
}

func (p *defaultProxy) Start() (err error) {
	p.logger.InfoWithFields(logger.Fields{
		"branch":  global.Branch,
		"commit":  global.Commit,
		"version": global.Version,
	}, "Starting owl proxy...")

	// 开启RPC监听端口
	p.listener, err = net.Listen("tcp", p.conf.Listen)
	if err != nil {
		p.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling net.Listen.")
		return err
	}
	listenerAddrStr := p.listener.Addr().String()

	// 创建proxy的grpc server
	p.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(p.prom.GrpcInterceptor))
	// 注册GRPC服务
	proxypb.RegisterOwlProxyServiceServer(
		p.grpcServer,
		service.NewOwlProxyService(p.conf, p.logger),
	)

	p.logger.Info(fmt.Sprintf("Owl proxy listening on: %s", listenerAddrStr))

	// 启动Prometheus的metrics http server
	go func() {
		p.wg.Add(1)
		defer p.Stop()
		defer p.wg.Done()

		p.logger.Info(fmt.Sprintf("Owl proxy's metrics http server listening on: %s", p.conf.MetricListen))
		if err = p.prom.MetricServerStart(); err != nil {
			p.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling proxy.prom.ListenAndServe.")
			return
		}
		p.logger.Info("Owl proxy's metrics server closed.")
	}()

	return p.grpcServer.Serve(p.listener)
}

func (p *defaultProxy) Stop() {
	defer p.wg.Wait()

	// 关闭grpc server
	if p.grpcServer != nil {
		p.grpcServer.Stop()
	}

	// 关闭metrics http server
	if p.prom != nil {
		ctx, cancel := context.WithTimeout(context.Background(), p.conf.Const.MetricServerShutdownTimeoutSecs)
		defer cancel()
		p.prom.MetricServerStop(ctx)
	}

	// 等待所有任务结束
	if p.listener != nil {
		_ = p.listener.Close()
	}
	p.cancelFunc()
}

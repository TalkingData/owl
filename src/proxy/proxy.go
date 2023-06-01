package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"owl/common/global"
	"owl/common/logger"
	"owl/proxy/conf"
	proxypb "owl/proxy/proto"
	"owl/proxy/service"
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

	conf   *conf.Conf
	logger *logger.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewProxy(ctx context.Context, conf *conf.Conf, logger *logger.Logger) Proxy {
	pCtx, pCancel := context.WithCancel(ctx)

	return &defaultProxy{
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
		}, "An error occurred while net.Listen.")
		return err
	}
	listenerAddrStr := p.listener.Addr().String()

	// 创建proxy的grpc server
	p.grpcServer = grpc.NewServer()
	// 注册GRPC服务
	proxypb.RegisterOwlProxyServiceServer(
		p.grpcServer,
		service.NewOwlProxyService(p.conf, p.logger),
	)

	p.logger.Info(fmt.Sprintf("Owl proxy listening on: %s", listenerAddrStr))

	return p.grpcServer.Serve(p.listener)
}

func (p *defaultProxy) Stop() {
	defer p.logger.Info("Owl proxy Stopped.")

	// 关闭grpc server
	if p.grpcServer != nil {
		p.grpcServer.Stop()
	}

	// 等待所有任务结束
	if p.listener != nil {
		_ = p.listener.Close()
	}
	p.cancelFunc()
}

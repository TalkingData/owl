package component

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"owl/cfc/biz"
	"owl/cfc/conf"
	cfcProto "owl/cfc/proto"
	"owl/cfc/service"
	"owl/common/global"
	"owl/common/logger"
	"owl/dao"
)

// cfc struct
type cfc struct {
	dao *dao.Dao

	biz *biz.Biz

	grpcServer *grpc.Server
	listener   net.Listener

	conf       *conf.Conf
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *logger.Logger
}

func newCfc(ctx context.Context, dao *dao.Dao, conf *conf.Conf, lg *logger.Logger) *cfc {
	c := &cfc{
		dao: dao,

		biz: biz.NewBiz(dao, conf, lg),

		conf:   conf,
		logger: lg,
	}

	c.ctx, c.cancelFunc = context.WithCancel(ctx)

	return c
}

func (cfc *cfc) Start() (err error) {
	cfc.logger.Info(fmt.Sprintf("Starting owl cfc %s...", global.Version))

	// 开启RPC监听端口
	cfc.listener, err = net.Listen("tcp", cfc.conf.Listen)
	if err != nil {
		cfc.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while net.Listen.")
		return err
	}

	cfc.grpcServer = grpc.NewServer()
	cfcProto.RegisterOwlCfcServiceServer(
		cfc.grpcServer,
		service.NewOwlCfcService(cfc.biz, cfc.dao, cfc.conf, cfc.logger),
	)

	cfc.logger.Info(fmt.Sprintf("Owl cfc listening on: %s", cfc.listener.Addr().String()))

	// 启动一些常驻后台的操作
	go cfc.biz.RefreshHostStatus(cfc.ctx)
	go cfc.biz.CleanExpiredMetric(cfc.ctx)

	return cfc.grpcServer.Serve(cfc.listener)
}

func (cfc *cfc) Stop() {
	defer cfc.logger.Info("Owl cfc stopped.")

	// 关闭grpc server
	if cfc.grpcServer != nil {
		cfc.grpcServer.Stop()
	}

	// 等待所有任务结束
	if cfc.listener != nil {
		_ = cfc.listener.Close()
	}

	cfc.cancelFunc()
}

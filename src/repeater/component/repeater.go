package component

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"owl/common/global"
	"owl/common/logger"
	"owl/repeater/backend"
	"owl/repeater/conf"
	repProto "owl/repeater/proto"
	"owl/repeater/service"
)

// repeater struct
type repeater struct {
	grpcServer *grpc.Server
	listener   net.Listener
	backend    backend.Backend

	conf *conf.Conf

	logger *logger.Logger
}

func newRepeater(conf *conf.Conf, lg *logger.Logger) *repeater {
	bk, err := backend.NewBackend(conf)
	if err != nil {
		lg.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while newDefaultRepeater.")
		panic(err)
	}

	return &repeater{
		backend: bk,
		conf:    conf,
		logger:  lg,
	}
}

func (rep *repeater) Start() (err error) {
	rep.logger.Info(fmt.Sprintf("Starting owl repeater %s", global.Version))

	// 开启RPC监听端口
	rep.listener, err = net.Listen("tcp", rep.conf.Listen)
	if err != nil {
		rep.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while net.Listen.")
		return err
	}

	rep.grpcServer = grpc.NewServer()
	repProto.RegisterOwlRepeaterServiceServer(
		rep.grpcServer,
		service.NewOwlRepeaterService(rep.backend, rep.logger),
	)

	rep.logger.Info(fmt.Sprintf("Owl repeater starting at: %s", rep.listener.Addr().String()))

	return rep.grpcServer.Serve(rep.listener)
}

func (rep *repeater) Stop() {
	defer rep.logger.Info("Owl repeater stopped.")

	// 关闭grpc server
	if rep.grpcServer != nil {
		rep.grpcServer.Stop()
	}

	// 等待所有任务结束
	if rep.listener != nil {
		_ = rep.listener.Close()
	}
}

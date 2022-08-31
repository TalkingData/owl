package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"os"
	"os/signal"
	"owl/common/logger"
	"owl/repeater/conf"
	repProto "owl/repeater/proto"
	"owl/repeater/service"
	"runtime"
	"syscall"
)

var (
	repSrv micro.Service

	repConf *conf.Conf
	repLg   *logger.Logger
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(repConf.EtcdAddresses...),
		etcdv3.Auth(repConf.EtcdUsername, repConf.EtcdPassword),
	)

	repSrv = micro.NewService(
		micro.Name(repConf.Const.ServiceName),
		micro.Address(repConf.Listen),
		micro.Version("v1"),
		micro.Registry(etcdReg),
		micro.RegisterTTL(repConf.MicroRegisterTtl),
		micro.RegisterInterval(repConf.MicroRegisterInterval),
		micro.Context(ctx),
	)

	_ = repProto.RegisterOwlRepeaterHandler(repSrv.Server(), service.NewOwlRepeaterService(repConf, repLg))

	e := make(chan error)
	go func() {
		e <- repSrv.Run()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				repLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while repSrv.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			repLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			cancel()
			return
		}
	}
}

// closeAll
func closeAll() {
	if repLg != nil {
		repLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	repConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(repConf.LogLevel),
		logger.LogPath(repConf.LogPath),
		logger.ServiceName(repConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	repLg = lg
}

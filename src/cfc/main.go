package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"os"
	"os/signal"
	"owl/cfc/conf"
	cfcProto "owl/cfc/proto"
	"owl/cfc/service"
	"owl/common/logger"
	"owl/common/orm"
	"owl/dao"
	"runtime"
	"syscall"
)

var (
	cfcSrv micro.Service

	cfcDao  *dao.Dao
	cfcConf *conf.Conf
	cfcLg   *logger.Logger
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(cfcConf.EtcdAddresses...),
		etcdv3.Auth(cfcConf.EtcdUsername, cfcConf.EtcdPassword),
	)

	cfcSrv = micro.NewService(
		micro.Name(cfcConf.Const.ServiceName),
		micro.Address(cfcConf.Listen),
		micro.Version("v1"),
		micro.Registry(etcdReg),
		micro.RegisterTTL(cfcConf.MicroRegisterTtl),
		micro.RegisterInterval(cfcConf.MicroRegisterInterval),
		micro.Context(ctx),
	)

	_ = cfcProto.RegisterOwlCfcServiceHandler(cfcSrv.Server(), service.NewOwlCfcService(cfcDao, cfcConf, cfcLg))

	e := make(chan error)
	go func() {
		e <- cfcSrv.Run()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				cfcLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while cfcSrv.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			cfcLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			cancel()
			return
		}
	}
}

// closeAll
func closeAll() {
	if cfcDao != nil {
		cfcDao.Close()
	}
	if cfcLg != nil {
		cfcLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	cfcConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(cfcConf.LogLevel),
		logger.LogPath(cfcConf.LogPath),
		logger.ServiceName(cfcConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	cfcLg = lg

	cfcDao = dao.NewDao(orm.NewMysqlGorm(
		cfcConf.MysqlAddress,
		cfcConf.MysqlUser,
		cfcConf.MysqlPassword,
		cfcConf.MysqlDbName,
		orm.MysqlMaxIdleConns(cfcConf.MysqlMaxIdleConns),
		orm.MysqlMaxOpenConns(cfcConf.MysqlMaxOpenConns),
	), cfcLg)
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"owl/cfc/component"
	"owl/cfc/conf"
	"owl/common/logger"
	"owl/common/orm"
	"owl/dao"
	"runtime"
	"syscall"
)

var (
	cfc component.Component

	cfcDao  *dao.Dao
	cfcConf *conf.Conf
	cfcLg   *logger.Logger
)

func main() {
	cfc = component.NewCfcComponent(context.Background(), cfcDao, cfcConf, cfcLg)
	if cfc == nil {
		cfcLg.ErrorWithFields(logger.Fields{
			"error": fmt.Errorf("nil cfc error"),
		}, "An error occurred while main.")
		return
	}

	e := make(chan error)
	go func() {
		e <- cfc.Start()
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
				}, "An error occurred while cfc.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			cfcLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if cfc != nil {
		cfc.Stop()
	}
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

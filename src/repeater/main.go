package main

import (
	"fmt"
	"os"
	"os/signal"
	"owl/common/logger"
	"owl/repeater/conf"
	"runtime"
	"syscall"
)

var (
	rep Repeater

	repConf *conf.Conf
	repLg   *logger.Logger
)

func main() {
	rep = NewRepeater(repConf, repLg)

	e := make(chan error)
	go func() {
		e <- rep.Start()
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
				}, "An error occurred while calling repSrv.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			repLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if rep != nil {
		rep.Stop()
	}
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
		logger.Filename(repConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while calling logger.NewLogger, error:", err.Error())
		panic(err)
	}
	repLg = lg
}

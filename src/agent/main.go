package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"owl/agent/conf"
	"owl/common/logger"
	"runtime"
	"syscall"
)

var (
	agt Agent

	agtConf *conf.Conf
	agtLg   *logger.Logger
)

func main() {
	var err error

	agt, err = NewAgent(context.Background(), agtConf, agtLg)
	if err != nil {
		agtLg.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while main.")
		return
	}
	if agt == nil {
		agtLg.ErrorWithFields(logger.Fields{
			"error": fmt.Errorf("nil agent error"),
		}, "An error occurred while main.")
		return
	}

	e := make(chan error)
	go func() {
		e <- agt.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err = <-e:
			if err != nil {
				agtLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while agent.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			agtLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

func closeAll() {
	if agt != nil {
		agt.Stop()
	}
	if agtLg != nil {
		agtLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	agtConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(agtConf.LogLevel),
		logger.LogPath(agtConf.LogPath),
		logger.Filename(agtConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	agtLg = lg
}

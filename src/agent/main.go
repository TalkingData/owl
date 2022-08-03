package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"owl/agent/component"
	"owl/agent/conf"
	"owl/common/logger"
	"runtime"
	"syscall"
)

var (
	agent component.Component

	agentConf *conf.Conf
	agentLg   *logger.Logger
)

func main() {
	var err error

	agent, err = component.NewAgentComponent(context.Background(), agentConf, agentLg)
	if err != nil {
		agentLg.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while main.")
		return
	}
	if agent == nil {
		agentLg.ErrorWithFields(logger.Fields{
			"error": fmt.Errorf("nil agent error"),
		}, "An error occurred while main.")
		return
	}

	e := make(chan error)
	go func() {
		e <- agent.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err = <-e:
			if err != nil {
				agentLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while agent.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			agentLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

func closeAll() {
	if agent != nil {
		agent.Stop()
	}
	if agentLg != nil {
		agentLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	agentConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(agentConf.LogLevel),
		logger.LogPath(agentConf.LogPath),
		logger.ServiceName(agentConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	agentLg = lg
}

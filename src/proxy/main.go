package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"owl/common/logger"
	"owl/proxy/conf"
	"runtime"
	"syscall"
)

var (
	proxy Proxy

	proxyConf *conf.Conf
	proxyLg   *logger.Logger
)

func main() {
	proxy = NewProxy(context.Background(), proxyConf, proxyLg)
	if proxy == nil {
		proxyLg.ErrorWithFields(logger.Fields{
			"error": fmt.Errorf("nil proxy error"),
		}, "An error occurred while main.")
		return
	}

	e := make(chan error)
	go func() {
		e <- proxy.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				proxyLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while proxy.Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			proxyLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if proxy != nil {
		proxy.Stop()
	}
	if proxyLg != nil {
		proxyLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	proxyConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(proxyConf.LogLevel),
		logger.LogPath(proxyConf.LogPath),
		logger.Filename(proxyConf.Const.ServiceName),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	proxyLg = lg
}

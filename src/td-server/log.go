package main

//服务端日志
//数据日志
//网络设备日志
//端口扫描日志
//数据库操作日志

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var (
	slog, nlog, dlog, sqllog *logs.BeeLogger
)

func init_log(cfg *Config) error {
	//数据接收日志
	dlog = logs.NewLogger(100)
	dlog.EnableFuncCallDepth(true)
	dlog.SetLevel(cfg.LOG_LEVEL)
	if err := dlog.SetLogger("file", fmt.Sprintf(`{"filename":"%s/data.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}
	//网络设备日志
	nlog = logs.NewLogger(100)
	nlog.EnableFuncCallDepth(true)
	nlog.SetLevel(cfg.LOG_LEVEL)
	if err := nlog.SetLogger("file", fmt.Sprintf(`{"filename":"%s/netdevice.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}

	slog = logs.NewLogger(100)
	slog.EnableFuncCallDepth(true)
	slog.SetLevel(cfg.LOG_LEVEL)
	if err := slog.SetLogger("file", fmt.Sprintf(`{"filename":"%s/server.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}
	return nil
	//客户端日志
	//
}

package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var (
	log  *logs.BeeLogger
	dlog *logs.BeeLogger
)

func init_log(cfg *Config) error {
	//数据接收日志
	log = logs.NewLogger(100)
	log.EnableFuncCallDepth(true)
	log.SetLevel(cfg.LOG_LEVEL)
	if err := log.SetLogger("file", fmt.Sprintf(`{"filename":"%s/proxy.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}
	dlog = logs.NewLogger(100)
	dlog.EnableFuncCallDepth(true)
	dlog.SetLevel(cfg.LOG_LEVEL)
	if err := dlog.SetLogger("file", fmt.Sprintf(`{"filename":"%s/data.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}

	return nil
}

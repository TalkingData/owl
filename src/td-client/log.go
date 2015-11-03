package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var (
	log *logs.BeeLogger
)

func init_log() error {
	log = logs.NewLogger(100)
	log.EnableFuncCallDepth(true)
	log.SetLevel(cfg.LOG_LEVEL)
	if err := log.SetLogger("file", fmt.Sprintf(`{"filename":"%s/client.log","rotate":true,"maxdays":%d}`,
		cfg.LOG_DIR, cfg.LOG_EXPIRE_DAYS)); err != nil {
		return err
	}
	return nil
}

package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var (
	lg *logs.BeeLogger
)

func InitLog() error {
	lg = logs.NewLogger(100)
	lg.EnableFuncCallDepth(true)
	lg.SetLogger("console", "")
	param := fmt.Sprintf(`{"filename":"%s","rotate":true,"maxdays":%d}`,
		GlobalConfig.LogFile, GlobalConfig.LogExpireDays)
	if err := lg.SetLogger("file", param); err != nil {
		return err
	}
	lg.SetLevel(GlobalConfig.LogLevel)
	return nil
}

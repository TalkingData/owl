package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	os.Chdir(filepath.Dir(os.Args[0]))
	os.Mkdir("logs", 0755)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if err := InitGlobalConfig(); err != nil {
		fmt.Println("failed to init global config:", err)
		return
	}
	if err := InitMysqlConnPool(); err != nil {
		fmt.Println("failed to init mysql connection pool:", err)
		return
	}
	if GlobalConfig.AUTO_BUILD_METRIC_TAG_INDEX {
		autoBuildMetricAndTagIndex()
	}
	if err := InitServer(); err != nil {
		fmt.Println("failed to init http server:", err)
		return
	}
}

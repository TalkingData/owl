package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
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
	var err error
	if err = InitGlobalConfig(); err != nil {
		fmt.Println(err)
		return
	}
	if err = InitLog(); err != nil {
		fmt.Println("InitLog ", err)
		return
	}
	go func() {
		fmt.Printf("start metric interface %s\n", GlobalConfig.MetricBind)
		fmt.Printf("%s\n", http.ListenAndServe(GlobalConfig.MetricBind, nil))
	}()
	if err = InitMysqlConnPool(); err != nil {
		fmt.Println("init mysql error: ", err.Error())
		return
	}
	if err = InitCFC(); err != nil {
		fmt.Println(err)
		return
	}
	lg.Info("start cfc on %s ", GlobalConfig.TCPBind)
	go updatHostStatus()
	if GlobalConfig.EnableCleanupExpiredMetric {
		go cleanupExpiredMetrics()
	}
	select {}
}

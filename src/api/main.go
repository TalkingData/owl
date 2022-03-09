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
		fmt.Println("failed to init global config:", err)
		return
	}
	if err = initPublicKeyAndPrivateKey(); err != nil {
		panic(err)
	}
	go func() {
		fmt.Printf("start metric interface %s\n", config.MetricBind)
		fmt.Printf("%s\n", http.ListenAndServe(config.MetricBind, nil))
	}()
	if err = initTSDB(); err != nil {
		fmt.Println("failed to init tsdb storage:", err)
		return
	}
	if err = InitMysqlConnPool(); err != nil {
		fmt.Println("failed to init mysql connection pool:", err)
		return
	}
	if err := InitServer(); err != nil {
		fmt.Println("failed to init http server:", err)
		return
	}
}

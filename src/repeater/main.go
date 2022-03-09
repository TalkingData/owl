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
		fmt.Println("failed to init log.")
		return
	}
	go func() {
		fmt.Printf("start metric interface %s\n", GlobalConfig.MetricBind)
		fmt.Printf("%s\n", http.ListenAndServe(GlobalConfig.MetricBind, nil))
	}()
	if err = InitRepeater(); err != nil {
		lg.Error("init repeater error, %s", err)
		return
	}
	go repeater.Forward()
	select {}
}

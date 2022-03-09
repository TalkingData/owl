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
		fmt.Println("failed to init logfile.")
		return
	}
	go func() {
		fmt.Printf("start metric interface %s\n", GlobalConfig.MetricBind)
		fmt.Printf("%s\n", http.ListenAndServe(GlobalConfig.MetricBind, nil))
	}()
	if err = InitAgent(); err != nil {
		fmt.Println(err)
		return
	}

	lg.Info("listen on %s", GlobalConfig.TCPBind)
	agent.register()
	go agent.watchConnLoop()
	go agent.watchHostConfig()
	go agent.StartTimer()

	go agent.SendTSD2Repeater()
	go agent.sendHeartbeat2Repeater()
	go agent.syncMetricToCFC()
	select {}
}

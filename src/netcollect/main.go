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
		fmt.Printf("start pprof failed %s\n", http.ListenAndServe(":10071", nil))
	}()
	if err = InitNetCollect(); err != nil {
		fmt.Println(err)
		return
	}

	lg.Debug("dial cfc %v", netCollect.dialCFC())
	lg.Debug("dial repeater %v", netCollect.dialRepeater())
	go netCollect.watchConnLoop()

	if err = InitIpRange(); err != nil {
		fmt.Println(err)
		return
	}

	go netCollect.SendTSD2Repeater()
	select {}
}

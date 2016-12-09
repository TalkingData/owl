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
	var err error
	if err = InitGlobalConfig(); err != nil {
		fmt.Println(err)
		return
	}
	if err = InitLog(); err != nil {
		fmt.Println("failed to init log.")
		return
	}

	if err = InitNetCollect(); err != nil {
		fmt.Println(err)
		return
	}

	go netCollect.Dial("cfc")
	go netCollect.Dial("repeater")

	if err = InitIpRange(); err != nil {
		fmt.Println(err)
		return
	}

	go netCollect.SendTSD2Repeater()
	select {}
}

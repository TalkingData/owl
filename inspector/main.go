/*
author: chao.ma
*/
package main

import (
	"fmt"
	"os"
	chm "owl/common/chanMonitor"
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
	if err := InitLog(); err != nil {
		fmt.Println("failed to init log:", err)
		return
	}
	if err := InitTsdb(); err != nil {
		fmt.Println("failed to init tsdb:", err)
		return
	}
	if err := InitInspector(); err != nil {
		fmt.Println("failed to init inspector:", err)
		return
	}

	chm.AddNamed("inspector.resultPool.results", "owl-inspector", inspector.resultPool.results)
	chm.AddNamed("inspector.taskPool.tasks", "owl-inspector", inspector.taskPool.tasks)
	chm.New("owl-inspector", ":20001").Start()
	select {}
}

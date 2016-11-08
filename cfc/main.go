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

	if err = InitMysqlConnPool(); err != nil {
		lg.Error("init mysql error: %s", err.Error())
		return
	}

	if err = InitCfc(); err != nil {
		fmt.Println(err)
		return
	}
	UpdatHostAive()
	select {}
}

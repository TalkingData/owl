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
		fmt.Println("failed to init logfile.")
		return
	}

	if err = InitAgent(); err != nil {
		fmt.Println(err)
		return
	}
	lg.Info("listen on %s", GlobalConfig.TCP_BIND)
	go agent.Dial("cfc")
	go agent.Dial("repeater")
	agent.SendConfig2CFC()
	agent.GetPluginList()
	agent.SendTSD2Repeater()
	agent.RunBuiltinMetric()
	agent.SendAgentAlive2Repeater()
	agent.SendHostAlive2CFC()
	select {}
}

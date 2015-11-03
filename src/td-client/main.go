package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

var (
	host        *Host = &Host{}
	DataBuffer  chan []byte
	version     = "2.0"
	cfg         *Config
	MetricCache = make(map[string][2]float64)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	file, _ := exec.LookPath(os.Args[0])
	pth, _ := path.Split(file)
	os.Chdir(pth)
	//加载配置文件
	var err error
	cfg, err = load_config("./conf/client.conf")
	if err != nil {
		fmt.Printf("load config file error , %v\n", err)
		os.Exit(1)
	}
	//创建目录
	for _, dir := range []string{cfg.LOG_DIR, "var", "update", "plugins"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("create dir(%s) error %v", dir, err)
			os.Exit(1)
		}
	}
	fd, err := os.Create("var/tdclient.pid")
	if err != nil {
		log.Error("create pid file error <%v> ", err)
		return
	}
	fd.WriteString(strconv.Itoa(os.Getpid()))
	fd.Close()
	//初始化日志文件
	if err := init_log(); err != nil {
		fmt.Println("main.init_log error %s", err)
		os.Exit(1)
	}
	err = host.LoadFromFile()
	if err != nil {
		log.Error("load host from file error(%v)", err)
		host.UUID = NewUUID()
		log.Info("init new uuid %s", host.UUID)
		host.SaveToFile()
		log.Info("save to file cache.")
	}
	host.Ip = ""
	DataBuffer = make(chan []byte, cfg.BUFFER_SIZE)

	server := NewTCPServe()
	server.SetPacketLimitSize(1024 * 1024)
	go server.Start()

	go host.Loop(server)
	go host.Monitor(DataBuffer)
	go GuardHandle()
	host.ServerHB()
}

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
	DataBuffer chan []byte
	cfg        *Config
	version    string
	hostmap    map[string]*Host = make(map[string]*Host)

	devicemap map[int]*NetDevice
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	file, _ := exec.LookPath(os.Args[0])
	pth, _ := path.Split(file)
	os.Chdir(pth)
	//加载配置文件
	var err error
	cfg, err := load_config("./conf/proxy.conf")
	if err != nil {
		fmt.Printf("load config file error , %v\n", err)
		os.Exit(1)
	}

	//创建目录
	for _, dir := range []string{cfg.LOG_DIR, "var", "update"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("create dir(%s) error %v", dir, err)
			os.Exit(1)
		}
	}

	//初始化日志文件
	if err := init_log(cfg); err != nil {
		fmt.Println("main.init_log error %s", err)
		os.Exit(1)
	}

	fd, err := os.Create("var/td-proxy.pid")
	if err != nil {
		log.Error("create pid file error <%v> ", err)
		return
	}

	fd.WriteString(strconv.Itoa(os.Getpid()))
	fd.Close()

	DataBuffer = make(chan []byte, cfg.BUFFER_SIZE)
	devicemap = make(map[int]*NetDevice)

	proxy := NewProxy(cfg)
	proxy.SetPacketLimitSize(uint32(cfg.MAX_PACKET_SIZE))
	go proxy.Start()

	go proxy.HandleLoop()
	go proxy.ForwardingData(cfg)

	select {}
}

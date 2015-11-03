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
	DataBuffer  chan []byte
	RedisBuffer chan string
	cfg         *Config
	mysql       *db
	devicemap   map[int]*NetDevice
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	file, _ := exec.LookPath(os.Args[0])
	pth, _ := path.Split(file)
	os.Chdir(pth)

	//加载配置文件
	var err error
	if cfg, err = load_config("./conf/server.conf"); err != nil {
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
	//写入pid
	fd, err := os.Create("var/td-server.pid")
	if err != nil {
		slog.Error("create pid file error <%v> ", err)
		return
	}

	fd.WriteString(strconv.Itoa(os.Getpid()))
	fd.Close()

	DataBuffer = make(chan []byte, cfg.BUFFER_SIZE)
	RedisBuffer = make(chan string, cfg.BUFFER_SIZE)
	devicemap = make(map[int]*NetDevice)

	//初始化数据库链接
	slog.Info("connect to database ...")
	mysql, err = NewMysqlConnPool(cfg)
	if err != nil {
		slog.Error("connect database error, %v\n", err)
		os.Exit(1)
	} else {
		slog.Info("connect database ok..")
	}

	//加载网络设备信息，并采集数据
	devices, err := mysql.GetDeviceByProxy("")
	for _, dev := range devices {
		go dev.Run()
		devicemap[dev.ID] = dev
	}
	//维护网络设备列表 deviceMap
	go HandleDevicesLoop(devicemap)

	//启动tcp server
	go StartTCPServe(cfg)
	go StartHttpServe(cfg)
	go PortCheck()
	go HostLoop()
	if cfg.ENABLE_REDIS {
		go WriteToRedis(RedisBuffer, cfg)
	}
	WiteToTSDB(DataBuffer, cfg)

}

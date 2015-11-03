package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"tcpserver"
	"time"
)

func main() {
	//set workpath
	file, _ := exec.LookPath(os.Args[0])
	pth, _ := path.Split(file)
	os.Chdir(pth)

	cfg, err := load_config("./conf/client.conf")
	if err != nil {
		fmt.Printf("load config file error , %v\n", err)
		os.Exit(1)
	}

	for _, dir := range []string{cfg.LOG_DIR, "var", "update", "plugins"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("create dir(%s) error %v", dir, err)
			os.Exit(1)
		}
	}

	if err := init_log(cfg); err != nil {
		fmt.Println("main.init_log error %s", err)
		os.Exit(1)
	}

	fd, err := os.Create("var/guard.pid")
	if err != nil {
		log.Error("create pid file error <%v> ", err)
		return
	}
	fd.WriteString(strconv.Itoa(os.Getpid()))
	fd.Close()
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.GUARDBIND)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	server := tcpserver.NewServer(listener, &handle{})
	go server.Start()
START:

	conn, err := net.DialTimeout("tcp", cfg.TCPBIND, time.Second*5)
	if err != nil {
		log.Warn("connect to guard error %s", err.Error())
		StartClient()
		time.Sleep(time.Second * 10)
		goto START
	}
	for {
		var buf bytes.Buffer
		head := make([]byte, 4)
		binary.BigEndian.PutUint32(head, uint32(1))
		binary.Write(&buf, binary.BigEndian, head)
		binary.Write(&buf, binary.BigEndian, byte(GUARDHB))
		lenght, err := conn.Write(buf.Bytes())
		if err != nil {
			log.Warn("send packet to guard error %s", err.Error())
			StartClient()
			time.Sleep(time.Second * 10)
			goto START
		}
		log.Info("send hb packet to client doen. %d bytes", lenght)
		time.Sleep(time.Second * 1)
	}

}

func StartClient() {
	cmd := exec.Command("sh", "-c", "nohup ./td-client >/dev/null 2>/dev/null &")
	err := cmd.Run()
	if err != nil {
		log.Error("start client error %s", err.Error())
	}
	log.Info("start client done. pit(%d)", cmd.ProcessState.Pid())

}

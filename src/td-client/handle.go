package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"tcpserver"
	"time"
)

type handle struct {
}

//数据包逻辑处理
func (this *handle) HandlerMessage(conn *tcpserver.Conn, data []byte) {
	defer func() {
		recover()
	}()
	log.Debug("receive data: %s|%s", data[0], string(data[1:]))
	switch data[0] {
	case HOSTCONFIGRESP: //get host config
		if err := json.Unmarshal(data[1:], &host); err != nil {
			log.Error("unmarshal data to host error(%s) data(%s)", err.Error(), string(data[1:]))
			return
		}
		host.SaveToFile()
	case CLIENTVERSIONRESP:
		var vp map[string]string
		if err := json.Unmarshal(data[1:], &vp); err != nil {
			log.Error("unmarshal error %s", err.Error())
			return
		}
		if v, ok := vp["version"]; ok {
			filename := "td-client-" + vp["version"] + ".tar.gz"
			if v != version {
				_, err := os.Stat("./update/" + filename)
				if err == nil {
					return
				}
				url := "http://" + cfg.HTTPSERVER + "/" + filename
				if err := DownloadFile(url); err != nil {
					log.Error("update client version error :%s", err.Error())
					return
				}
				if err := Unzip(filename); err != nil {
					log.Error(err.Error())
					return
				}
				log.Info("client update done, version:%s,  wait exit.", vp["version"])
				time.Sleep(time.Second * 5)
				os.Exit(0)
			}
		}
	case GUARDHB:
		//log.Trace("receive guard heatbeart packet")
	default:
		log.Critical(string(data[1:]))
	}
}

func (this *handle) Connect(conn *tcpserver.Conn) {
	log.Info("%s connected", conn.GetLocalIp())
}

func (this *handle) Disconnect(conn *tcpserver.Conn) {
	log.Info("%s disconnect ", conn.GetLocalIp())
}

func StartGuard() {
	cmd := exec.Command("sh", "-c", "nohup ./td-guard >/dev/null 2>/dev/null &")
	err := cmd.Run()
	if err != nil {
		log.Error("start guard error %s", err.Error())
	}
	cmd.Output()
	log.Info("start guard done. %s", cmd.ProcessState.String())
}

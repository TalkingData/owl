package main

import (
	"encoding/json"
	"net"
	"tcpserver"
	"time"
)

type Proxy struct {
	*tcpserver.Server
	Ip string
}

func NewProxy(cfg *Config) *Proxy {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.TCPBIND)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		panic(err)
	}
	return &Proxy{tcpserver.NewServer(listener, &handle{}), ""}
}

func (proxy *Proxy) HandleLoop() {

	for {
		if proxy.Ip == "" {
			time.Sleep(time.Second * 30)
			continue
		}
		req := ProxyReq{Proxy: proxy.Ip}
		select {
		case <-time.Tick(time.Minute * 5):
			DataBuffer <- NewPacket(CLIENTVERSION, []byte(""))
			DataBuffer <- NewPacket(GETDEVICES, req.ToJson())
			DataBuffer <- NewPacket(GETPORTS, req.ToJson())
			devices := []*NetDevice{}
			for _, dev := range devicemap {
				if len(dev.DeviceInterfaces) > 0 {
					devices = append(devices, dev)
				}
			}
			req, _ := json.Marshal(devices)
			DataBuffer <- NewPacket(UPDATEDEVICES, req)
		}
	}

}

func (proxy *Proxy) ForwardingData(cfg *Config) {
START:
	conn, err := proxy.Connect(cfg.TCPSERVER)
	if err != nil {
		log.Error("连接服务端失败, 地址:%s 错误信息:%s", cfg.TCPSERVER, err.Error())
		time.Sleep(time.Second * 10)
		goto START
	}
	go conn.Run()
	proxy.Ip = conn.GetLocalIp()

	for {
		select {
		case msg, ok := <-DataBuffer:
			if ok {
				if err := conn.AsyncWriteData(msg); err != nil {
					log.Critical("write data to server error message(%v)", err)
					DataBuffer <- msg
					goto START
				}
				dlog.Info("send data: %s|%s", PROTOCOLTYPE[msg[4]], string(msg[5:]))
			}
		}
	}
}

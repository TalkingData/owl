package main

import (
	"owl/common/types"
	"sync"
	"time"

	"github.com/wuyingsong/tcp"
)

var (
	proxy *CFCProxy
)

type CFCProxy struct {
	srv *tcp.AsyncTCPServer
	cfc *tcp.TCPConn

	connMap map[string]string //hostid -> tcpaddr1
	lock    sync.RWMutex
}

func InitCfcProxy() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MaxPacketSize))
	s := tcp.NewAsyncTCPServer(GlobalConfig.TCPBind, &callback{}, protocol)
	proxy = &CFCProxy{
		s,
		&tcp.TCPConn{},
		make(map[string]string),
		sync.RWMutex{},
	}
	return proxy.srv.ListenAndServe()
}

// 发送插件同步请求
func (proxy *CFCProxy) sendSyncPluginRequest(hostID string, p types.Plugin) error {
	lg.Debug("send sync plugin request, %s", p.Path)
	spr := types.SyncPluginRequest{
		hostID,
		p,
	}
	return proxy.cfc.AsyncWritePacket(
		tcp.NewDefaultPacket(types.MsgAgentRequestSyncPlugins, spr.Encode()),
	)
}

func (proxy *CFCProxy) DialCFC() {
	var (
		tempDelay time.Duration
		err       error
		cfc       *tcp.TCPConn
		reconnect bool
	)

retry:
	if reconnect {
		lg.Error("cfc session is closed, retry.")
	}
	cfc, err = proxy.srv.Connect(GlobalConfig.CFCAddr, nil, nil)
	if err != nil {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 5 * time.Second; tempDelay > max {
			tempDelay = max
		}
		time.Sleep(tempDelay)
		reconnect = true
		goto retry
	}
	lg.Info("connect cfc %s successfully.", GlobalConfig.CFCAddr)
	proxy.cfc = cfc
	for {
		if proxy.cfc.IsClosed() {
			reconnect = true
			goto retry
		}
		time.Sleep(time.Second)
	}
}

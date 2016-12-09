package main

import (
	"owl/common/tcp"
	"owl/common/types"
	"sync"
	"time"
)

var (
	proxy *CfcProxy
)

type CfcProxy struct {
	srv        *tcp.Server
	cfc        *tcp.Session
	PluginList HostPluginList
}

//存放主机需要执行的插件列表
type HostPluginList struct {
	M    map[string][]types.Plugin
	lock sync.RWMutex
}

func InitCfcProxy() error {
	s := tcp.NewServer(GlobalConfig.TCP_BIND, &handle{})
	proxy = &CfcProxy{
		s,
		&tcp.Session{},
		HostPluginList{M: make(map[string][]types.Plugin)},
	}
	return proxy.srv.ListenAndServe()
}

func (this *CfcProxy) DialCFC() {
	var (
		tempDelay time.Duration
		err       error
		cfc       *tcp.Session
		reconnect bool
	)

retry:
	if reconnect {
		lg.Error("cfc session is closed, retry.")
	}
	cfc, err = this.srv.Connect(GlobalConfig.CFC_ADDR, nil)
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
	lg.Info("connect cfc %s successfully.", GlobalConfig.CFC_ADDR)
	this.cfc = cfc
	for {
		if this.cfc.IsClosed() {
			reconnect = true
			goto retry
		}
		time.Sleep(time.Second)
	}
}

func (this *HostPluginList) Get(host_id string) []types.Plugin {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if l, ok := this.M[host_id]; ok {
		return l
	}
	return nil
}

func (this *HostPluginList) Set(host_id string, list []types.Plugin) {
	this.lock.Lock()
	this.M[host_id] = list
	this.lock.Unlock()
}

func (this *HostPluginList) All() map[string][]types.Plugin {
	m := make(map[string][]types.Plugin, len(this.M))
	this.lock.RLock()
	m = this.M
	this.lock.RUnlock()
	return m
}

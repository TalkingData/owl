package main

import (
	"time"

	"github.com/wuyingsong/tcp"
)

const (
	timeFomart = "2006-01-02 15:04:05"
)

type CFC struct {
	*tcp.AsyncTCPServer
}

var cfc = new(CFC)

// InitCFC Initialize the cfc server
func InitCFC() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MaxPacketSize))
	cfc.AsyncTCPServer = tcp.NewAsyncTCPServer(GlobalConfig.TCPBind, &callback{}, protocol)
	return cfc.ListenAndServe()
}

func updatHostStatus() {
	for {
		for _, host := range mydb.getAllHosts() {
			timeDiff := time.Now().Sub(host.UpdateAt).Seconds()
			if timeDiff > 120 {
				if host.IsAlive() {
					lg.Info("set host(%s %s) status down, time difference:%0.2fs", host.IP, host.Hostname, timeDiff)
					mydb.setHostAlive(host.ID, "0")
				}
			} else {
				if !host.IsAlive() {
					lg.Info("set host(%s %s) status ok, time difference:%0.2fs ", host.IP, host.Hostname, timeDiff)
					mydb.setHostAlive(host.ID, "1")
				}
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func cleanupExpiredMetrics() {
	for {
		time.Sleep(time.Minute * time.Duration(GlobalConfig.CleanupExpiredMetricIntervalMinutes))
		if err := mydb.cleanupExpiredMetrics(); err != nil {
			lg.Error("cleanupExpiredMetrics failed, error:%s", err)
		}
	}
}

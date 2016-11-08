package main

import (
	"owl/common/tcp"
	"time"
)

var cfcServer *tcp.Server

func InitCfc() error {
	cfcServer = tcp.NewServer(GlobalConfig.TCP_BIND, &handle{})
	return cfcServer.ListenAndServe()
}

func UpdatHostAive() {
	go func() {
		for {
			for _, host := range mydb.GetNoMaintainHost() {
				time_diff := time.Now().Sub(host.UpdateAt).Seconds()
				if time_diff > 60 {
					if host.IsAlive() {
						lg.Info("set host(%s %s) status down, time difference:%0.2fs", host.IP, host.Hostname, time_diff)
						mydb.SetHostAlive(host.ID, "0")
					}
				} else {
					if !host.IsAlive() {
						lg.Info("set host(%s %s) status ok, time difference:%0.2fs ", host.IP, host.Hostname, time_diff)
						mydb.SetHostAlive(host.ID, "1")
					}
				}
			}
			time.Sleep(time.Minute * 1)
		}
	}()
}

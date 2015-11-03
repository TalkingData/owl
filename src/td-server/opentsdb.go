package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func WiteToTSDB(channel chan []byte, cfg *Config) {
START:
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.OPENTSDB_ADDR)
	if err != nil {
		panic(fmt.Sprintf("error opentsdb address(%s)", cfg.OPENTSDB_ADDR))
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		slog.Error("connect tsdb server(%s) error(%s)", tcpAddr, err)
		time.Sleep(time.Second * 30)
		goto START
	}
	slog.Info("connected to opentsdb server %s", cfg.OPENTSDB_ADDR)
	conn.SetKeepAlive(true)
	defer conn.Close()
	for {
		select {
		case data, ok := <-channel:
			if ok {
				data = append(data, '\n')
				length, err := conn.Write(data)
				if err != nil {
					slog.Error("write opentsdb error %s", err)
					channel <- bytes.Trim(data, "\n")
					time.Sleep(time.Second * 30)
					goto START
				}
				dlog.Info("write opentsdb %d bytes, data:(%s)", length, string(bytes.Trim(data, "\n")))
			}
		}
	}
}

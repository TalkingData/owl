package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"tcpserver"
)

func StartTCPServe(cfg *Config) *tcpserver.Server {
	addr, err := net.ResolveTCPAddr("tcp", cfg.TCPBIND)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	s := tcpserver.NewServer(listener, &handle{})
	s.SetPacketLimitSize(uint32(cfg.MAX_PACKET_SIZE))
	go s.Start()
	return s
}

func StartHttpServe(cfg *Config) {
	h := http.FileServer(http.Dir("./update"))
	slog.Info("start http file server on ./update")
	err := http.ListenAndServe(cfg.HTTPBIND, h)
	if err != nil {
		fmt.Println("start http file server failed ", err.Error())
		os.Exit(1)
	}
	slog.Info("start http listen %s done", cfg.HTTPBIND)
}

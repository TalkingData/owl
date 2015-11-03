package main

import (
	"net"
	"tcpserver"
)

func NewTCPServe() *tcpserver.Server {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.TCPBIND)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		panic(err)
	}
	return tcpserver.NewServer(listener,  &handle{})
}

package main

import "tcpserver"

type handle struct {
}

//数据包逻辑处理
func (this *handle) HandlerMessage(conn *tcpserver.Conn, data []byte) {
	defer func() {
		recover()
	}()
	switch data[0] {
	case HOSTHB:
		//log.Info("receive host heatbeart packet")
	default:
		log.Critical(string(data[1:]))
	}
}

func (this *handle) Connect(conn *tcpserver.Conn) {
	log.Info("%s connected", conn.RemoteAddr())
}

func (this *handle) Disconnect(conn *tcpserver.Conn) {
	log.Info("%s connected", conn.RemoteAddr())
}

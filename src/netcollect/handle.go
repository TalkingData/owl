package main

import (
	"owl/common/types"

	"github.com/wuyingsong/tcp"
)

// type handle struct {
// }

// func (this *handle) MakeSession(sess *tcp.Session) {
// 	lg.Info("%s new connection ", sess.RemoteAddr())
// }

// func (this *handle) LostSession(sess *tcp.Session) {
// 	lg.Info("%s disconnect ", sess.RemoteAddr())
// }

// //数据包逻辑处理
// func (this *handle) Handle(sess *tcp.Session, data []byte) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			lg.Error("Recovered in Handle", r)
// 		}
// 	}()
// 	mt := types.MessageType(data[0])
// 	lg.Info("received %v %v", types.MessageTypeText[mt], string(data[1:]))
// }

type callback struct {
}

func (cb *callback) OnConnected(conn *tcp.TCPConn) {
	lg.Info("callback:%s connected", conn.GetRemoteAddr().String())
}

//链接断开回调
func (cb *callback) OnDisconnected(conn *tcp.TCPConn) {
	lg.Info("callback:%s disconnect ", conn.GetRemoteAddr().String())
}

//错误回调
func (cb *callback) OnError(err error) {
	lg.Error("callback: %s", err)
}

//消息处理回调
func (cb *callback) OnMessage(conn *tcp.TCPConn, p tcp.Packet) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in OnMessage", r)
		}
	}()
	pkt := p.(*tcp.DefaultPacket)
	lg.Info("recived %v %v", types.MsgTextMap[pkt.Type], string(pkt.Body))
}

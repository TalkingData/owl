package main

import (
	"owl/common/types"

	"github.com/wuyingsong/tcp"
)

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
	switch pkt.Type {
	case types.ALAR_MESS_INSPECTOR_TASKS:
		at := &types.AlarmTasks{}
		if err := at.Decode(pkt.Body); err != nil {
			lg.Error(err.Error())
			return
		}
		if len(at.Tasks) == 0 {
			return
		}
		lg.Info("receive %v %v", types.AlarmMessageTypeText[pkt.Type], string(pkt.Body))
		inspector.taskPool.PutTasks(at.Tasks)
	default:
		lg.Error("unknown option: %v", pkt.Type)
	}
}

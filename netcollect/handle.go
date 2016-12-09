package main

import (
	"owl/common/tcp"
	"owl/common/types"
)

type handle struct {
}

func (this *handle) MakeSession(sess *tcp.Session) {
	lg.Info("%s new connection ", sess.RemoteAddr())
}

func (this *handle) LostSession(sess *tcp.Session) {
	lg.Info("%s disconnect ", sess.RemoteAddr())
}

//数据包逻辑处理
func (this *handle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in Handle", r)
		}
	}()
	mt := types.MessageType(data[0])
	lg.Info("received %v %v", types.MessageTypeText[mt], string(data[1:]))
}

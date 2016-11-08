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

type cfcHandle struct {
	handle
}

//数据包逻辑处理
func (this *cfcHandle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in HandlerMessage", r)
		}
	}()
	mt := types.MessageType(data[0])
	lg.Info("receive %v %v", types.MessageTypeText[mt], string(data[1:]))
	var err error
	switch mt {
	case types.MESS_GET_HOST_PLUGIN_LIST_RESP:
		resp := types.GetPluginResp{}
		if err = resp.Decode(data[1:]); err != nil {
			lg.Error("decode plugin response error %s %s ", err.Error(), string(data[1:]))
			return
		}
		DelNotUsePlugin(resp.Plugins)
		MergePlugin(resp.Plugins)
	default:
	}
}

type repeaterHandle struct {
	handle
}

func (this *repeaterHandle) Handle(sess *tcp.Session, data []byte) {
}

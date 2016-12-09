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
	var err error
	mt := types.MessageType(data[0])
	lg.Info("receive %v %v", types.MessageTypeText[mt], string(data[1:]))
	if proxy.cfc.IsClosed() {
		return
	}
	switch mt {

	case types.MESS_POST_METRIC, types.MESS_POST_HOST_CONFIG, types.MESS_POST_HOST_ALIVE:
		proxy.cfc.Send(data)
	case types.MESS_GET_HOST_PLUGIN_LIST:
		proxy.cfc.Send(data)
		host := types.Host{}
		if err = host.Decode(data[1:]); err != nil {
			lg.Error("decode host error ", err.Error())
			return
		}
		if resp := proxy.PluginList.Get(host.ID); len(resp) > 0 {
			sess.Send(
				types.Pack(
					types.MESS_GET_HOST_PLUGIN_LIST_RESP,
					&types.GetPluginResp{host.ID, resp},
				))
		}
	case types.MESS_GET_HOST_PLUGIN_LIST_RESP:
		resp := types.GetPluginResp{}
		if err = resp.Decode(data[1:]); err != nil {
			lg.Error("decode plugin response error %s", err)
			return
		}
		proxy.PluginList.Set(resp.HostID, resp.Plugins)
	default:
		lg.Error("Unknown Option: %v", mt)
	}
}

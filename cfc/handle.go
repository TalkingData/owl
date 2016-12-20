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
//TODO: 获取所有插件md5、更新插件
func (this *handle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in HandlerMessage", r)
		}
	}()
	var err error
	mt := types.MessageType(data[0])
	lg.Info("receive %v %v", types.MessageTypeText[mt], string(data[1:]))
	switch mt {
	//agent上传基础信息, ip hostname agentVersion
	case types.MESS_POST_HOST_CONFIG:
		host := &types.Host{}
		err = host.Decode(data[1:])
		if err != nil {
			lg.Error(err.Error())
			sess.Close()
			return
		}
		if h := mydb.GetHost(host.ID); h == nil { //主机不存在
			err := mydb.CreateHost(host.ID, host.SN, host.IP, host.Hostname, host.AgentVersion)
			if err != nil {
				lg.Error("create host error: %s", err.Error())
			} else {
				lg.Info("create host:%v", host)
			}
		} else if h.IP != host.IP ||
			h.Hostname != host.Hostname ||
			h.AgentVersion != host.AgentVersion ||
			h.SN != host.SN {
			lg.Info("update host: %v->%v", h, host)
			mydb.UpdateHost(host)
		}
		//同步metric
	case types.MESS_POST_METRIC:
		cfg := types.MetricConfig{}
		err = cfg.Decode(data[1:])
		if err != nil {
			lg.Error(err.Error())
			sess.Close()
		}
		host := mydb.GetHost(cfg.HostID)
		if host == nil {
			return
		}
		if mydb.MetricIsExists(cfg.HostID, cfg.SeriesData.GetMetric()) {
			return
		}
		if err = mydb.CreateMetric(cfg.HostID, cfg.SeriesData); err != nil {
			lg.Error("create metric error %v", err.Error())
		}
		//获取需要执行的插件列表
	case types.MESS_GET_HOST_PLUGIN_LIST:
		host := &types.Host{}
		err = host.Decode(data[1:])
		if err != nil {
			return
		}
		sess.Send(types.Pack(types.MESS_GET_HOST_PLUGIN_LIST_RESP,
			&types.GetPluginResp{HostID: host.ID, Plugins: mydb.GetPlugins(host.ID)}),
		)
	case types.MESS_POST_HOST_ALIVE:
		cfg := &types.Host{}
		err = cfg.Decode(data[1:])
		if err != nil {
			lg.Error(err.Error())
			sess.Close()
		}
		host := mydb.GetHost(cfg.ID)
		if host == nil {
			return
		}
		mydb.UpdateHost(host)
	default:
		lg.Error("Unknown Option: %v", mt)
	}
}

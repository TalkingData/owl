package main

import (
	"os"
	"owl/common/types"
	"path/filepath"

	"github.com/wuyingsong/tcp"
)

type handle struct {
}

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

func (cb *callback) OnMessage(conn *tcp.TCPConn, p tcp.Packet) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in OnMessage", r)
		}
	}()
	pkt := p.(*tcp.DefaultPacket)
	switch pkt.Type {
	case types.MsgCFCSendPluginsList:
		resp := types.GetPluginResp{}
		if err := resp.Decode(pkt.Body); err != nil {
			lg.Error("decode plugin response error %s", err)
			return
		}
		lg.Debug("recive message, type:%s, body:%s", types.MsgTextMap[pkt.Type], string(pkt.Body))
		removeNoUsePlugin(resp.Plugins)
		mergePlugin(resp.Plugins)
	case types.MsgCFCSendReconnect:
		conn.Close()
	case types.MsgCFCSendPlugin:
		sp := types.SyncPluginResponse{}
		if err := sp.Decode(pkt.Body); err != nil {
			lg.Error("%s", err)
			return
		}
		lg.Debug("recive message, %s %s", types.MsgTextMap[pkt.Type], sp.Path)
		filename := filepath.Join(GlobalConfig.PluginDir, sp.Path)
	retry:
		fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			if os.IsNotExist(err) {
				dir := filepath.Dir(filename)
				lg.Warn("create plugin failed, dir(%s) is not exists, create", dir)
				if err = os.MkdirAll(dir, 0755); err != nil {
					lg.Warn("mkdir %s failed, error:%s", dir, err.Error())
					return
				}
				goto retry
			}
			lg.Error("%s", err)
			return
		}
		defer fd.Close()
		writeLen, err := fd.Write(sp.Body)
		if err != nil {
			lg.Error("create plugin error(%s), %s", err, sp.Path)
			return
		}
		lg.Debug("create plugin(%s) successfully, write %d bytes.", sp.Path, writeLen)
	default:
		lg.Error("unsupport packet type %v", pkt.Type)
		conn.Close()
	}

}

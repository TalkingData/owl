package main

import (
	"io/ioutil"
	"os"
	"owl/common/types"
	"path/filepath"

	"github.com/wuyingsong/tcp"
	"github.com/wuyingsong/utils"
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
	case types.MsgAgentRequestSyncPlugins:
		spr := types.SyncPluginRequest{}
		if err := spr.Decode(pkt.Body); err != nil {
			lg.Error("decode SyncPluginRequest error", err)
			return
		}
		lg.Debug("%s, %s", types.MsgTextMap[pkt.Type], spr)
		pth := filepath.Join(GlobalConfig.PluginDir, spr.Path)
		md5String, err := utils.GetFileMD5(pth)
		if os.IsExist(err) || md5String != spr.Checksum {
			proxy.sendSyncPluginRequest(spr.HostID, spr.Plugin)
			return
		}
		fd, err := os.Open(pth)
		if err != nil {
			lg.Error("%s", err)
			return
		}
		sp := types.SyncPluginResponse{
			Path: spr.Path,
		}
		fileBytes, err := ioutil.ReadAll(fd)
		if err != nil {
			lg.Error("%s", err)
			return
		}
		sp.Body = fileBytes
		conn.AsyncWritePacket(tcp.NewDefaultPacket(
			types.MsgCFCSendPlugin,
			sp.Encode(),
		))
	case types.MsgAgentRegister:
		host := types.Host{}
		if err := host.Decode(pkt.Body); err != nil {
			lg.Error("decode host error(%s)", err)
			return
		}
		lg.Debug("%s, %v", types.MsgTextMap[pkt.Type], host)
		proxy.lock.Lock()
		if addr, ok := proxy.connMap[host.ID]; !ok || addr != conn.GetRemoteAddr().String() {
			lg.Debug("set conn map  %s -> %s", host.ID, conn.GetRemoteAddr().String())
			proxy.connMap[host.ID] = conn.GetRemoteAddr().String()
		}
		proxy.lock.Unlock()
		proxy.cfc.AsyncWritePacket(pkt)

	case types.MsgAgentSendHeartbeat, types.MsgAgentSendMetricInfo, types.MsgAgentGetPluginsList:
		lg.Debug("%s, %s", types.MsgTextMap[pkt.Type], string(pkt.Body))
		proxy.cfc.AsyncWritePacket(pkt)
	case types.MsgCFCSendPluginsList:
		res := types.GetPluginResp{}
		if err := res.Decode(pkt.Body); err != nil {
			lg.Error("decode plugin response error(%s)", err)
			return
		}
		lg.Debug("%s, %s", types.MsgTextMap[pkt.Type], string(pkt.Body))
		proxy.lock.RLock()
		if addr, ok := proxy.connMap[res.HostID]; ok {
			if c := proxy.srv.GetTCPConn(addr); c != nil {
				c.AsyncWritePacket(pkt)
			}
		}
		proxy.lock.RUnlock()
	case types.MsgCFCSendPlugin:
		res := types.SyncPluginResponse{}
		if err := res.Decode(pkt.Body); err != nil {
			lg.Error("decode sync plugin response error(%s)", err)
			return
		}
		lg.Debug("%s, %s", types.MsgTextMap[pkt.Type], res.Path)
		proxy.lock.RLock()
		if addr, ok := proxy.connMap[res.HostID]; ok {
			if c := proxy.srv.GetTCPConn(addr); c != nil {
				lg.Info("sender plugin to client")
				c.AsyncWritePacket(pkt)
			} else {
				lg.Error("get tcp conn faild", res.HostID)
			}
		} else {
			lg.Warn("get host id failed", res.HostID)
		}
		proxy.lock.RUnlock()
		filename := filepath.Join(GlobalConfig.PluginDir, res.Path)
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
		writeLen, err := fd.Write(res.Body)
		if err != nil {
			lg.Error("create plugin error(%s), %s", err, res.Path)
			return
		}
		lg.Debug("create plugin(%s) successfully, write %d bytes.", res.Path, writeLen)
		fd.Close()
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"owl/common/types"
	"path/filepath"
	"strings"

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
			lg.Error("Recovered in OnMessage %v", r)
		}
	}()
	pkt := p.(*tcp.DefaultPacket)
	lg.Debug("%s, %s", types.MsgTextMap[pkt.Type], string(pkt.Body))
	cb.dispatch(conn, pkt)
}

func (cb *callback) dispatch(conn *tcp.TCPConn, pkt *tcp.DefaultPacket) {
	switch pkt.Type {
	// 客户端注册
	case types.MsgAgentRegister:
		host := &types.Host{}
		if err := host.Decode(pkt.Body); err != nil {
			lg.Error("decode host error", err)
			conn.Close()
			return
		}
		if err := registerAgent(host); err != nil {
			lg.Error("register agent failed, error:%s", err)
			return
		}
		lg.Info("register host %v", host)

	// 客户端获取需要执行的插件列表
	case types.MsgAgentGetPluginsList:
		var (
			err  error
			host types.Host
		)
		// 反序列化
		if err = host.Decode(pkt.Body); err != nil {
			lg.Error("decode host error %s", err)
			conn.Close()
			return
		}
		pluginsResp, err := agentGetPluginList(host)
		if err != nil {
			lg.Error("get host plugin error %s host:%v", err, host)
			return
		}
		// 发送结果集到 agent
		if err = conn.AsyncWritePacket(
			tcp.NewDefaultPacket(
				types.MsgCFCSendPluginsList,
				pluginsResp.Encode(),
			),
		); err != nil {
			lg.Error("send plugin list to agent failed, error:%s", err)
			return
		}
		lg.Info("send plugin list to agent, response:%s", pluginsResp.Encode())

	case types.MsgAgentSendHeartbeat:
		host := &types.Host{}
		if err := host.Decode(pkt.Body); err != nil {
			lg.Error("decode host error %s", err)
			return
		}
		if host.ID == "" {
			lg.Warning("host id is empty %v", host)
			return
		}
		if err := mydb.createOrUpdateHost(host); err != nil {
			lg.Error("update host ")
		}

	case types.MsgAgentRequestSyncPlugins:
		spr := types.SyncPluginRequest{}
		if err := spr.Decode(pkt.Body); err != nil {
			lg.Error("decode SyncPluginRequest error", err)
			return
		}
		resp, err := agentRequestSyncPlugins(spr)
		if err != nil {
			lg.Error("get plugin failed, request:%s, error:%s", spr.Encode(), err)
			return
		}
		if err = conn.AsyncWritePacket(tcp.NewDefaultPacket(
			types.MsgCFCSendPlugin,
			resp.Encode(),
		)); err != nil {
			lg.Error("send sync plugin response to agent failed, request:%s, error:%s",
				spr.HostID, spr.Encode(), err)
			return
		}
		lg.Info("send sync plugin response to agent, request:%s, response:%s",
			spr.Encode(), resp.Encode())

		// 客户端上传metric
	case types.MsgAgentSendMetricInfo:
		metricConfig := types.MetricConfig{}
		// 反序列化
		if err := metricConfig.Decode(pkt.Body); err != nil {
			lg.Error("decode metricConfig error", err)
			return
		}
		agentSendMetricInfo(metricConfig)
	default:
		lg.Warn("%v no callback", types.MsgTextMap[pkt.Type])
	}
}

func getHostIDByHostname(hostname string) string {
	host, err := mydb.getHostByHostname(hostname)
	if err != nil {
		lg.Error("get host id by hostname failed, hostname:%q, error:%s", hostname, err)
	}
	if host == nil {
		return ""
	}
	return host.ID
}

//每次启动时发送
func registerAgent(host *types.Host) error {
	if err := mydb.createOrUpdateHost(host); err != nil {
		return err
	}
	//添加到产品线
	if host.GetMetadata("product") == "" {
		return nil
	}
	productName := strings.TrimSpace(host.GetMetadata("product"))
	var (
		product *types.Product
		err     error
	)

	//产品线不存在并且开启了自动创建
	if product, err = mydb.getProductByName(productName); err != nil {
		return fmt.Errorf("get product %s failed, %s", productName, err)
	}
	// 产品线不存在
	if product == nil {
		// 已开启自动创建
		if GlobalConfig.AutoCreateProduct {
			//创建产品线
			product, err = mydb.createProduct(productName)
			if err != nil {
				return fmt.Errorf("auto create product failed, message:%s", err)
			}
			lg.Info("create product %d:%s", product.ID, product.Name)
		} else {
			lg.Info("auto_create_product is disabled, product %s not created ", productName)
			return nil
		}
	}

	//添加主机到产品线
	if err = mydb.addHostToProduct(host.ID, product.ID); err != nil {
		return fmt.Errorf("add host to product failed, host:%s, product:%s", host.ID, product.Name)
	}
	lg.Info("add host %s:%s to product %s ", host.Hostname, host.IP, product.Name)

	//自动分配主机组
	groupName := host.GetMetadata("group")
	if groupName == "" {
		return nil
	}
	group, err := mydb.getProductHostGroup(product.ID, groupName)
	if err != nil {
		return fmt.Errorf("get product %s host group %s failed, error:%s", productName, groupName, err)
	}
	if group == nil {
		group, err = mydb.createProductHostGroup(product.ID, groupName)
		if err != nil {
			return fmt.Errorf("create host group %s in product %s failed, error:%s", groupName, productName, err)
		}

	}

	if err = mydb.addHost2Group(group.ID, host.ID); err != nil {
		return fmt.Errorf("add host %v to group %s failed, error:%s", host, groupName, err)
	}
	lg.Info("add host %v to product %s host group %s", host, productName, groupName)
	return nil
}

//当插件不存在或者 md5 变化时请求
func agentRequestSyncPlugins(req types.SyncPluginRequest) (*types.SyncPluginResponse, error) {
	pth := filepath.Join(GlobalConfig.PluginDir, req.Path)
	md5String, err := utils.GetFileMD5(pth)
	if err != nil {
		return nil, err
	}
	if md5String != req.Checksum {
		return nil, err
	}
	fd, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	sp := &types.SyncPluginResponse{
		HostID: req.HostID,
		Path:   req.Path,
	}
	fileBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	sp.Body = fileBytes
	return sp, nil
}

//客户端定时请求
func agentGetPluginList(host types.Host) (*types.GetPluginResp, error) {
	// 获取 plugin
	plugins, err := mydb.getHostPlugins(host.ID)
	if err != nil {
		return nil, err
	}
	return &types.GetPluginResp{
		HostID:  host.ID,
		Plugins: plugins,
	}, nil
}

//客户端定时同步
func agentSendMetricInfo(info types.MetricConfig) {
	if info.HostID == "" {
		hostname := info.SeriesData.Tags["host"]
		info.HostID = getHostIDByHostname(hostname)
	}
	info.SeriesData.RemoveTag("host")
	info.SeriesData.RemoveTag("uuid")
	seriesData := info.SeriesData
	seriesData.RemoveTag("host")
	seriesData.RemoveTag("uuid")
	//创建 metric
	if err := mydb.createOrUpdateMetric(
		info.HostID,
		seriesData,
	); err != nil {
		lg.Error("create or update metric error %s metric:%v", err, info)
		return
	}
	lg.Debug("create metric %s", info.Encode())
}

package main

import (
	"fmt"
	"owl/client/builtin"
	"owl/common/types"
	"reflect"
	"time"

	"github.com/wuyingsong/tcp"
)

const (
	//GetHostPluginListInterval 定义agent获取插件间隔, 单位秒
	GetHostPluginListInterval = 5 * 60

	// RunBuiltinMetricCycle 定义了内置指标采集时间间隔,单位秒
	RunBuiltinMetricCycle = 1 * 60

	// HeartbeatInterval 定义了agent同步心跳的时间间隔,单位秒
	HeartbeatInterval = 1 * 60
)

type Agent struct {
	*tcp.AsyncTCPServer
	hostcfg    *types.Host
	cfc        *tcp.TCPConn
	repeater   *tcp.TCPConn
	SendChan   chan types.TimeSeriesData
	tsdHistory map[string]float64
}

var (
	agent Agent
	// AgentVersion 定义owl agent版本号
	AgentVersion = "5.0.0"
)

func InitAgent() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MaxPacketSize))
	s := tcp.NewAsyncTCPServer(GlobalConfig.TCPBind, &callback{}, protocol)
	agent = Agent{
		s,
		&types.Host{},
		&tcp.TCPConn{},
		&tcp.TCPConn{},
		make(chan types.TimeSeriesData, GlobalConfig.BufferSize),
		make(map[string]float64),
	}
	return agent.ListenAndServe()
}

// 连接cfc
func (agent *Agent) dialCFC() error {
	if !agent.cfc.IsClosed() {
		return fmt.Errorf("cfc is already connected")
	}
	conn, err := agent.Connect(GlobalConfig.CfcAddr, nil, nil)
	if err != nil {
		return err
	}
	agent.cfc = conn
	return nil
}

// 连接repeater
func (agent *Agent) dialRepeater() error {
	if !agent.repeater.IsClosed() {
		return fmt.Errorf("repeater is already connected")
	}
	conn, err := agent.Connect(GlobalConfig.RepeaterAddr, nil, nil)
	if err != nil {
		return err
	}
	agent.repeater = conn
	return nil

}

// 守护cfc和repeater连接，失败重连
func (agent *Agent) watchConnLoop() {
	for {
		if agent.cfc.IsClosed() {
			lg.Error("cfc reconnect %v", agent.dialCFC())
		}
		if agent.repeater.IsClosed() {
			lg.Error("repeater reconnect %v", agent.dialRepeater())
		}
		time.Sleep(time.Second * 5)
	}
}

func (agent *Agent) watchHostConfig() {
	for {
		if agent.hostcfg == nil {
			agent.hostcfg = newHostConfig()
			agent.register()
		} else {
			hostcfg := newHostConfig()
			if !reflect.DeepEqual(agent.hostcfg, hostcfg) {
				agent.hostcfg = hostcfg
				agent.register()
			}
		}
		time.Sleep(time.Second * 30)
	}
}

// 发送插件同步请求
func (agent *Agent) sendSyncPluginRequest(p types.Plugin) error {
	lg.Debug("send sync plugin request, %s", p.Path)
	spr := types.SyncPluginRequest{
		agent.hostcfg.ID,
		p,
	}
	return agent.cfc.AsyncWritePacket(
		tcp.NewDefaultPacket(types.MsgAgentRequestSyncPlugins, spr.Encode()),
	)
}

// 发送主机配置信息
func (agent *Agent) register() error {
	lg.Debug("send host config %+v", agent.hostcfg)
	return agent.cfc.AsyncWritePacket(
		tcp.NewDefaultPacket(
			types.MsgAgentRegister,
			agent.hostcfg.Encode(),
		),
	)
}

//发送心跳
func (agent *Agent) sendHeartbeat2CFC() error {
	pkt := tcp.NewDefaultPacket(
		types.MsgAgentSendHeartbeat,
		agent.hostcfg.Encode(),
	)
	lg.Debug("send heartbeat to cfc %v", agent.hostcfg)
	return agent.cfc.AsyncWritePacket(
		pkt,
	)
}

func (agent *Agent) sendGetPluginList() error {
	lg.Debug("send get plugin list request")
	return agent.cfc.AsyncWritePacket(
		tcp.NewDefaultPacket(
			types.MsgAgentGetPluginsList,
			newHostConfig().Encode(),
		))
}

func (agent *Agent) sendMetric2CFC(tsd types.TimeSeriesData) {
	cfg := types.MetricConfig{
		agent.hostcfg.ID,
		tsd,
	}
	if err := agent.cfc.AsyncWritePacket(
		tcp.NewDefaultPacket(
			types.MsgAgentSendMetricInfo,
			cfg.Encode(),
		),
	); err == nil {
		lg.Info("send to cfc %s", tsd)
	} else {
		lg.Error("send time series data to cfc error %s", err)
	}
}

func (agent *Agent) StartTimer() {
	t1 := time.NewTicker(time.Second * GetHostPluginListInterval)
	t3 := time.NewTicker(time.Second * HeartbeatInterval)
	go agent.sendGetPluginList()
	go agent.runBuiltinMetricCollect()
	go agent.sendHeartbeat2CFC()
	for {
		select {
		case <-t1.C:
			agent.sendGetPluginList()
		case <-t3.C:
			agent.sendHeartbeat2CFC()
		}
	}
}

func (agent *Agent) getLocalIPAddress() string {
	if !agent.cfc.IsClosed() {
		return agent.cfc.GetLocalIPAddress()
	}
	if !agent.repeater.IsClosed() {
		return agent.repeater.GetLocalIPAddress()
	}
	return ""
}

func (agent *Agent) SendTSD2Repeater() {
	var err error
	for {
		if len(agent.hostcfg.ID) > 0 && len(agent.hostcfg.Hostname) > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	for tsd := range agent.SendChan {

		tags := map[string]string{"uuid": agent.hostcfg.ID, "host": agent.hostcfg.Hostname}
		pk := tsd.PK()
		history, exist := agent.tsdHistory[pk]
		currValue := tsd.Value
		agent.tsdHistory[pk] = currValue
		if !exist {
			agent.sendMetric2CFC(tsd)
			switch tsd.DataType {
			case "counter", "COUNTER", "DERIVE", "derive":
				continue
			}
		}
		switch tsd.DataType {
		// 计算速率
		case "COUNTER", "counter":
			rate := (tsd.Value - history) / float64(tsd.Cycle)
			if rate < 0 {
				continue
			}
			tsd.Value = rate
		// 计算差值
		case "DERIVE", "derive":
			tsd.Value = tsd.Value - history
		}
		tsd.AddTags(tags)
		for {
			if agent.repeater.IsClosed() {
				err = tcp.ErrConnClosing
				goto retry
			}
			err = agent.repeater.AsyncWritePacket(
				tcp.NewDefaultPacket(
					types.MsgAgentSendTimeSeriesData,
					tsd.Encode(),
				),
			)
			if err == nil {
				break
			}
		retry:
			lg.Warn("send to rpeater error(%s), retry after 5 seconds", err.Error())
			time.Sleep(time.Second * 5)
		}
		lg.Info("send to repeater %s", tsd)
	}
}

func (agent *Agent) sendHeartbeat2Repeater() {
	tsd := types.TimeSeriesData{
		Metric:   "agent.alive",
		DataType: "GAUGE",
		Value:    1,
		Cycle:    HeartbeatInterval,
	}
	for {
		tsd.Timestamp = time.Now().Unix()
		agent.SendChan <- tsd
		time.Sleep(time.Second * HeartbeatInterval)
	}
}

func (agent *Agent) runBuiltinMetricCollect() {
	lg.Debug("run built-in metric collect")
	now := time.Now().Unix()
	diff := 60 - (now % 60)
	time.Sleep(time.Second * time.Duration(diff))
	for {
		go builtin.MemoryMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.SwapMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.LoadMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.NetMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.DiskMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.FdMetrics(RunBuiltinMetricCycle, agent.SendChan)
		go builtin.CpuMetrics(RunBuiltinMetricCycle, agent.SendChan)
		time.Sleep(time.Second * RunBuiltinMetricCycle)
	}
}

func newHostConfig() *types.Host {
	host := &types.Host{}
	host.ID = getHostID()
	host.Uptime, host.IdlePct = getHostUptime()
	host.Hostname = getHostname()
	host.AgentVersion = AgentVersion
	host.IP = agent.getLocalIPAddress()
	return host
}

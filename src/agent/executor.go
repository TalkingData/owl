package main

import (
	"github.com/shopspring/decimal"
	"owl/common/global"
)

func (a *agent) refreshAgentInfo() {
	a.agentInfo.HostId = a.executor.GetHostID()
	a.agentInfo.Hostname = a.executor.GetHostname()
	a.agentInfo.AgentVersion = global.Version
	a.agentInfo.Uptime, a.agentInfo.IdlePct = a.executor.GetHostUptimeAndIdle()

	// 使IdlePct只保留两位小数
	a.agentInfo.IdlePct, _ = decimal.NewFromFloat(a.agentInfo.IdlePct).Round(2).Float64()

	// Get local ip with proxy
	a.agentInfo.Ip = a.executor.GetLocalIp(a.conf.ProxyAddress)
}

func (a *agent) execBuiltinMetrics(ts int64) {
	cycle := int32(a.conf.ExecBuiltinMetricsIntervalSecs)

	a.preprocessTsData(a.executor.ExecCollectAlive(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectCpu(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectDisk(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectFd(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectLoad(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectMem(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectNet(ts, cycle), true)
	a.preprocessTsData(a.executor.ExecCollectSwap(ts, cycle), true)
}

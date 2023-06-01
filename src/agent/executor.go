package main

import (
	"github.com/shopspring/decimal"
	"owl/common/global"
)

func (agent *agent) refreshAgentInfo() {
	agent.agentInfo.HostId = agent.executor.GetHostID()
	agent.agentInfo.Hostname = agent.executor.GetHostname()
	agent.agentInfo.AgentVersion = global.Version
	agent.agentInfo.Uptime, agent.agentInfo.IdlePct = agent.executor.GetHostUptimeAndIdle()

	// 使IdlePct只保留两位小数
	agent.agentInfo.IdlePct, _ = decimal.NewFromFloat(agent.agentInfo.IdlePct).Round(2).Float64()

	// Get local ip with proxy
	agent.agentInfo.Ip = agent.executor.GetLocalIp(agent.conf.ProxyAddress)
}

func (agent *agent) execBuiltinMetrics() {
	cycle := int32(agent.conf.ExecBuiltinMetricsIntervalSecs)

	agent.preprocessTsData(agent.executor.ExecCollectAlive(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectCpu(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectDisk(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectFd(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectLoad(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectMem(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectNet(cycle), true)
	agent.preprocessTsData(agent.executor.ExecCollectSwap(cycle), true)
}

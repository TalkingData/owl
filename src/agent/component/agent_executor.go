package component

func (agent *agent) execBuiltinMetrics() {
	cycle := int32(agent.conf.ExecBuiltinMetricsIntervalSecs)

	agent.sendTsDataArray(agent.executor.ExecCollectAlive(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectCpu(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectDisk(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectFd(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectLoad(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectMem(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectNet(cycle), true)
	agent.sendTsDataArray(agent.executor.ExecCollectSwap(cycle), true)
}

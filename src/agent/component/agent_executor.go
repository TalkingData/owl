package component

func (agent *agent) execBuiltinMetrics(cycle int32) {
	agent.sendManyTsData(agent.executor.ExecCollectAlive(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectCpu(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectDisk(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectFd(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectLoad(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectMem(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectNet(cycle))
	agent.sendManyTsData(agent.executor.ExecCollectSwap(cycle))
}

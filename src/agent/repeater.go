package main

import (
	"context"
	"owl/common/logger"
	proxyProto "owl/proxy/proto"
)

func (agent *agent) sendTimeSeriesData(tsData *proxyProto.TsData) {
	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	agent.logger.DebugWithFields(logger.Fields{
		"metric":    tsData.Metric,
		"data_type": tsData.DataType,
		"value":     tsData.Value,
		"timestamp": tsData.Timestamp,
		"cycle":     tsData.Cycle,
		"tags":      tsData.Tags,
	}, "agent.sendTimeSeriesData called.")

	if _, err := agent.proxyCli.ReceiveTimeSeriesData(ctx, tsData); err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"host_id": agent.agentInfo.HostId,
			"metrics": tsData.Metric,
			"error":   err,
		}, "An error occurred while proxyCli.ReceiveTimeSeriesData in agent.sendTimeSeriesData.")
	}
}

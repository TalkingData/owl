package main

import (
	"context"
	"owl/common/logger"
	commonpb "owl/common/proto"
)

func (agent *agent) sendTsData(force bool) {
	buffLen := agent.tsDataBuff.Len()
	if buffLen < 1 || (!force && buffLen < agent.conf.SendTsDataBatchSize) {
		return
	}

	tsDataArr := agent.tsDataBuff.Get(agent.conf.SendTsDataBatchSize)
	agent.sendTimeSeriesDataArray(&commonpb.TsDataArray{Data: tsDataArr})

	agent.sendTsData(force)
}

func (agent *agent) sendTimeSeriesDataArray(tsDataArr *commonpb.TsDataArray) {
	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	if _, err := agent.proxyCli.ReceiveTimeSeriesDataArray(ctx, tsDataArr); err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"host_id":        agent.agentInfo.HostId,
			"ts_data_length": len(tsDataArr.Data),
			"error":          err,
		}, "An error occurred while proxyCli.ReceiveTimeSeriesDataArray in agent.sendTimeSeriesDataArray.")
	}
}

func (agent *agent) sendTimeSeriesData(tsData *commonpb.TsData) {
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

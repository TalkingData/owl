package component

import (
	"context"
	"owl/common/logger"
	repProto "owl/repeater/proto"
)

func (agent *agent) sendTimeSeriesData(tsData *repProto.TsData) {
	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallRepeaterTimeoutSecs)
	defer cancel()

	agent.logger.DebugWithFields(logger.Fields{
		"metric":    tsData.Metric,
		"data_type": tsData.DataType,
		"value":     tsData.Value,
		"timestamp": tsData.Timestamp,
		"cycle":     tsData.Cycle,
		"tags":      tsData.Tags,
	}, "Finished repCli.ReceiveTimeSeriesData.")

	if _, err := agent.repCli.ReceiveTimeSeriesData(ctx, tsData); err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"host_id": agent.agentInfo.HostId,
			"metrics": tsData.Metric,
			"error":   err,
		}, "An error occurred while repCli.ReceiveTimeSeriesData in agent.sendTimeSeriesData.")
	}
}

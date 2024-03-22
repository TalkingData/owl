package main

import (
	"context"
	"owl/common/logger"
	commonpb "owl/common/proto"
)

func (a *agent) sendTsData(force bool) {
	a.logger.DebugWithFields(logger.Fields{
		"force":                 force,
		"ts_data_buffer_length": a.tsDataBuffer.Len(),
	}, "agent.sendTsData called, will send time series data form agent.tsDataBuffer.")
	defer a.logger.Debug("agent.sendTsData end.")

	buffLen := a.tsDataBuffer.Len()
	if buffLen < 1 || (!force && buffLen < a.conf.SendTsDataBatchSize) {
		return
	}

	tsDataArr := a.tsDataBuffer.Get(a.conf.SendTsDataBatchSize)
	a.sendTimeSeriesDataArray(&commonpb.TsDataArray{Data: tsDataArr})

	a.sendTsData(force)
}

func (a *agent) sendTimeSeriesDataArray(tsDataArr *commonpb.TsDataArray) {
	ctx, cancel := context.WithTimeout(a.ctx, a.conf.SendTsDataArrayTimeoutSecs)
	defer cancel()

	if _, err := a.proxyCli.ReceiveTimeSeriesDataArray(ctx, tsDataArr); err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"host_id":            a.agentInfo.HostId,
			"ts_data_arr_length": len(tsDataArr.Data),
			"error":              err,
		}, "An error occurred while calling proxyCli.ReceiveTimeSeriesDataArray.")
	}
}

func (a *agent) sendTimeSeriesData(tsData *commonpb.TsData) {
	ctx, cancel := context.WithTimeout(a.ctx, a.conf.SendTsDataTimeoutSecs)
	defer cancel()

	a.logger.DebugWithFields(logger.Fields{
		"metric":    tsData.Metric,
		"data_type": tsData.DataType,
		"value":     tsData.Value,
		"timestamp": tsData.Timestamp,
		"cycle":     tsData.Cycle,
		"tags":      tsData.Tags,
	}, "agent.sendTimeSeriesData called.")

	if _, err := a.proxyCli.ReceiveTimeSeriesData(ctx, tsData); err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"host_id": a.agentInfo.HostId,
			"metric":  tsData.Metric,
			"error":   err,
		}, "An error occurred while calling proxyCli.ReceiveTimeSeriesData.")
	}
}

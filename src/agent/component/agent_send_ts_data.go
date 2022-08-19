package component

import (
	"owl/common/logger"
	"owl/dto"
	metricList "owl/dto/metric_list"
)

// sendTsDataArray 发送数组形式的Ts数据，若fillAgentInfo为true将向tags中覆盖填充当前agentInfo的host和uuid
func (agent *agent) sendTsDataArray(in dto.TsDataArray, fillAgentInfo bool) {
	agent.logger.InfoWithFields(logger.Fields{
		"length": len(in),
	}, "agent.sendTsDataArray called.")
	defer agent.logger.Info("agent.sendTsDataArray end.")

	go func() {
		for _, currTsData := range in {
			if err := currTsData.Validate(); err != nil {
				agent.logger.WarnWithFields(logger.Fields{
					"current_ts_data_metric": currTsData.Metric,
					"current_ts_data_type":   currTsData.DataType,
					"current_ts_data_value":  currTsData.Value,
					"current_ts_data_cycle":  currTsData.Cycle,
					"current_ts_data_tags":   currTsData.Tags,
					"error":                  err,
				}, "Ts data validate failed in agent.sendTsDataArray, skipped this.")
			}

			// 取出旧数据
			previousMetric, exist := agent.metricList.Get(currTsData.GetPk())
			newMetric := metricList.NewMetric(
				agent.agentInfo.HostId,
				currTsData.Metric,
				currTsData.DataType,
				currTsData.Value,
				currTsData.Timestamp,
				currTsData.Cycle,
				currTsData.Tags,
			)

			// 对于不存在的Metric，立即上报
			if !exist {
				agent.reportAgentMetric(newMetric.ToCfcMetric())
				// 更新metric list
				agent.metricList.Put(newMetric.GetPk(), newMetric)
			} else if previousMetric.Timestamp+int64(currTsData.Cycle) <= currTsData.Timestamp {
				// 当有旧数据且其 ts+采集间隔 的值，仍然<=最新数据的ts值，则更新数据
				// 此判断是防止COUNTER、DERIVE类型数据被超频率采集后导致Value接近于0的保护处理
				// 更新metric list
				agent.metricList.Put(newMetric.GetPk(), newMetric)
			}

			switch currTsData.DataType {
			// 对于DataType==COUNTER的数据，需要与上一次数据处理后才可发送
			case dto.TsDataTypeCounter:
				if !exist {
					agent.logger.WarnWithFields(logger.Fields{
						"current_ts_data_metric": currTsData.Metric,
						"current_ts_data_type":   currTsData.DataType,
						"current_ts_data_value":  currTsData.Value,
						"current_ts_data_cycle":  currTsData.Cycle,
						"current_ts_data_tags":   currTsData.Tags,
					}, "Ts data type is 'COUNTER', but previous ts data not found, skipped this.")
					continue
				}
				if currTsData.Cycle == 0 {
					agent.logger.ErrorWithFields(logger.Fields{
						"current_ts_data_metric": currTsData.Metric,
						"current_ts_data_type":   currTsData.DataType,
						"current_ts_data_value":  currTsData.Value,
						"current_ts_data_cycle":  currTsData.Cycle,
						"current_ts_data_tags":   currTsData.Tags,
						"previous_ts_data_value": previousMetric.Value,
					}, "An error occurred while agent.sendTsDataArray, the cycle value of ts data can not be 'ZERO'.")
					continue
				}

				rate := (currTsData.Value - previousMetric.Value) / float64(currTsData.Cycle)
				if rate < 0 {
					agent.logger.WarnWithFields(logger.Fields{
						"current_ts_data_metric": currTsData.Metric,
						"current_ts_data_type":   currTsData.DataType,
						"current_ts_data_value":  currTsData.Value,
						"current_ts_data_cycle":  currTsData.Cycle,
						"current_ts_data_tags":   currTsData.Tags,
						"previous_ts_data_value": previousMetric.Value,
						"rate":                   rate,
					}, "Ts data type is 'COUNTER', but the rate<=0, skipped this.")
					continue
				}
				currTsData.Value = rate
			// 对于DataType==DERIVE的数据，需要与上一次数据处理后才可发送
			case dto.TsDataTypeDerive:
				if !exist {
					agent.logger.WarnWithFields(logger.Fields{
						"current_ts_data_metric": currTsData.Metric,
						"current_ts_data_type":   currTsData.DataType,
						"current_ts_data_value":  currTsData.Value,
						"current_ts_data_cycle":  currTsData.Cycle,
						"current_ts_data_tags":   currTsData.Tags,
					}, "Ts data type is 'DERIVE', but the previous ts data not found, skipped this.")
					continue
				}
				currTsData.Value = currTsData.Value - previousMetric.Value
			}

			fillTags := map[string]string{}
			// 向tags中追加uuid和host属性
			if fillAgentInfo {
				fillTags["uuid"] = agent.agentInfo.HostId
				fillTags["host"] = agent.agentInfo.Hostname
			}
			currTsData.MergeTags(fillTags)
			agent.sendTimeSeriesData(currTsData.Trans2RepTsData())
		}
	}()
}

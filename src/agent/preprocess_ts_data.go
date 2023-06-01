package main

import (
	"owl/common/logger"
	"owl/dto"
)

// preprocessTsData 预处理TsData，若fillAgentInfo为true将向tags中覆盖填充当前agentInfo的host和uuid
func (agent *agent) preprocessTsData(in dto.TsDataArray, fillAgentInfo bool) {
	agent.logger.InfoWithFields(logger.Fields{
		"length": len(in),
	}, "agent.preprocessTsData called.")
	defer agent.logger.Info("agent.preprocessTsData end.")

	go func() {
		for _, currData := range in {
			if err := currData.Validate(); err != nil {
				agent.logger.WarnWithFields(logger.Fields{
					"current_data_metric": currData.Metric,
					"current_data_type":   currData.DataType,
					"current_data_value":  currData.Value,
					"current_data_cycle":  currData.Cycle,
					"current_data_tags":   currData.Tags,
					"error":               err,
				}, "Ts data validate failed in agent.preprocessTsData, skipped this.")
			}

			// 取出旧数据
			prevData, exist := agent.tsDataMap.Get(currData.GetPk())
			newTsData := currData.DeepCopyTsData()

			// 对于不存在的Metric，立即上报
			if !exist {
				agent.reportAgentMetric(newTsData.ToCommonMetric(agent.agentInfo.HostId))
				// 更新tsDataMap
				agent.tsDataMap.Put(newTsData.GetPk(), newTsData)
			} else if prevData.Timestamp+int64(currData.Cycle) <= currData.Timestamp {
				// 当有旧数据且其 ts+采集间隔 的值，仍然<=最新数据的ts值，则更新数据
				// 此判断是防止COUNTER、DERIVE类型数据被超频率采集后导致Value接近于0的保护处理
				// 更新tsDataMap
				agent.tsDataMap.Put(newTsData.GetPk(), newTsData)
			}

			if err := agent.dtHandlerMap.Get(currData.DataType)(exist, currData, prevData); err != nil {
				f := logger.Fields{
					"current_data_metric":    currData.Metric,
					"current_data_type":      currData.DataType,
					"current_data_value":     currData.Value,
					"current_data_timestamp": currData.Timestamp,
					"current_data_cycle":     currData.Cycle,
					"current_data_tags":      currData.Tags,
					"previous_exist":         exist,
					"error":                  err,
				}
				if exist {
					f["previous_data_value"] = prevData.Value
					f["previous_data_timestamp"] = prevData.Timestamp
					f["previous_data_cycle"] = prevData.Cycle
				}
				agent.logger.ErrorWithFields(f, "An error occurred while dtHandler in agent.preprocessTsData")
				continue
			}

			fillTags := map[string]string{}
			// 向tags中追加uuid和host属性
			if fillAgentInfo {
				fillTags["uuid"] = agent.agentInfo.HostId
				fillTags["host"] = agent.agentInfo.Hostname
			}
			currData.MergeTags(fillTags)
			agent.tsDataBuff.Put(currData.Trans2CommonTsData())
			agent.sendTsData(false)
		}
	}()
}

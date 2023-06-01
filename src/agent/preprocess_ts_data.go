package main

import (
	"owl/common/logger"
	"owl/dto"
)

// preprocessTsData 预处理TsData，若fillAgentInfo为true将向tags中覆盖填充当前agentInfo的host和uuid
func (a *agent) preprocessTsData(in dto.TsDataArray, fillAgentInfo bool) {
	a.logger.InfoWithFields(logger.Fields{
		"length": len(in),
	}, "agent.preprocessTsData called.")
	defer a.logger.Info("agent.preprocessTsData end.")

	go func() {
		for _, currData := range in {
			if err := currData.Validate(); err != nil {
				a.logger.WarnWithFields(logger.Fields{
					"current_data_metric": currData.Metric,
					"current_data_type":   currData.DataType,
					"current_data_value":  currData.Value,
					"current_data_cycle":  currData.Cycle,
					"current_data_tags":   currData.Tags,
					"error":               err,
				}, "Ts data validate failed in agent.preprocessTsData, skipped this.")
			}

			// 取出旧数据
			currDataPk := currData.GetPk()
			prevData, exist := a.tsDataMap.Get(currDataPk)
			newTsData := currData.DeepCopyTsData()

			// 对于不存在的Metric，立即上报
			if !exist {
				a.reportAgentMetric(newTsData.ToCommonMetric(a.agentInfo.HostId))
				// 更新tsDataMap
				a.tsDataMap.Put(newTsData.GetPk(), newTsData)
			} else if prevData.Timestamp+int64(currData.Cycle) <= currData.Timestamp {
				// 当有旧数据且其 ts+采集间隔 的值，仍然<=最新数据的ts值，则更新数据
				// 此判断是防止COUNTER、DERIVE类型数据被超频率采集后导致Value接近于0的保护处理
				// 更新tsDataMap
				a.tsDataMap.Put(newTsData.GetPk(), newTsData)
			}

			if err := a.dtHandlerMap.Get(currData.DataType)(exist, currData, prevData); err != nil {
				f := logger.Fields{
					"current_data_metric":    currData.Metric,
					"current_data_type":      currData.DataType,
					"current_data_value":     currData.Value,
					"current_data_timestamp": currData.Timestamp,
					"current_data_cycle":     currData.Cycle,
					"current_data_tags":      currData.Tags,
					"current_data_pk":        currDataPk,
					"previous_exist":         exist,
					"error":                  err,
				}
				if exist {
					f["previous_data_value"] = prevData.Value
					f["previous_data_timestamp"] = prevData.Timestamp
					f["previous_data_cycle"] = prevData.Cycle
				}
				a.logger.WarnWithFields(f, "An error occurred while calling dtHandler.")
				continue
			}

			fillTags := map[string]string{}
			// 向tags中追加uuid和host属性
			if fillAgentInfo {
				fillTags["uuid"] = a.agentInfo.HostId
				fillTags["host"] = a.agentInfo.Hostname
			}
			currData.MergeTags(fillTags)
			a.tsDataBuffer.Put(currData.Trans2CommonTsData())
			a.sendTsData(false)
		}
	}()
}

package main

import (
	"context"
	"owl/common/logger"
	commonpb "owl/common/proto"
	"owl/common/utils"
	pluginList "owl/dto/plugin_list"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// reportHeartbeat 发送当前Agent心跳数据
func (a *agent) reportHeartbeat() {
	a.logger.Info("agent.reportHeartbeat called.")
	defer a.logger.Info("agent.reportHeartbeat end.")

	req := &commonpb.AgentInfo{
		HostId:       a.agentInfo.HostId,
		Ip:           a.agentInfo.Ip,
		Hostname:     a.agentInfo.Hostname,
		AgentVersion: a.agentInfo.AgentVersion,
		AgentOs:      runtime.GOOS,
		AgentArch:    runtime.GOARCH,
		Uptime:       a.agentInfo.Uptime,
		IdlePct:      a.agentInfo.IdlePct,
		Metadata:     a.agentInfo.Metadata,
	}

	ctx, cancel := context.WithTimeout(a.ctx, a.conf.ReportHeartbeatTimeoutSecs)
	defer cancel()

	_, err := a.proxyCli.ReceiveAgentHeartbeat(ctx, req)
	if err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling agent.proxyCli.ReceiveAgentHeartbeat.")
	}
}

func (a *agent) reportAgentAllMetrics() {
	a.logger.InfoWithFields(logger.Fields{
		"length": a.tsDataMap.Len(),
	}, "agent.reportAgentAllMetrics called.")
	defer a.logger.Info("agent.reportAgentAllMetrics end.")

	// 将tsDataMap中的数据分成多批，每批agent.conf.ReportMetricBatchSize个，并分别作为参数交给reportAgentMetrics处理
	commonMetricList := make([]*commonpb.Metric, 0, a.conf.ReportMetricBatchSize)
	counter := 0

	currTs := time.Now().Unix()
	for _, tsData := range a.tsDataMap.List() {
		// 如果当前时间戳减去tsData的时间戳大于两倍的采集周期，则删除该tsData
		if int64(tsData.Cycle)*int64(a.conf.CleanExpiredMetricCycleExpiredRatio)+tsData.Timestamp < currTs {
			a.tsDataMap.Remove(tsData.GetPk())
			continue
		}

		commonMetricList = append(commonMetricList, tsData.ToCommonMetric(a.agentInfo.HostId))
		counter++

		if counter >= a.conf.ReportMetricBatchSize {
			// 处理当前批次的逻辑，例如发送到网络或进行其他处理
			a.reportAgentMetrics(&commonpb.Metrics{Metrics: commonMetricList})

			// 重置计数器和批次列表
			counter = 0
			commonMetricList = make([]*commonpb.Metric, 0, a.conf.ReportMetricBatchSize)
		}
	}

	// 处理剩余的不足一批的数据
	if counter > 0 {
		// 处理剩余批次的逻辑
		a.reportAgentMetrics(&commonpb.Metrics{Metrics: commonMetricList})
	}
}

func (a *agent) reportAgentMetrics(in *commonpb.Metrics) {
	ctx, cancel := context.WithTimeout(a.ctx, a.conf.ReportMetricsTimeoutSecs)
	defer cancel()

	if _, err := a.proxyCli.ReceiveAgentMetrics(ctx, in); err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"metrics_length": len(in.Metrics),
			"error":          err,
		}, "An error occurred while calling proxyCli.ReceiveAgentMetrics.")
	}
}

func (a *agent) reportAgentMetric(in *commonpb.Metric) {
	ctx, cancel := context.WithTimeout(a.ctx, a.conf.ReportMetricTimeoutSecs)
	defer cancel()

	if _, err := a.proxyCli.ReceiveAgentMetric(ctx, in); err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"request": in,
			"error":   err,
		}, "An error occurred while calling proxyCli.ReceiveAgentMetric.")
	}
}

// registerAgent 注册当前Agent
func (a *agent) registerAgent() error {
	a.logger.Info("agent.registerAgent called.")
	defer a.logger.Info("agent.registerAgent end.")

	req := &commonpb.AgentInfo{
		HostId:       a.agentInfo.HostId,
		Ip:           a.agentInfo.Ip,
		Hostname:     a.agentInfo.Hostname,
		AgentVersion: a.agentInfo.AgentVersion,
		Uptime:       a.agentInfo.Uptime,
		IdlePct:      a.agentInfo.IdlePct,
		Metadata:     a.agentInfo.Metadata,
	}

	ctx, cancel := context.WithTimeout(a.ctx, a.conf.ReportHeartbeatIntervalSecs)
	defer cancel()

	_, err := a.proxyCli.RegisterAgent(ctx, req)
	if err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling agent.proxyCli.RegisterAgent.")
		return err
	}

	return nil
}

// listPluginsProcess Agent请求Proxy列出自身所需的Plugins，并做相应处理
func (a *agent) listPluginsProcess() {
	a.logger.Info("agent.listPluginsProcess called.")
	defer a.logger.Info("agent.listPluginsProcess end.")

	ctx, cancel := context.WithTimeout(a.ctx, a.conf.ListPluginsTimeoutSecs)
	defer cancel()

	plugins, err := a.proxyCli.ListAgentPlugins(ctx, &commonpb.HostIdReq{HostId: a.agentInfo.HostId})
	if err != nil {
		a.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling agent.proxyCli.ListAgentPlugins.")
		return
	}

	newPluginPkMap := map[string]struct{}{}
	for _, p := range plugins.Plugins {
		a.logger.DebugWithFields(logger.Fields{
			"rsp_plugin_id":       p.Id,
			"rsp_plugin_name":     p.Name,
			"rsp_plugin_path":     p.Path,
			"rsp_plugin_args":     p.Args,
			"rsp_plugin_checksum": p.Checksum,
			"rsp_plugin_interval": p.Interval,
			"rsp_plugin_timeout":  p.Timeout,
		}, "Got agent plugin from proxy.")

		localAbsPath, err := filepath.Abs(path.Join(a.conf.PluginDir, p.Path))
		if err != nil {
			a.logger.ErrorWithFields(logger.Fields{
				"plugin_dir":  a.conf.PluginDir,
				"plugin_path": p.Path,
				"error":       err,
			}, "An error occurred while calling filepath.Abs.")
		}

		newPlugin := pluginList.NewPlugin(
			a.ctx,
			p.Id,
			p.Name, localAbsPath, p.Checksum,
			utils.ParseCommandArgs(p.Args),
			p.Interval, p.Timeout,
			a.conf.ExecuteUntrustedPlugin,
			func(ctx context.Context, ts int64, cycle int32, command string, args ...string) {
				dataArr := a.executor.ExecCollectCmd(ctx, ts, command, args...)
				for _, data := range dataArr {
					data.Cycle = cycle
				}
				a.preprocessTsData(dataArr, true)
			},
		)

		newPluginPk := newPlugin.GetPk()
		newPluginPkMap[newPluginPk] = struct{}{}

		fileChecksum := newPlugin.GetFileChecksum()
		//  插件文件的校验和与本地文件不一致时，需要重新下载插件文件
		if fileChecksum != p.Checksum {
			a.logger.WarnWithFields(logger.Fields{
				"plugin_id":            newPlugin.Id,
				"plugin_name":          newPlugin.Name,
				"response_plugin_path": p.Path,
				"plugin_pathname":      newPlugin.LocalPath,
				"plugin_checksum":      newPlugin.Checksum,
				"plugin_file_checksum": fileChecksum,
			}, "Valid plugin checksum failed, Prepare for download plugin file.")

			// 从proxy下载插件文件，如果失败则跳过
			if err = a.downloadPluginFile(p.Path, newPlugin.LocalPath); err != nil {
				a.logger.ErrorWithFields(logger.Fields{
					"plugin_id":            newPlugin.Id,
					"plugin_name":          newPlugin.Name,
					"response_plugin_path": p.Path,
					"plugin_pathname":      newPlugin.LocalPath,
					"plugin_checksum":      newPlugin.Checksum,
					"error":                err,
				}, "An error occurred while calling agent.proxyCli.downloadPluginFile, Skipped this plugin and task.")
				continue
			}
		}
		if !a.pluginList.Exists(newPluginPk) {
			a.logger.DebugWithFields(logger.Fields{
				"plugin_id":            newPlugin.Id,
				"plugin_name":          newPlugin.Name,
				"response_plugin_path": p.Path,
				"plugin_pathname":      newPlugin.LocalPath,
				"plugin_checksum":      newPlugin.Checksum,
			}, "Put agent plugin to newPluginList.")
			a.pluginList.Put(newPluginPk, newPlugin)
			// 启动新Plugin采集任务
			newPlugin.StartTask()
		}
	}

	// 对于不存在的插件，则移除任务
	for k := range a.pluginList.List() {
		if _, ok := newPluginPkMap[k]; !ok {
			a.pluginList.StopAndRemoveTask(k)
		}
	}
}

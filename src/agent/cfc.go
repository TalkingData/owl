package main

import (
	"context"
	"owl/common/logger"
	pluginList "owl/dto/plugin_list"
	proxyProto "owl/proxy/proto"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// reportHeartbeat 发送当前Agent心跳数据
func (agent *agent) reportHeartbeat() {
	agent.logger.Info("agent.reportHeartbeat called.")
	defer agent.logger.Info("agent.reportHeartbeat end.")

	req := &proxyProto.AgentInfo{
		HostId:       agent.agentInfo.HostId,
		Ip:           agent.agentInfo.Ip,
		Hostname:     agent.agentInfo.Hostname,
		AgentVersion: agent.agentInfo.AgentVersion,
		AgentOs:      runtime.GOOS,
		AgentArch:    runtime.GOARCH,
		Uptime:       agent.agentInfo.Uptime,
		IdlePct:      agent.agentInfo.IdlePct,
		Metadata:     agent.agentInfo.Metadata,
	}

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	_, err := agent.proxyCli.ReceiveAgentHeartbeat(ctx, req)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.proxyCli.ReceiveAgentHeartbeat.")
	}
}

func (agent *agent) reportAgentAllMetrics() {
	agent.logger.InfoWithFields(logger.Fields{
		"length": agent.metricList.Len(),
	}, "agent.reportAgentAllMetrics called.")
	defer agent.logger.Info("agent.reportAgentAllMetrics end.")

	for _, v := range agent.metricList.List() {
		agent.reportAgentMetric(v.ToProxyMetric())
	}
}

func (agent *agent) reportAgentMetric(in *proxyProto.Metric) {
	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	if _, err := agent.proxyCli.ReceiveAgentMetric(ctx, in); err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"request": in,
			"error":   err,
		}, "An error occurred while proxyCli.ReceiveAgentMetric in agent.reportAgentMetric")
	}
}

// registerAgent 注册当前Agent
func (agent *agent) registerAgent() error {
	agent.logger.Info("agent.registerAgent called.")
	defer agent.logger.Info("agent.registerAgent end.")

	req := &proxyProto.AgentInfo{
		HostId:       agent.agentInfo.HostId,
		Ip:           agent.agentInfo.Ip,
		Hostname:     agent.agentInfo.Hostname,
		AgentVersion: agent.agentInfo.AgentVersion,
		Uptime:       agent.agentInfo.Uptime,
		IdlePct:      agent.agentInfo.IdlePct,
		Metadata:     agent.agentInfo.Metadata,
	}

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	_, err := agent.proxyCli.RegisterAgent(ctx, req)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.proxyCli.RegisterAgent.")
		return err
	}

	return nil
}

// listPluginsProcess Agent请求Proxy列出自身所需的Plugins，并做相应处理
func (agent *agent) listPluginsProcess() {
	agent.logger.Info("agent.listPluginsProcess called.")
	defer agent.logger.Info("agent.listPluginsProcess end.")

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallProxyTimeoutSecs)
	defer cancel()

	plugins, err := agent.proxyCli.ListAgentPlugins(ctx, &proxyProto.HostIdReq{HostId: agent.agentInfo.HostId})
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.proxyCli.ListAgentPlugins.")
		return
	}

	newPluginPkMap := map[string]struct{}{}
	for _, p := range plugins.Plugins {
		agent.logger.DebugWithFields(logger.Fields{
			"rsp_plugin_id":       p.Id,
			"rsp_plugin_name":     p.Name,
			"rsp_plugin_path":     p.Path,
			"rsp_plugin_args":     p.Args,
			"rsp_plugin_checksum": p.Checksum,
			"rsp_plugin_interval": p.Interval,
			"rsp_plugin_timeout":  p.Timeout,
		}, "Got agent plugin from proxy.")

		localAbsPath, err := filepath.Abs(path.Join(agent.conf.PluginDir, p.Path))
		if err != nil {
			agent.logger.ErrorWithFields(logger.Fields{
				"plugin_dir":  agent.conf.PluginDir,
				"plugin_path": p.Path,
				"error":       err,
			}, "An error occurred while filepath.Abs in agent.listPluginsProcess.")
		}

		newPlugin := pluginList.NewPlugin(
			agent.ctx,
			p.Id,
			p.Name, localAbsPath, p.Checksum,
			strings.Split(p.Args, " "),
			p.Interval, p.Timeout,
			agent.conf.ExecuteUntrustedPlugin,
			func(ctx context.Context, cycle int32, command string, args ...string) {
				dataArr := agent.executor.ExecCollectCmd(ctx, command, args...)
				for _, data := range dataArr {
					data.Cycle = cycle
				}
				agent.sendTsDataArray(dataArr, true)
			},
		)

		newPluginPk := newPlugin.GetPk()
		newPluginPkMap[newPluginPk] = struct{}{}

		fileChecksum := newPlugin.GetFileChecksum()
		//  插件文件的校验和与本地文件不一致时，需要重新下载插件文件
		if fileChecksum != p.Checksum {
			agent.logger.WarnWithFields(logger.Fields{
				"plugin_id":            newPlugin.Id,
				"plugin_name":          newPlugin.Name,
				"response_plugin_path": p.Path,
				"plugin_pathname":      newPlugin.LocalPath,
				"plugin_checksum":      newPlugin.Checksum,
				"plugin_file_checksum": fileChecksum,
			}, "Valid plugin checksum failed, Prepare for download plugin file.")

			// 从proxy下载插件文件，如果失败则跳过
			if err = agent.downloadPluginFile(p.Path, newPlugin.LocalPath); err != nil {
				agent.logger.ErrorWithFields(logger.Fields{
					"plugin_id":            newPlugin.Id,
					"plugin_name":          newPlugin.Name,
					"response_plugin_path": p.Path,
					"plugin_pathname":      newPlugin.LocalPath,
					"plugin_checksum":      newPlugin.Checksum,
					"error":                err,
				}, "An error occurred while agent.proxyCli.downloadPluginFile, Skipped this plugin and task.")
				continue
			}
		}
		if !agent.pluginList.Exists(newPluginPk) {
			agent.logger.DebugWithFields(logger.Fields{
				"plugin_id":            newPlugin.Id,
				"plugin_name":          newPlugin.Name,
				"response_plugin_path": p.Path,
				"plugin_pathname":      newPlugin.LocalPath,
				"plugin_checksum":      newPlugin.Checksum,
			}, "Put agent plugin to newPluginList.")
			agent.pluginList.Put(newPluginPk, newPlugin)
			// 启动新Plugin采集任务
			newPlugin.StartTask()
		}
	}

	// 对于不存在的插件，则移除任务
	for k := range agent.pluginList.List() {
		if _, ok := newPluginPkMap[k]; !ok {
			agent.pluginList.StopTaskAndRemove(k)
		}
	}
}

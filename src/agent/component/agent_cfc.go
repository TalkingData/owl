package component

import (
	"context"
	"google.golang.org/grpc/status"
	"os"
	cfcProto "owl/cfc/proto"
	"owl/common/logger"
	"owl/common/utils"
	pluginList "owl/dto/plugin_list"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// listPluginsProcess Agent请求CFC列出自身所需的Plugins，并做相应处理
func (agent *agent) listPluginsProcess() {
	agent.logger.Info("agent.listPluginsProcess called.")
	defer agent.logger.Info("agent.listPluginsProcess end.")

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallCfcTimeoutSecs)
	defer cancel()

	plugins, err := agent.cfcCli.ListAgentPlugins(ctx, &cfcProto.HostIdReq{HostId: agent.agentInfo.HostId})
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.cfcCli.ListAgentPlugins.")
		return
	}

	newPluginPkMap := map[string]struct{}{}
	for _, p := range plugins.Plugins {
		agent.logger.DebugWithFields(logger.Fields{
			"response_plugin_id":       p.Id,
			"response_plugin_name":     p.Name,
			"response_plugin_path":     p.Path,
			"response_plugin_args":     p.Args,
			"response_plugin_checksum": p.Checksum,
			"response_plugin_interval": p.Interval,
			"response_plugin_timeout":  p.Timeout,
		}, "Got agent plugin from CFC.")

		pluginPathname, err := filepath.Abs(path.Join(agent.conf.PluginDir, p.Path))
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
			p.Name, pluginPathname, p.Checksum,
			strings.Split(p.Args, " "),
			p.Interval, p.Timeout,
			agent.conf.ExecuteUntrustedPlugin,
			func(ctx context.Context, cycle int32, command string, args ...string) {
				dataArr := agent.executor.ExecCollectCmd(ctx, command, args...)
				for _, data := range dataArr {
					data.Cycle = cycle
				}
				agent.sendManyTsData(dataArr)
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
				"plugin_pathname":      newPlugin.Pathname,
				"plugin_checksum":      newPlugin.Checksum,
				"plugin_file_checksum": fileChecksum,
			}, "Valid plugin checksum failed, Prepare for download plugin file.")

			// 从CFC下载插件文件，如果失败则跳过
			if err = agent.downloadPluginFile(p.Id, newPlugin.Pathname); err != nil {
				agent.logger.ErrorWithFields(logger.Fields{
					"plugin_id":            newPlugin.Id,
					"plugin_name":          newPlugin.Name,
					"response_plugin_path": p.Path,
					"plugin_pathname":      newPlugin.Pathname,
					"plugin_checksum":      newPlugin.Checksum,
					"error":                err,
				}, "An error occurred while agent.cfcCli.downloadPluginFile, Skipped this plugin and task.")
				continue
			}
		}
		if !agent.pluginList.Exists(newPluginPk) {
			agent.logger.DebugWithFields(logger.Fields{
				"plugin_id":            newPlugin.Id,
				"plugin_name":          newPlugin.Name,
				"response_plugin_path": p.Path,
				"plugin_pathname":      newPlugin.Pathname,
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

// reportHeartbeat 发送当前Agent心跳数据
func (agent *agent) reportHeartbeat() {
	agent.logger.Info("agent.reportHeartbeat called.")
	defer agent.logger.Info("agent.reportHeartbeat end.")

	req := &cfcProto.AgentInfo{
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

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallCfcTimeoutSecs)
	defer cancel()

	_, err := agent.cfcCli.ReceiveAgentHeartbeat(ctx, req)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.cfcCli.ReceiveAgentHeartbeat.")
	}
}

// registerAgent 注册当前Agent
func (agent *agent) registerAgent() error {
	agent.logger.Info("agent.registerAgent called.")
	defer agent.logger.Info("agent.registerAgent end.")

	req := &cfcProto.AgentInfo{
		HostId:       agent.agentInfo.HostId,
		Ip:           agent.agentInfo.Ip,
		Hostname:     agent.agentInfo.Hostname,
		AgentVersion: agent.agentInfo.AgentVersion,
		Uptime:       agent.agentInfo.Uptime,
		IdlePct:      agent.agentInfo.IdlePct,
		Metadata:     agent.agentInfo.Metadata,
	}

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallCfcTimeoutSecs)
	defer cancel()

	_, err := agent.cfcCli.RegisterAgent(ctx, req)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.cfcCli.RegisterAgent.")
		return err
	}

	return nil
}

// downloadPluginFile 下载指定插件文件
func (agent *agent) downloadPluginFile(pluginId uint32, pathname string) error {
	agent.logger.Info("agent.downloadPluginFile called.")
	defer agent.logger.Info("agent.downloadPluginFile end.")

	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.DownloadPluginTimeoutSecs)
	defer cancel()

	stream, err := agent.cfcCli.DownloadPluginFile(ctx, &cfcProto.PluginIdReq{PluginId: pluginId})
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while agent.cfcCli.DownloadPluginFile.")
		return err
	}

	_ = os.MkdirAll(path.Dir(pathname), 0755)
	fp, err := os.OpenFile(pathname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"plugin_id": pluginId,
			"pathname":  pathname,
			"error":     err,
		}, "An error occurred while os.OpenFile in agent.downloadPluginFile.")
		return err
	}

	defer func() {
		_ = fp.Close()
	}()

	// 开始分块循环接收文件
	for {
		rsp, err := stream.Recv()
		if err != nil {
			sts := status.Convert(err)
			if sts.Code() == utils.DefaultDownloaderEndOfFileExitCode {
				agent.logger.InfoWithFields(logger.Fields{
					"plugin_id":      pluginId,
					"pathname":       pathname,
					"status_code":    sts.Code(),
					"status_message": sts.Message(),
				}, "agent.downloadPluginFile success by EOF status.Code.")
				return nil
			}
			agent.logger.ErrorWithFields(logger.Fields{
				"plugin_id": pluginId,
				"pathname":  pathname,
				"error":     err,
			}, "An error occurred while stream.Recv in agent.downloadPluginFile.")
			return err
		}

		mob, err := fp.Write(rsp.Buffer)
		if err != nil {
			agent.logger.ErrorWithFields(logger.Fields{
				"plugin_id": pluginId,
				"pathname":  pathname,
				"error":     err,
			}, "An error occurred while fp.Write in agent.downloadPluginFile.")
			return err
		}
		agent.logger.DebugWithFields(logger.Fields{
			"plugin_id":       pluginId,
			"pathname":        pathname,
			"number_of_bytes": mob,
		}, "agent.downloadPluginFile received some data.")
	}
}

func (agent *agent) reportAgentAllMetrics() {
	agent.logger.InfoWithFields(logger.Fields{
		"length": agent.metricList.Len(),
	}, "agent.reportAgentAllMetrics called.")
	defer agent.logger.Info("agent.reportAgentAllMetrics end.")

	for _, v := range agent.metricList.List() {
		agent.reportAgentMetric(v.ToCfcMetric())
	}
}

func (agent *agent) reportAgentMetric(in *cfcProto.Metric) {
	ctx, cancel := context.WithTimeout(agent.ctx, agent.conf.CallCfcTimeoutSecs)
	defer cancel()

	if _, err := agent.cfcCli.ReceiveAgentMetric(ctx, in); err != nil {
		agent.logger.ErrorWithFields(logger.Fields{
			"request": in,
			"error":   err,
		}, "An error occurred while cfcCli.ReceiveAgentMetric in agent.reportAgentMetric")
	}
}

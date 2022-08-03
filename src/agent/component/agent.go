package component

import (
	"context"
	"fmt"
	"os"
	"owl/agent/conf"
	"owl/agent/executor"
	cfcProto "owl/cfc/proto"
	"owl/cli"
	"owl/common/global"
	"owl/common/logger"
	metricList "owl/dto/metric_list"
	pluginList "owl/dto/plugin_list"
	repProto "owl/repeater/proto"
	"time"
)

// agent struct
type agent struct {
	cfcCli cfcProto.OwlCfcServiceClient
	repCli repProto.OwlRepeaterServiceClient

	executor *executor.Executor

	agentInfo  cfcProto.AgentInfo
	pluginList *pluginList.PluginList
	metricList *metricList.MetricList

	conf   *conf.Conf
	logger *logger.Logger

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newAgent(ctx context.Context, conf *conf.Conf, lg *logger.Logger) (*agent, error) {
	agt := &agent{
		executor: executor.NewExecutor(lg),

		pluginList: pluginList.NewPluginList(),
		metricList: metricList.NewMetricList(),

		conf:   conf,
		logger: lg,
	}

	agt.ctx, agt.cancelFunc = context.WithCancel(ctx)

	// 初始化失败时，不返回agent对象
	if err := agt.init(); err != nil {
		return nil, err
	}

	return agt, nil
}

func (agent *agent) Start() error {
	agent.logger.Info(fmt.Sprintf("Starting owl agent %s", global.Version))

	// 首次启动首先需要注册Agent
	_ = agent.registerAgent()

	reportHbTk := time.Tick(agent.conf.ReportHeartbeatIntervalSecs)
	listPluginsTk := time.Tick(agent.conf.ListPluginsIntervalSecs)
	reportAgentMetricsTk := time.Tick(agent.conf.ReportMetricIntervalSecs)
	execBuiltinMetricsTk := time.Tick(time.Duration(agent.conf.ExecBuiltinMetricsIntervalSecs) * time.Second)

	for {
		select {
		case <-reportHbTk:
			go func() {
				agent.refreshAgentInfo()
				agent.reportHeartbeat()
			}()

		case <-listPluginsTk:
			go agent.listPluginsProcess()

		case <-reportAgentMetricsTk:
			go agent.reportAgentAllMetrics()

		case <-execBuiltinMetricsTk:
			go agent.execBuiltinMetrics(int32(agent.conf.ExecBuiltinMetricsIntervalSecs))

		case <-agent.ctx.Done():
			agent.logger.InfoWithFields(logger.Fields{
				"context_error": agent.ctx.Err(),
			}, "owl agent exited by context done.")
			return agent.ctx.Err()
		}
	}
}

func (agent *agent) Stop() {
	if agent.pluginList != nil {
		agent.pluginList.StopAllPluginTask()
	}
	agent.cancelFunc()
}

func (agent *agent) refreshAgentInfo() {
	agent.agentInfo.HostId = agent.executor.GetHostID()
	agent.agentInfo.Hostname = agent.executor.GetHostname()
	agent.agentInfo.AgentVersion = global.Version
	agent.agentInfo.Uptime, agent.agentInfo.IdlePct = agent.executor.GetHostUptimeAndIdle()

	// Get local ip with cfc
	if ip := agent.executor.GetLocalIp(agent.conf.CfcAddress); len(ip) > 0 {
		agent.agentInfo.Ip = ip
		return
	}

	// Get local ip with repeater
	agent.agentInfo.Ip = agent.executor.GetLocalIp(agent.conf.RepeaterAddress)
}

func (agent *agent) init() (err error) {
	// 连接Cfc
	agent.cfcCli, err = cli.NewCfcClient(agent.conf.CfcAddress)
	if err != nil {
		return err
	}

	// 连接Repeater
	agent.repCli, err = cli.NewRepeaterClient(agent.conf.RepeaterAddress)
	if err != nil {
		return err
	}

	// 创建插件路径
	if err = os.Mkdir(agent.conf.PluginDir, 0755); err != nil {
		agent.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while os.Mkdir for create plugins dir in agent.init, Skipped it.")
	}

	agent.refreshAgentInfo()
	agent.listPluginsProcess()

	return nil
}

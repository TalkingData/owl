package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"owl/agent/conf"
	"owl/agent/executor"
	"owl/cli"
	"owl/common/global"
	"owl/common/logger"
	metricList "owl/dto/metric_list"
	pluginList "owl/dto/plugin_list"
	proxyProto "owl/proxy/proto"
	"sync"
	"time"
)

// Agent 客户端
type Agent interface {
	// Start 启动Component服务
	Start() error
	// Stop 关闭Component服务
	Stop()
}

// NewAgent 创建Agent组件
func NewAgent(ctx context.Context, conf *conf.Conf, lg *logger.Logger) (Agent, error) {
	return newAgent(ctx, conf, lg)
}

// agent struct
type agent struct {
	httpServer *http.Server

	proxyCli proxyProto.OwlProxyServiceClient

	executor *executor.Executor

	agentInfo  proxyProto.AgentInfo
	pluginList *pluginList.PluginList
	metricList *metricList.MetricList

	conf   *conf.Conf
	logger *logger.Logger

	wg         *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newAgent(ctx context.Context, conf *conf.Conf, lg *logger.Logger) (*agent, error) {
	a := &agent{
		executor: executor.NewExecutor(lg),

		pluginList: pluginList.NewPluginList(),
		metricList: metricList.NewMetricList(),

		conf:   conf,
		logger: lg,

		wg: new(sync.WaitGroup),
	}

	a.ctx, a.cancelFunc = context.WithCancel(ctx)

	// 初始化失败时，不返回agent对象
	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (agent *agent) Start() error {
	agent.logger.Info(fmt.Sprintf("Starting owl agent %s...", global.Version))

	agent.wg.Add(1)
	defer agent.wg.Done()

	// 首次启动首先需要注册Agent
	_ = agent.registerAgent()

	reportHbTk := time.Tick(agent.conf.ReportHeartbeatIntervalSecs)
	listPluginsTk := time.Tick(agent.conf.ListPluginsIntervalSecs)
	reportAgentMetricsTk := time.Tick(agent.conf.ReportMetricIntervalSecs)
	execBuiltinMetricsTk := time.Tick(
		time.Duration(agent.conf.ExecBuiltinMetricsIntervalSecs) * time.Second,
	)

	// 启动httpServer
	go func() {
		agent.wg.Add(1)
		defer agent.Stop()
		defer agent.wg.Done()

		agent.logger.Info(fmt.Sprintf("Owl agent's http server listening on: %s", agent.conf.Listen))

		if err := agent.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			agent.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while httpServer.ListenAndServe.")
			return
		}

		agent.logger.Info("Owl agent's http server closed.")
	}()

	// 启动主服务
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
			go agent.execBuiltinMetrics()

		case <-agent.ctx.Done():
			agent.logger.InfoWithFields(logger.Fields{
				"context_error": agent.ctx.Err(),
			}, "Owl agent exited by context done.")
			return agent.ctx.Err()
		}
	}
}

func (agent *agent) Stop() {
	defer agent.wg.Wait()

	// 关闭httpServer
	if agent.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), agent.conf.Const.HttpServerShutdownTimeoutSecs)
		defer cancel()
		_ = agent.httpServer.Shutdown(ctx)
	}

	// 关闭插件采集任务
	if agent.pluginList != nil {
		agent.pluginList.StopAllPluginTask()
	}

	agent.cancelFunc()
}

func (agent *agent) init() (err error) {
	agent.httpServer = &http.Server{
		Addr:    agent.conf.Listen,
		Handler: agent.newHttpHandler(),
	}

	// 连接Proxy
	agent.proxyCli, err = cli.NewProxyClient(agent.conf.ProxyAddress)
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

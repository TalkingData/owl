package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"owl/agent/conf"
	dtHandler "owl/agent/dt_handler"
	"owl/agent/executor"
	"owl/cli"
	"owl/common/global"
	"owl/common/logger"
	commonpb "owl/common/proto"
	"owl/dto"
	pluginList "owl/dto/plugin_list"
	proxypb "owl/proxy/proto"
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

	proxyCli proxypb.OwlProxyServiceClient

	executor *executor.Executor

	agentInfo  commonpb.AgentInfo
	pluginList *pluginList.PluginList
	tsDataMap  *dto.TsDataMap

	tsDataBuffer *dto.TsDataBuffer

	dtHandlerMap dtHandler.DtHandlerMap

	metricServer *http.Server

	conf   *conf.Conf
	logger *logger.Logger

	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newAgent(ctx context.Context, conf *conf.Conf, lg *logger.Logger) (*agent, error) {
	a := &agent{
		executor: executor.NewExecutor(lg),

		pluginList: pluginList.NewPluginList(),
		tsDataMap:  dto.NewTsDataMap(),

		tsDataBuffer: dto.NewTsDataBuffer(),

		dtHandlerMap: dtHandler.NewDtHandlerMap(),

		conf:   conf,
		logger: lg,
	}

	a.ctx, a.cancelFunc = context.WithCancel(ctx)

	// 初始化失败时，不返回agent对象
	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *agent) Start() error {
	a.logger.InfoWithFields(logger.Fields{
		"branch":  global.Branch,
		"commit":  global.Commit,
		"version": global.Version,
	}, "Starting owl agent...")

	a.wg.Add(1)
	defer a.wg.Done()

	// 首次启动首先需要注册Agent
	_ = a.registerAgent()

	reportHbTk := time.Tick(a.conf.ReportHeartbeatIntervalSecs)
	listPluginsTk := time.Tick(a.conf.ListPluginsIntervalSecs)
	reportAgentMetricsTk := time.Tick(a.conf.ReportMetricsIntervalSecs)
	execBuiltinMetricsTk := time.Tick(
		time.Duration(a.conf.ExecBuiltinMetricsIntervalSecs) * time.Second,
	)
	forceSendTsDataTk := time.Tick(a.conf.ForceSendTsDataIntervalSecs)

	// 启动httpServer
	go func() {
		a.wg.Add(1)
		defer a.Stop()
		defer a.wg.Done()

		a.logger.Info(fmt.Sprintf("Owl agent's http server listening on: %s", a.conf.Listen))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling httpServer.ListenAndServe.")
			return
		}

		a.logger.Info("Owl agent's http server closed.")
	}()

	// 启动Prometheus的metrics http server
	go func() {
		a.wg.Add(1)
		defer a.Stop()
		defer a.wg.Done()

		a.logger.Info(fmt.Sprintf(
			"Owl agent's metrics http server listening on: %s", a.conf.MetricListen,
		))
		if err := a.metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.ErrorWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling metricServer.ListenAndServe.")
			return
		}

		a.logger.Info("Owl agent's metrics server closed.")
	}()

	// 启动主服务
	for {
		select {
		case <-reportHbTk:
			go func() {
				a.refreshAgentInfo()
				a.reportHeartbeat()
			}()

		case <-listPluginsTk:
			go a.listPluginsProcess()

		case <-reportAgentMetricsTk:
			go a.reportAgentAllMetrics()

		case c := <-execBuiltinMetricsTk:
			go a.execBuiltinMetrics(c.Unix())

		case <-forceSendTsDataTk:
			go a.sendTsData(true)

		case <-a.ctx.Done():
			a.logger.InfoWithFields(logger.Fields{
				"context_error": a.ctx.Err(),
			}, "Owl agent exited by context done.")
			return a.ctx.Err()
		}
	}
}

func (a *agent) Stop() {
	defer a.wg.Wait()

	// 结束前强制发送数据定时器
	a.sendTsData(true)

	// 关闭httpServer
	if a.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), a.conf.Const.HttpServerShutdownTimeoutSecs)
		defer cancel()
		_ = a.httpServer.Shutdown(ctx)
	}

	// 关闭插件采集任务
	if a.pluginList != nil {
		a.pluginList.StopAllPluginTask()
	}

	// 关闭metrics http server
	if a.metricServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), a.conf.Const.HttpServerShutdownTimeoutSecs)
		defer cancel()
		_ = a.metricServer.Shutdown(ctx)
	}

	a.cancelFunc()
}

func (a *agent) init() (err error) {
	a.httpServer = &http.Server{
		Addr:    a.conf.Listen,
		Handler: a.newHttpHandler(),
	}

	a.metricServer = &http.Server{
		Addr:    a.conf.MetricListen,
		Handler: a.newMetricHttpHandler(),
	}

	// 连接Proxy
	a.proxyCli, err = cli.NewProxyClient(a.conf.ProxyAddress)
	if err != nil {
		return err
	}

	// 创建插件路径
	if err = os.Mkdir(a.conf.PluginDir, 0755); err != nil {
		a.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling os.Mkdir for create plugins dir, Skipped it.")
	}

	a.refreshAgentInfo()
	a.listPluginsProcess()

	return nil
}

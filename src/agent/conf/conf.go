package conf

import (
	"github.com/Unknwon/goconfig"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                         string
	MetricListen                   string
	ExecBuiltinMetricsIntervalSecs int

	LogPath  string
	LogLevel string

	ProxyAddress string

	ListPluginsTimeoutSecs    time.Duration
	DownloadPluginTimeoutSecs time.Duration
	ListPluginsIntervalSecs   time.Duration

	CleanExpiredMetricCycleExpiredRatio int

	ReportMetricTimeoutSecs   time.Duration
	ReportMetricsTimeoutSecs  time.Duration
	ReportMetricsIntervalSecs time.Duration
	ReportMetricBatchSize     int

	RegisterAgentTimeoutSecs    time.Duration
	ReportHeartbeatTimeoutSecs  time.Duration
	ReportHeartbeatIntervalSecs time.Duration

	SendTsDataTimeoutSecs       time.Duration
	SendTsDataArrayTimeoutSecs  time.Duration
	SendTsDataBatchSize         int
	ForceSendTsDataIntervalSecs time.Duration

	PluginDir              string
	ExecuteUntrustedPlugin bool
}

func NewConfig(options ...Option) *Conf {
	opts := newOptions(options...)

	cfg, err := goconfig.LoadConfigFile(opts.ConfFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen:       cfg.MustValue("main", "listen", defaultListen),
		MetricListen: cfg.MustValue("main", "metric_listen", defaultMetricListen),
		ExecBuiltinMetricsIntervalSecs: cfg.MustInt(
			"main", "exec_builtin_metrics_interval_secs", defaultExecBuiltinMetricsIntervalSecs,
		),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		ProxyAddress: cfg.MustValue("proxy", "address", defaultProxyAddress),

		ListPluginsTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "list_plugin_timeout_secs", defaultListPluginTimeoutSecs,
		)) * time.Second,
		DownloadPluginTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "download_plugin_timeout_secs", defaultDownloadPluginTimeoutSecs,
		)) * time.Second,
		ListPluginsIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "list_plugins_interval_secs", defaultListPluginsIntervalSecs,
		)) * time.Second,

		CleanExpiredMetricCycleExpiredRatio: cfg.MustInt(
			"main", "clean_expired_metric_cycle_expired_ratio", defaultCleanExpiredMetricCycleExpiredRatio,
		),

		ReportMetricTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "report_metric_timeout_secs", defaultReportMetricTimeoutSecs,
		)) * time.Second,
		ReportMetricsTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "report_metrics_timeout_secs", defaultReportMetricsTimeoutSecs,
		)) * time.Second,
		ReportMetricsIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "report_metrics_interval_secs", defaultReportMetricsIntervalSecs,
		)) * time.Second,
		ReportMetricBatchSize: cfg.MustInt(
			"proxy", "report_metric_batch_size", defaultReportMetricBatchSize,
		),

		RegisterAgentTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "register_agent_timeout_secs", defaultRegisterAgentTimeoutSecs,
		)) * time.Second,
		ReportHeartbeatTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "report_heartbeat_timeout_secs", defaultReportHeartbeatTimeoutSecs,
		)) * time.Second,
		ReportHeartbeatIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "report_heartbeat_interval_secs", defaultReportHeartbeatIntervalSecs,
		)) * time.Second,

		SendTsDataTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "send_ts_data_timeout_secs", defaultSendTsDataTimeoutSecs,
		)) * time.Second,
		SendTsDataArrayTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "send_ts_data_array_timeout_secs", defaultSendTsDataArrayTimeoutSecs,
		)) * time.Second,
		SendTsDataBatchSize: cfg.MustInt(
			"proxy", "send_ts_data_batch_size", defaultSendTsDataBatchSize,
		),
		ForceSendTsDataIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "force_send_ts_data_interval_secs", defaultForceSendTsDataIntervalSecs,
		)) * time.Second,

		PluginDir: cfg.MustValue("plugin", "plugin_dir", defaultPluginDir),
		ExecuteUntrustedPlugin: cfg.MustBool(
			"plugin", "execute_untrusted_plugin", defaultExecuteUntrustedPlugin,
		),
	}
}

package conf

import (
	"github.com/Unknwon/goconfig"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                         string
	ExecBuiltinMetricsIntervalSecs int

	LogPath  string
	LogLevel string

	ProxyAddress                string
	CallProxyTimeoutSecs        time.Duration
	DownloadPluginTimeoutSecs   time.Duration
	ListPluginsIntervalSecs     time.Duration
	ReportMetricIntervalSecs    time.Duration
	ReportMetricBatchSize       int
	ReportHeartbeatIntervalSecs time.Duration
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

		Listen: cfg.MustValue("main", "listen", defaultListen),
		ExecBuiltinMetricsIntervalSecs: cfg.MustInt(
			"main", "exec_builtin_metrics_interval_secs", defaultExecBuiltinMetricsIntervalSecs,
		),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		ProxyAddress: cfg.MustValue("proxy", "address", defaultProxyAddress),
		CallProxyTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "call_proxy_timeout_secs", defaultCallProxyTimeoutSecs,
		)) * time.Second,
		DownloadPluginTimeoutSecs: time.Duration(cfg.MustInt(
			"proxy", "download_plugin_timeout_secs", defaultDownloadPluginTimeoutSecs,
		)) * time.Second,
		ListPluginsIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "list_plugins_interval_secs", defaultListPluginsIntervalSecs,
		)) * time.Second,
		ReportMetricIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "report_metric_interval_secs", defaultReportMetricIntervalSecs,
		)) * time.Second,
		ReportMetricBatchSize: cfg.MustInt(
			"proxy", "report_metric_batch_size", defaultReportMetricBatchSize,
		),
		ReportHeartbeatIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "report_heartbeat_interval_secs", defaultReportHeartbeatIntervalSecs,
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

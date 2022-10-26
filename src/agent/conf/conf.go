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
	ReportHeartbeatIntervalSecs time.Duration

	PluginDir              string
	ExecuteUntrustedPlugin bool
}

func NewConfig() *Conf {
	cfg, err := goconfig.LoadConfigFile(defaultConfigFilePathname)
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
		ReportHeartbeatIntervalSecs: time.Duration(cfg.MustInt(
			"proxy", "report_heartbeat_interval_secs", defaultReportHeartbeatIntervalSecs,
		)) * time.Second,

		PluginDir: cfg.MustValue("plugin", "plugin_dir", defaultPluginDir),
		ExecuteUntrustedPlugin: cfg.MustBool(
			"plugin", "execute_untrusted_plugin", defaultExecuteUntrustedPlugin,
		),
	}
}

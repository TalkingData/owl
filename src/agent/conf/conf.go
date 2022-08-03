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

	CfcAddress                  string
	CallCfcTimeoutSecs          time.Duration
	DownloadPluginTimeoutSecs   time.Duration
	ListPluginsIntervalSecs     time.Duration
	ReportMetricIntervalSecs    time.Duration
	ReportHeartbeatIntervalSecs time.Duration

	RepeaterAddress         string
	CallRepeaterTimeoutSecs time.Duration

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

		LogLevel: cfg.MustValue("log", "log_level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "log_path", defaultLogPath),

		CfcAddress: cfg.MustValue("cfc", "address", defaultCfcAddress),
		CallCfcTimeoutSecs: time.Duration(cfg.MustInt(
			"cfc", "call_cfc_timeout_secs", defaultCallCfcTimeoutSecs,
		)) * time.Second,
		DownloadPluginTimeoutSecs: time.Duration(cfg.MustInt(
			"cfc", "download_plugin_timeout_secs", defaultDownloadPluginTimeoutSecs,
		)) * time.Second,
		ListPluginsIntervalSecs: time.Duration(cfg.MustInt(
			"cfc", "list_plugins_interval_secs", defaultListPluginsIntervalSecs,
		)) * time.Second,
		ReportMetricIntervalSecs: time.Duration(cfg.MustInt(
			"cfc", "report_metric_interval_secs", defaultReportMetricIntervalSecs,
		)) * time.Second,
		ReportHeartbeatIntervalSecs: time.Duration(cfg.MustInt(
			"cfc", "report_heartbeat_interval_secs", defaultReportHeartbeatIntervalSecs,
		)) * time.Second,

		RepeaterAddress: cfg.MustValue("repeater", "address", defaultRepeaterAddress),
		CallRepeaterTimeoutSecs: time.Duration(cfg.MustInt(
			"repeater", "call_repeater_timeout_secs", defaultCallRepeaterTimeoutSecs,
		)) * time.Second,

		PluginDir: cfg.MustValue("plugin", "plugin_dir", defaultPluginDir),
		ExecuteUntrustedPlugin: cfg.MustBool(
			"plugin", "execute_untrusted_plugin", defaultExecuteUntrustedPlugin,
		),
	}
}

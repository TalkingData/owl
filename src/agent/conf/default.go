package conf

const (
	defaultConfigFilePathname = "../conf/owl_agent.conf"
)

const (
	defaultListen                         = "127.0.0.1:19575"
	defaultExecBuiltinMetricsIntervalSecs = 60

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	defaultCfcAddress                  = "127.0.0.1:19576"
	defaultCallCfcTimeoutSecs          = 10
	defaultDownloadPluginTimeoutSecs   = 90
	defaultListPluginsIntervalSecs     = 300
	defaultReportMetricIntervalSecs    = 300
	defaultReportHeartbeatIntervalSecs = 60

	defaultRepeaterAddress         = "127.0.0.1:19577"
	defaultCallRepeaterTimeoutSecs = 10

	defaultPluginDir              = "../plugins"
	defaultExecuteUntrustedPlugin = false
)

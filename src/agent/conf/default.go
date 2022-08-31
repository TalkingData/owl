package conf

const (
	defaultConfigFilePathname = "../conf/owl_agent.conf"
)

const (
	defaultListen                         = "127.0.0.1:19575"
	defaultExecBuiltinMetricsIntervalSecs = 60

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	defaultProxyAddress                = "127.0.0.1:19577"
	defaultCallProxyTimeoutSecs        = 10
	defaultDownloadPluginTimeoutSecs   = 90
	defaultListPluginsIntervalSecs     = 300
	defaultReportMetricIntervalSecs    = 300
	defaultReportHeartbeatIntervalSecs = 60

	defaultPluginDir              = "../plugins"
	defaultExecuteUntrustedPlugin = false
)

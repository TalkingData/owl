package conf

const (
	defaultListen                         = "127.0.0.1:19575"
	defaultMetricListen                   = "127.0.0.1:19672"
	defaultExecBuiltinMetricsIntervalSecs = 60

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	defaultProxyAddress = "127.0.0.1:19570"

	defaultListPluginTimeoutSecs     = 300
	defaultDownloadPluginTimeoutSecs = 300
	defaultListPluginsIntervalSecs   = 300

	defaultCleanExpiredMetricCycleExpiredRatio = 60

	defaultReportMetricTimeoutSecs   = 90
	defaultReportMetricsTimeoutSecs  = 300
	defaultReportMetricsIntervalSecs = 300
	defaultReportMetricBatchSize     = 500

	defaultRegisterAgentTimeoutSecs    = 90
	defaultReportHeartbeatTimeoutSecs  = 90
	defaultReportHeartbeatIntervalSecs = 60

	defaultSendTsDataTimeoutSecs       = 90
	defaultSendTsDataArrayTimeoutSecs  = 120
	defaultSendTsDataBatchSize         = 200
	defaultForceSendTsDataIntervalSecs = 25

	defaultPluginDir              = "../plugins"
	defaultExecuteUntrustedPlugin = false
)

package conf

const (
	defaultListen                         = "127.0.0.1:19575"
	defaultExecBuiltinMetricsIntervalSecs = 60

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	defaultProxyAddress                = "127.0.0.1:19570"
	defaultCallProxyTimeoutSecs        = 10
	defaultDownloadPluginTimeoutSecs   = 90
	defaultListPluginsIntervalSecs     = 300
	defaultReportMetricIntervalSecs    = 300
	defaultReportMetricBatchSize       = 500
	defaultReportHeartbeatIntervalSecs = 60
	defaultSendTsDataBatchSize         = 50
	defaultForceSendTsDataIntervalSecs = 20

	defaultPluginDir              = "../plugins"
	defaultExecuteUntrustedPlugin = false
)

package conf

const (
	defaultConfigFilePathname = "../conf/owl_cfc.conf"
)

const (
	defaultListen                              = "127.0.0.1:19576"
	defaultPluginDir                           = "../plugins"
	defaultRefreshHostStatusIntervalSecs       = 90
	defaultHostDownStatusThresholdSecs         = 90
	defaultCleanExpiredMetricIntervalSecs      = 300
	defaultCleanExpiredMetricCycleExpiredRatio = 5
	defaultAllowCreateProductAuto              = true

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	// Mysql默认配置
	defaultMysqlAddress      = "127.0.0.1:3306"
	defaultMysqlDbName       = "owl_v5"
	defaultMysqlUser         = "owl"
	defaultMysqlPassword     = "owl"
	defaultMysqlMaxIdleConns = 20
	defaultMysqlMaxOpenConns = 500
)

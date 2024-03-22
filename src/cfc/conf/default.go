package conf

const (
	defaultListen                              = "127.0.0.1:0"
	defaultMetricListen                        = "127.0.0.1:19671"
	defaultMicroRegisterTtl                    = 10
	defaultMicroRegisterInterval               = 3
	defaultRefreshHostStatusIntervalSecs       = 90
	defaultHostDownStatusThresholdSecs         = 90
	defaultCleanExpiredMetricIntervalSecs      = 300
	defaultCleanExpiredMetricCycleExpiredRatio = 60
	defaultAllowCreateProductAuto              = true

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	// Mysql默认配置
	defaultMysqlAddress      = "127.0.0.1:3306"
	defaultMysqlDbName       = "owl_v5"
	defaultMysqlUser         = "owl"
	defaultMysqlPassword     = "owl"
	defaultMysqlMaxIdleConns = 20
	defaultMysqlMaxOpenConns = 500
)

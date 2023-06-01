package conf

const (
	defaultListen                = "127.0.0.1:0"
	defaultMetricListen          = "127.0.0.1:19672"
	defaultMicroRegisterTtl      = 10
	defaultMicroRegisterInterval = 3
	defaultBackend               = "kairosdb"

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	defaultKairosDbAddress      = "127.0.0.1:4242"
	defaultKairosDbMaxIdleConns = 10
	defaultKairosDbMaxOpenConns = 20

	defaultKairosDbRedundantAddress1     = "127.0.0.1:4242"
	defaultKairosDbRedundantAddress2     = "127.0.0.1:4242"
	defaultKairosDbRedundantMaxIdleConns = 10
	defaultKairosDbRedundantMaxOpenConns = 20
)

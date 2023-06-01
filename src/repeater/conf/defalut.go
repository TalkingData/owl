package conf

const (
	defaultListen                = "127.0.0.1:0"
	defaultMicroRegisterTtl      = 10
	defaultMicroRegisterInterval = 3
	defaultBackend               = "kairosdb"

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	// Etcd默认配置
	defaultEtcdUsername = ""
	defaultEtcdPassword = ""

	defaultKairosDbAddress      = "127.0.0.1:4242"
	defaultKairosDbMaxIdleConns = 30
	defaultKairosDbMaxOpenConns = 100
)

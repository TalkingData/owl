package conf

const (
	defaultConfigFilePathname = "../conf/owl_repeater.conf"
)

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

	defaultKairosDbAddress = "127.0.0.1:4242"
)

package conf

const (
	defaultConfigFilePathname = "../conf/owl_repeater.conf"
)

const (
	defaultListen = "127.0.0.1:19577"

	defaultLogLevel = "debug"
	defaultLogPath  = "../logs"

	defaultBackend             = "opentsdb"
	defaultOpentsdbAddress     = "127.0.0.1:4242"
	defaultKairosdbAddress     = "127.0.0.1:4242"
	defaultKairosdbRestAddress = "127.0.0.1:8080"
	defaultKfkAddresses        = "127.0.0.1:9092,127.0.0.1:9092,127.0.0.1:9092"
	defaultKfkTopic            = "owl"
)

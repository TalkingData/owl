package conf

import (
	"github.com/Unknwon/goconfig"
)

type Conf struct {
	Const *constConf

	Listen string

	LogPath  string
	LogLevel string

	Backend             string
	OpentsdbAddress     string
	KairosdbAddress     string
	KairosdbRestAddress string
	KafkaAddresses      []string
	KafkaTopic          string
}

func NewConfig() *Conf {
	cfg, err := goconfig.LoadConfigFile(defaultConfigFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen: cfg.MustValue("main", "listen", defaultListen),

		LogLevel: cfg.MustValue("log", "log_level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "log_path", defaultLogPath),

		Backend:             cfg.MustValue("backend", "backend", defaultBackend),
		OpentsdbAddress:     cfg.MustValue("backend", "opentsdb_address", defaultOpentsdbAddress),
		KairosdbAddress:     cfg.MustValue("backend", "kairosdb_address", defaultKairosdbAddress),
		KairosdbRestAddress: cfg.MustValue("backend", "rest_kairosdb_address", defaultKairosdbRestAddress),
		KafkaAddresses:      cfg.MustValueArray("backend", "kafka_address", defaultKfkAddresses),
		KafkaTopic:          cfg.MustValue("backend", "kafka_topic", defaultKfkTopic),
	}
}

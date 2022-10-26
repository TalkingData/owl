package conf

import (
	"github.com/Unknwon/goconfig"
	"owl/common/global"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                string
	MicroRegisterTtl      time.Duration
	MicroRegisterInterval time.Duration
	Backend               string

	LogPath  string
	LogLevel string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	KairosDbAddress string
}

func NewConfig() *Conf {
	cfg, err := goconfig.LoadConfigFile(defaultConfigFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen: cfg.MustValue("main", "listen", defaultListen),
		MicroRegisterTtl: time.Duration(cfg.MustInt(
			"main", "micro_register_ttl", defaultMicroRegisterTtl,
		)) * time.Second,
		MicroRegisterInterval: time.Duration(cfg.MustInt(
			"main", "micro_register_interval", defaultMicroRegisterInterval,
		)) * time.Second,
		Backend: cfg.MustValue("main", "backend", defaultBackend),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),

		KairosDbAddress: cfg.MustValue("kairosdb", "kairosdb_address", defaultKairosDbAddress),
	}
}

package conf

import (
	"github.com/Unknwon/goconfig"
)

type Conf struct {
	Const *constConf

	Listen  string
	Backend string

	LogPath  string
	LogLevel string

	KairosDbAddress string
}

func NewConfig() *Conf {
	cfg, err := goconfig.LoadConfigFile(defaultConfigFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen:  cfg.MustValue("main", "listen", defaultListen),
		Backend: cfg.MustValue("main", "backend", defaultBackend),

		LogLevel: cfg.MustValue("log", "log_level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "log_path", defaultLogPath),

		KairosDbAddress: cfg.MustValue("kairosdb", "kairosdb_address", defaultKairosDbAddress),
	}
}

package conf

import (
	"github.com/Unknwon/goconfig"
	"owl/common/global"
)

type Conf struct {
	Const *constConf

	Listen    string
	PluginDir string

	LogPath  string
	LogLevel string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string
}

func NewConfig() *Conf {
	cfg, err := goconfig.LoadConfigFile(defaultConfigFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen:    cfg.MustValue("main", "listen", defaultListen),
		PluginDir: cfg.MustValue("main", "plugin_dir", defaultPluginDir),

		LogLevel: cfg.MustValue("log", "log_level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "log_path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),
	}
}

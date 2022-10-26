package conf

import (
	"github.com/Unknwon/goconfig"
	"owl/common/global"
)

type Conf struct {
	Const *constConf

	Listen              string
	CallCfcRetries      int
	CallRepeaterRetries int
	PluginDir           string

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

		Listen:              cfg.MustValue("main", "listen", defaultListen),
		CallCfcRetries:      cfg.MustInt("main", "call_cfc_retries", defaultCallCfcRetries),
		CallRepeaterRetries: cfg.MustInt("main", "call_repeater_retries", defaultCallRepeaterRetries),
		PluginDir:           cfg.MustValue("main", "plugin_dir", defaultPluginDir),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),
	}
}

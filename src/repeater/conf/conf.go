package conf

import (
	"github.com/Unknwon/goconfig"
	"owl/common/global"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                string
	MetricListen          string
	MicroRegisterTtl      time.Duration
	MicroRegisterInterval time.Duration
	Backend               string

	LogPath  string
	LogLevel string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	KairosDbAddress      string
	KairosDbMaxIdleConns int
	KairosDbMaxOpenConns int

	KairosDbRedundantAddress1     string
	KairosDbRedundantAddress2     string
	KairosDbRedundantMaxIdleConns int
	KairosDbRedundantMaxOpenConns int
}

func NewConfig(options ...Option) *Conf {
	opts := newOptions(options...)

	cfg, err := goconfig.LoadConfigFile(opts.ConfFilePathname)
	if err != nil {
		panic(err)
	}

	return &Conf{
		Const: newConstConf(),

		Listen:       cfg.MustValue("main", "listen", defaultListen),
		MetricListen: cfg.MustValue("main", "metric_listen", defaultMetricListen),
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

		KairosDbAddress:      cfg.MustValue("kairosdb", "kairosdb_address", defaultKairosDbAddress),
		KairosDbMaxIdleConns: cfg.MustInt("kairosdb", "kairosdb_max_idle_conns", defaultKairosDbMaxIdleConns),
		KairosDbMaxOpenConns: cfg.MustInt("kairosdb", "kairosdb_max_open_conns", defaultKairosDbMaxOpenConns),

		KairosDbRedundantAddress1: cfg.MustValue(
			"kairosdb_redundant", "kairosdb_redundant_address1", defaultKairosDbRedundantAddress1,
		),
		KairosDbRedundantAddress2: cfg.MustValue(
			"kairosdb_redundant", "kairosdb_redundant_address2", defaultKairosDbRedundantAddress2,
		),
		KairosDbRedundantMaxIdleConns: cfg.MustInt(
			"kairosdb_redundant", "kairosdb_redundant_max_idle_conns", defaultKairosDbRedundantMaxIdleConns,
		),
		KairosDbRedundantMaxOpenConns: cfg.MustInt(
			"kairosdb_redundant", "kairosdb_redundant_max_open_conns", defaultKairosDbRedundantMaxOpenConns,
		),
	}
}

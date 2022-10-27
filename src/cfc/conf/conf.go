package conf

import (
	"github.com/Unknwon/goconfig"
	"owl/common/global"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                              string
	MicroRegisterTtl                    time.Duration
	MicroRegisterInterval               time.Duration
	RefreshHostStatusIntervalSecs       time.Duration
	HostDownStatusThresholdSecs         int
	CleanExpiredMetricIntervalSecs      time.Duration
	CleanExpiredMetricCycleExpiredRatio int
	AllowCreateProductAuto              bool

	LogPath  string
	LogLevel string

	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	MysqlAddress      string
	MysqlDbName       string
	MysqlUser         string
	MysqlPassword     string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int
}

func NewConfig(options ...Option) *Conf {
	opts := newOptions(options...)

	cfg, err := goconfig.LoadConfigFile(opts.ConfFilePathname)
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
		RefreshHostStatusIntervalSecs: time.Duration(cfg.MustInt(
			"main", "refresh_host_status_interval_secs", defaultRefreshHostStatusIntervalSecs,
		)) * time.Second,
		HostDownStatusThresholdSecs: cfg.MustInt(
			"main", "host_down_status_threshold_secs", defaultHostDownStatusThresholdSecs,
		),
		CleanExpiredMetricIntervalSecs: time.Duration(cfg.MustInt(
			"main", "clean_expired_metric_interval_secs", defaultCleanExpiredMetricIntervalSecs,
		)) * time.Second,
		CleanExpiredMetricCycleExpiredRatio: cfg.MustInt(
			"main", "clean_expired_metric_cycle_expired_ratio", defaultCleanExpiredMetricCycleExpiredRatio,
		),
		AllowCreateProductAuto: cfg.MustBool(
			"main", "allow_create_product_auto", defaultAllowCreateProductAuto,
		),

		LogLevel: cfg.MustValue("log", "level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "path", defaultLogPath),

		EtcdAddresses: cfg.MustValueArray("etcd", "addresses", global.DefaultConfigSeparator),
		EtcdUsername:  cfg.MustValue("etcd", "username", defaultEtcdUsername),
		EtcdPassword:  cfg.MustValue("etcd", "password", defaultEtcdPassword),

		MysqlAddress:      cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:       cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:         cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:     cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxIdleConns: cfg.MustInt("mysql", "max_idle_conns", defaultMysqlMaxIdleConns),
		MysqlMaxOpenConns: cfg.MustInt("mysql", "max_open_conns", defaultMysqlMaxOpenConns),
	}
}

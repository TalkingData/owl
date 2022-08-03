package conf

import (
	"github.com/Unknwon/goconfig"
	"time"
)

type Conf struct {
	Const *constConf

	Listen                              string
	PluginDir                           string
	RefreshHostStatusIntervalSecs       time.Duration
	HostDownStatusThresholdSecs         int
	CleanExpiredMetricIntervalSecs      time.Duration
	CleanExpiredMetricCycleExpiredRatio int
	AllowCreateProductAuto              bool

	LogPath  string
	LogLevel string

	MysqlAddress      string
	MysqlDbName       string
	MysqlUser         string
	MysqlPassword     string
	MysqlMaxIdleConns int
	MysqlMaxOpenConns int
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

		LogLevel: cfg.MustValue("log", "log_level", defaultLogLevel),
		LogPath:  cfg.MustValue("log", "log_path", defaultLogPath),

		MysqlAddress:      cfg.MustValue("mysql", "address", defaultMysqlAddress),
		MysqlDbName:       cfg.MustValue("mysql", "db_name", defaultMysqlDbName),
		MysqlUser:         cfg.MustValue("mysql", "user", defaultMysqlUser),
		MysqlPassword:     cfg.MustValue("mysql", "password", defaultMysqlPassword),
		MysqlMaxIdleConns: cfg.MustInt("mysql", "max_idle_conns", defaultMysqlMaxIdleConns),
		MysqlMaxOpenConns: cfg.MustInt("mysql", "max_open_conns", defaultMysqlMaxOpenConns),
	}
}

package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath                       = "./conf/cfc.conf"
	DefaultTCPBind                       = "0.0.0.0:10020"
	DefaultMetricBind                    = "0.0.0.0:10021"
	DefaultMySQLAddr                     = "127.0.0.1:3306"
	DefaultMySQLDBName                   = "owl"
	DefaultMySQLUser                     = "owl"
	DefaultMySQLPassword                 = ""
	DefaultMySQLMaxConn                  = 20
	DefaultMySQLMaxIdleConn              = 5
	DefaultMaxPacketSize                 = 4096
	DefaultLogFile                       = "./logs/cfc.log"
	DefaultLogExpireDays                 = 7
	DefaultLogLevel                      = 3
	DefaultPluginDir                     = "./plugins"
	DefaultCleanupExpiredMetricsInterval = 10
	DefaultMetricExpiredCycle            = 10
)

var GlobalConfig *Config

type Config struct {
	//MYSQL CONFIG
	MySQLAddr        string //mysql ip地址
	MySQLUser        string //mysql 登陆用户名
	MySQLPassword    string //mysql 登陆密码
	MySQLDBName      string //mysql 数据库名称
	MySQLMaxIdleConn int    //mysql 最大空闲连接数
	MySQLMaxConn     int    //mysql 最大连接数

	//SERVER CONFIG
	TCPBind    string //tcp监听地址和端口
	MetricBind string

	//LOG CONFIG
	LogFile       string //日志保存目录
	LogLevel      int    //日志级别
	LogExpireDays int    //日志保留天数

	MaxPacketSize                       int
	PluginDir                           string
	CleanupExpiredMetricIntervalMinutes int // 清理过期 metric 时间间隔，单位分钟
	MetricExpiredCycle                  int //metric 过期时间
	EnableCleanupExpiredMetric          bool
	AutoCreateProduct                   bool
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		TCPBind:          cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DefaultTCPBind),
		MetricBind:       cfg.MustValue(goconfig.DEFAULT_SECTION, "metric_bind", DefaultMetricBind),
		MySQLAddr:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_addr", DefaultMySQLAddr),
		MySQLUser:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", DefaultMySQLUser),
		MySQLDBName:      cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", DefaultMySQLDBName),
		MySQLPassword:    cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", DefaultMySQLPassword),
		MySQLMaxConn:     cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", DefaultMySQLMaxConn),
		MySQLMaxIdleConn: cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", DefaultMySQLMaxIdleConn),
		MaxPacketSize:    cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DefaultMaxPacketSize),
		LogFile:          cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DefaultLogFile),
		LogExpireDays:    cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DefaultLogExpireDays),
		LogLevel:         cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DefaultLogLevel),
		PluginDir:        cfg.MustValue(goconfig.DEFAULT_SECTION, "plugin_dir", DefaultPluginDir),
		CleanupExpiredMetricIntervalMinutes: cfg.MustInt(goconfig.DEFAULT_SECTION, "cleanup_expired_metric_interval",
			DefaultCleanupExpiredMetricsInterval),
		MetricExpiredCycle:         cfg.MustInt(goconfig.DEFAULT_SECTION, "metric_expired_cycle", DefaultMetricExpiredCycle),
		EnableCleanupExpiredMetric: cfg.MustBool(goconfig.DEFAULT_SECTION, "enable_cleanup_expired_metric", false),
		AutoCreateProduct:          cfg.MustBool(goconfig.DEFAULT_SECTION, "auto_create_product", false),
	}
	return nil
}

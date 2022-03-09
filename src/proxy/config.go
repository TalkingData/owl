package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath       = "./conf/proxy.conf"
	DefaultTCPBind       = "0.0.0.0:10030"
	DefaultMetricBind    = "0.0.0.0:10031"
	DefaultCFCAddr       = "127.0.0.1:10020"
	DefaultMaxPacketSize = 4096
	DefaultLogFile       = "./logs/proxy.log"
	DefaultLogExipreDays = 7
	DefaultLogLevel      = 7
	DefaultPluginDir     = "./plugins"
)

var GlobalConfig *Config

type Config struct {
	CFCAddr string

	TCPBind    string //tcp监听地址和端口
	MetricBind string

	//LOG CONFIG
	LogFile       string //日志保存目录
	LogLevel      int    //日志级别
	LogExpireDays int    //日志保留天数

	MaxPacketSize int
	PluginDir     string
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		TCPBind:    cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DefaultTCPBind),
		CFCAddr:    cfg.MustValue(goconfig.DEFAULT_SECTION, "cfc_addr", DefaultCFCAddr),
		MetricBind: cfg.MustValue(goconfig.DEFAULT_SECTION, "metric_bind", DefaultMetricBind),

		MaxPacketSize: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DefaultMaxPacketSize),
		LogFile:       cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DefaultLogFile),
		LogExpireDays: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DefaultLogExipreDays),
		LogLevel:      cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DefaultLogLevel),
		PluginDir:     cfg.MustValue(goconfig.DEFAULT_SECTION, "plugin_dir", DefaultPluginDir),
	}
	return nil

}

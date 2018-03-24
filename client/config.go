package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath      = "./conf/client.conf"
	DefaultCfcAddr      = "127.0.0.1:10020"
	DefaultRepeaterAddr = "127.0.0.1:10040"
	DefaultTCPBind      = "127.0.0.1:10000"
	DefaultMetricBind   = "127.0.0.1:10001"

	// 0 代表不限制
	DefaultMaxPacketSize = 0
	DefaultLogExpireDays = 7

	DefaultLogFile    = "./logs/client.log"
	DefaultLogLevel   = 3
	DefaultBufferSize = 1024 * 1024
	DefaultPluginDir  = "/usr/local/owl-plugin"
)

var GlobalConfig *Config

type Config struct {
	CfcAddr string //tcp监听地址和端口

	RepeaterAddr string
	TCPBind      string
	MetricBind   string
	BufferSize   int

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
		CfcAddr:      cfg.MustValue(goconfig.DEFAULT_SECTION, "cfc_addr", DefaultCfcAddr),
		RepeaterAddr: cfg.MustValue(goconfig.DEFAULT_SECTION, "repeater_addr", DefaultRepeaterAddr),
		TCPBind:      cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DefaultTCPBind),
		MetricBind:   cfg.MustValue(goconfig.DEFAULT_SECTION, "metric_bind", DefaultMetricBind),
		BufferSize:   cfg.MustInt(goconfig.DEFAULT_SECTION, "buffer_size", DefaultBufferSize),

		MaxPacketSize: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DefaultMaxPacketSize),
		LogFile:       cfg.MustValue(goconfig.DEFAULT_SECTION, "log_FILE", DefaultLogFile),
		LogExpireDays: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DefaultLogExpireDays),
		LogLevel:      cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DefaultLogLevel),
		PluginDir:     cfg.MustValue(goconfig.DEFAULT_SECTION, "plugin_dir", DefaultPluginDir),
	}
	return nil

}

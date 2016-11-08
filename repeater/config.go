package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath          = "./conf/repeater.conf"
	DEFAULT_BACKEND         = "opentsdb"
	DEFAULT_TCP_BIND        = "0.0.0.0:10040"
	DEFAULT_OPENTSDB_ADDR   = "127.0.0.1:4242"
	DEFAULT_REPEATER_ADDR   = ""
	DEFAULT_MAX_PACKET_SIZE = 4096
	DEFAULT_BUFFER_SIZE     = 1 << 20
	DEFAULT_LOG_FILE        = "./logs/repeater.log"
	DEFAULT_LOG_EXPIRE_DAYS = 7
	DEFAULT_LOG_LEVEL       = 3
)

var GlobalConfig *Config

type Config struct {
	BACKEND  string
	TCP_BIND string

	OPENTSDB_ADDR string

	REPEATER_ADDR string
	//LOG CONFIG
	LOG_FILE        string //日志保存目录
	LOG_LEVEL       int    //日志级别
	LOG_EXPIRE_DAYS int    //日志保留天数

	MAX_PACKET_SIZE int
	BUFFER_SIZE     int64
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		BACKEND:       cfg.MustValue(goconfig.DEFAULT_SECTION, "backend", DEFAULT_BACKEND),
		TCP_BIND:      cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DEFAULT_TCP_BIND),
		OPENTSDB_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "opentsdb_addr", DEFAULT_OPENTSDB_ADDR),
		REPEATER_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "repeater_addr", DEFAULT_REPEATER_ADDR),

		MAX_PACKET_SIZE: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		BUFFER_SIZE:     cfg.MustInt64(goconfig.DEFAULT_SECTION, "buffer_size", DEFAULT_BUFFER_SIZE),
		LOG_FILE:        cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_FILE),
		LOG_EXPIRE_DAYS: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:       cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
	}
	return nil

}

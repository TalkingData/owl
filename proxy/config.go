package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath          = "./conf/proxy.conf"
	DEFAULT_TCP_BIND        = "0.0.0.0:10030"
	DEFAULT_CFC_ADDR        = ""
	DEFAULT_MAX_PACKET_SIZE = 4096
	DEFAULT_LOG_FILE        = "./logs/proxy.log"
	DEFAULT_LOG_EXPIRE_DAYS = 7
	DEFAULT_LOG_LEVEL       = 3
)

var GlobalConfig *Config

type Config struct {
	CFC_ADDR string

	TCP_BIND string //tcp监听地址和端口

	//LOG CONFIG
	LOG_FILE        string //日志保存目录
	LOG_LEVEL       int    //日志级别
	LOG_EXPIRE_DAYS int    //日志保留天数

	MAX_PACKET_SIZE int
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		TCP_BIND: cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DEFAULT_TCP_BIND),
		CFC_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "cfc_addr", DEFAULT_CFC_ADDR),

		MAX_PACKET_SIZE: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		LOG_FILE:        cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_FILE),
		LOG_EXPIRE_DAYS: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:       cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
	}
	return nil

}

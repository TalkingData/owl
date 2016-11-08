package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath          = "./conf/agent.conf"
	DEFAULT_LOG_PATH        = "./logs/agent.log"
	DEFAULT_CFC_ADDR        = "127.0.0.1:10020"
	DEFAULT_REPEATER_ADDR   = "127.0.0.1:10040"
	DEFAULT_TCP_BIND        = "127.0.0.1:10010"
	DEFAULT_MAX_PACKET_SIZE = 4096
	DEFAULT_LOG_EXPIRE_DAYS = 7
	DEFAULT_LOG_LEVEL       = 3
	DEFAULT_BUFFER_SIZE     = 1024 * 1024
	DEFAULT_CADVISOR_ADDR   = "http://127.0.0.1:8080/"
	DEFAULT_CADVISOR_ENABLE = false
)

var GlobalConfig *Config

type Config struct {
	CFC_ADDR      string //tcp监听地址和端口
	REPEATER_ADDR string
	TCP_BIND      string
	BUFFER_SIZE   int

	CADVISOR_ENABLE bool
	CADVISOR_ADDR   string

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
		CFC_ADDR:      cfg.MustValue(goconfig.DEFAULT_SECTION, "cfc_addr", DEFAULT_CFC_ADDR),
		REPEATER_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "repeater_addr", DEFAULT_REPEATER_ADDR),
		TCP_BIND:      cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DEFAULT_TCP_BIND),
		BUFFER_SIZE:   cfg.MustInt(goconfig.DEFAULT_SECTION, "buffer_size", DEFAULT_BUFFER_SIZE),

		CADVISOR_ENABLE: cfg.MustBool(goconfig.DEFAULT_SECTION, "cadvisor_enable", DEFAULT_CADVISOR_ENABLE),
		CADVISOR_ADDR:   cfg.MustValue(goconfig.DEFAULT_SECTION, "cadvisor_addr", DEFAULT_CADVISOR_ADDR),
		MAX_PACKET_SIZE: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		LOG_FILE:        cfg.MustValue(goconfig.DEFAULT_SECTION, "log_FILE", DEFAULT_LOG_PATH),
		LOG_EXPIRE_DAYS: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:       cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
	}
	return nil

}

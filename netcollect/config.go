package main

import (
	"fmt"
	"strings"

	"github.com/Unknwon/goconfig"
)

const (
	ConfigFilePath           = "./conf/netcollect.conf"
	DEFAULT_REPEATER_ADDR    = "127.0.0.1:10040"
	DEFAULT_CFC_ADDR         = "127.0.0.1:10020"
	DEFAULT_MAX_PACKET_SIZE  = 4096
	DEFAULT_SNMP_TIMEOUT     = 10
	DEFAULT_BUFFER_SIZE      = 1 << 20
	DEFAULT_LOG_FILE         = "./logs/netcollect.log"
	DEFAULT_LOG_EXPIRE_DAYS  = 7
	DEFAULT_LOG_LEVEL        = 3
	DEFAULT_COLLECT_INTERVAL = 60

	DEFAULT_SNMP_PORT      = 161
	DEFAULT_SNMP_VERSION   = "2c"
	DEFAULT_SNMP_COMMUNITY = "public"
)

var GlobalConfig *Config

type Config struct {
	REPEATER_ADDR string
	CFC_ADDR      string
	//LOG CONFIG
	LOG_FILE        string //日志保存目录
	LOG_LEVEL       int    //日志级别
	LOG_EXPIRE_DAYS int    //日志保留天数
	IP_RANGE        []string

	COLLECT_INTERVAL int
	LEGAL_PREFIX     []string

	MAX_PACKET_SIZE int
	BUFFER_SIZE     int64

	SNMP_PORT      int
	SNMP_VERSION   string
	SNMP_COMMUNITY string
	SNMP_TIMEOUT   int
	Metadata       map[string]string
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		IP_RANGE:      cfg.MustValueArray(goconfig.DEFAULT_SECTION, "ip_range", ","),
		REPEATER_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "repeater_addr", DEFAULT_REPEATER_ADDR),
		CFC_ADDR:      cfg.MustValue(goconfig.DEFAULT_SECTION, "cfc_addr", DEFAULT_CFC_ADDR),

		SNMP_PORT:      cfg.MustInt(goconfig.DEFAULT_SECTION, "snmp_port", DEFAULT_SNMP_PORT),
		SNMP_VERSION:   cfg.MustValue(goconfig.DEFAULT_SECTION, "snmp_version", DEFAULT_SNMP_VERSION),
		SNMP_COMMUNITY: cfg.MustValue(goconfig.DEFAULT_SECTION, "snmp_community", DEFAULT_SNMP_COMMUNITY),
		SNMP_TIMEOUT:   cfg.MustInt(goconfig.DEFAULT_SECTION, "snmp_timeout_seconds", DEFAULT_SNMP_TIMEOUT),

		COLLECT_INTERVAL: cfg.MustInt(goconfig.DEFAULT_SECTION, "collect_interval", DEFAULT_COLLECT_INTERVAL),
		LEGAL_PREFIX:     cfg.MustValueArray(goconfig.DEFAULT_SECTION, "legal_prefix", ","),
		MAX_PACKET_SIZE:  cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		BUFFER_SIZE:      cfg.MustInt64(goconfig.DEFAULT_SECTION, "buffer_size", DEFAULT_BUFFER_SIZE),
		LOG_FILE:         cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_FILE),
		LOG_EXPIRE_DAYS:  cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:        cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
	}
	metadata := cfg.MustValue(goconfig.DEFAULT_SECTION, "meta_data", "")
	if len(metadata) > 0 {
		GlobalConfig.Metadata = parseMedata(metadata)
	}
	fmt.Println("config %v", GlobalConfig)
	return nil
}

func parseMedata(val string) map[string]string {
	m := map[string]string{}
	arr1 := strings.Split(val, ",")
	for _, val := range arr1 {
		arr2 := strings.Split(val, "=")
		if len(arr2) != 2 {
			continue
		}
		m[arr2[0]] = arr2[1]
	}
	return m
}

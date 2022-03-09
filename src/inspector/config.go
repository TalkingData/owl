package main

import "github.com/Unknwon/goconfig"

const (
	CONFIG_FILE_PATH          = "./conf/inspector.conf"
	DEFAULT_LOG_PATH          = "./logs/inspector.log"
	DEFAULT_CONTROLER_ADDR    = "127.0.0.1:10050"
	DEFAULT_MAX_PACKET_SIZE   = 409600
	DEFAULT_LOG_EXPIRE_DAYS   = 7
	DEFAULT_LOG_LEVEL         = 3
	DEFAULT_MAX_TASK_BUFFER   = 4096
	DEFAULT_MAX_RESULT_BUFFER = 4096
	DEFAULT_WORKER_COUNT      = 20 //the worker for process tasks
	DEFAULT_TSDB_ADDR         = "127.0.0.1:4242"
	DEFAULT_TSDB_TIMEOUT      = 30
	DEFAULT_RESULT_BUFFER     = 10
	DEFAULT_BACKEND_TSDB      = "opentsdb"
)

var GlobalConfig *Config

type Config struct {
	CONTROLLER_ADDR   string
	LOG_FILE          string
	LOG_LEVEL         int
	LOG_EXPIRE_DAYS   int
	MAX_PACKET_SIZE   int
	MAX_TASK_BUFFER   int
	MAX_RESULT_BUFFER int
	WORKER_COUNT      int
	TSDB_ADDR         string
	TSDB_TIMEOUT      int
	RESULT_BUFFER     int
	BACKEND_TSDB      string
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATH)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		CONTROLLER_ADDR:   cfg.MustValue(goconfig.DEFAULT_SECTION, "controller_addr", DEFAULT_CONTROLER_ADDR),
		MAX_PACKET_SIZE:   cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		LOG_FILE:          cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_PATH),
		LOG_EXPIRE_DAYS:   cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:         cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
		MAX_TASK_BUFFER:   cfg.MustInt(goconfig.DEFAULT_SECTION, "max_task_buffer", DEFAULT_MAX_TASK_BUFFER),
		MAX_RESULT_BUFFER: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_result_buffer", DEFAULT_MAX_RESULT_BUFFER),
		WORKER_COUNT:      cfg.MustInt(goconfig.DEFAULT_SECTION, "worker_count", DEFAULT_WORKER_COUNT),
		TSDB_ADDR:         cfg.MustValue(goconfig.DEFAULT_SECTION, "tsdb_addr", DEFAULT_TSDB_ADDR),
		TSDB_TIMEOUT:      cfg.MustInt(goconfig.DEFAULT_SECTION, "tsdb_timeout", DEFAULT_TSDB_TIMEOUT),
		RESULT_BUFFER:     cfg.MustInt(goconfig.DEFAULT_SECTION, "result_buffer", DEFAULT_RESULT_BUFFER),
		BACKEND_TSDB:      cfg.MustValue(goconfig.DEFAULT_SECTION, "backend_tsdb", DEFAULT_BACKEND_TSDB),
	}
	return nil
}

package main

import (
	"fmt"

	"github.com/Unknwon/goconfig"
)

const (
	ConfigFilePath                  = "./conf/repeater.conf"
	DefaultBackend                  = "opentsdb"
	DefaultTCPBind                  = "0.0.0.0:10040"
	DefaultMetricBind               = "0.0.0.0:10041"
	DefaultOpentsdbAddr             = "127.0.0.1:4242"
	DefaultKairosdbAddr             = "127.0.0.1:4242"
	DefaultRepeaterAddr             = ""
	DefaultMaxPacketSize            = 4096
	DefaultBufferSize               = 1 << 20
	DefaultLogFile                  = "./logs/repeater.log"
	DefaultLogExpireDays            = 7
	DefaultLogLevel                 = 3
	DefaultBackendWriteFailedPolicy = "retry"
)

var (
	RetryPolicy   = "retry"
	SkippedPolicy = "skipped"
)

var GlobalConfig *Config

type Config struct {
	Backend                  []string
	BackendWriteFailedPolicy string
	TCPBind                  string
	MetricBind               string

	OpentsdbAddr string
	KafkaBrokers []string
	KafkaTopic   string

	KairosdbAddr string

	RepeaterAddr string
	//LOG CONFIG
	LogFile       string //日志保存目录
	LogLevel      int    //日志级别
	LogExpireDays int    //日志保留天数

	MaxPacketSize int
	BufferSize    int64
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		Backend:                  cfg.MustValueArray(goconfig.DEFAULT_SECTION, "backend", ","),
		BackendWriteFailedPolicy: cfg.MustValue(goconfig.DEFAULT_SECTION, "backend_write_failed_policy", DefaultBackendWriteFailedPolicy),
		TCPBind:                  cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DefaultTCPBind),
		MetricBind:               cfg.MustValue(goconfig.DEFAULT_SECTION, "metric_bind", DefaultMetricBind),
		OpentsdbAddr:             cfg.MustValue(goconfig.DEFAULT_SECTION, "opentsdb_addr", DefaultOpentsdbAddr),
		KairosdbAddr:             cfg.MustValue(goconfig.DEFAULT_SECTION, "kairosdb_addr", DefaultKairosdbAddr),
		RepeaterAddr:             cfg.MustValue(goconfig.DEFAULT_SECTION, "repeater_addr", DefaultRepeaterAddr),
		KafkaBrokers:             cfg.MustValueArray(goconfig.DEFAULT_SECTION, "kafka_brokers", ","),
		KafkaTopic:               cfg.MustValue(goconfig.DEFAULT_SECTION, "kafka_topic", "owl"),

		MaxPacketSize: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DefaultMaxPacketSize),
		BufferSize:    cfg.MustInt64(goconfig.DEFAULT_SECTION, "buffer_size", DefaultBufferSize),
		LogFile:       cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DefaultLogFile),
		LogExpireDays: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DefaultLogExpireDays),
		LogLevel:      cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DefaultLogLevel),
	}
	if len(GlobalConfig.Backend) == 0 {
		GlobalConfig.Backend = []string{DefaultBackend}
	}
	switch GlobalConfig.BackendWriteFailedPolicy {
	case RetryPolicy, SkippedPolicy:
	default:
		return fmt.Errorf("unsupport backend_write_failed_policy %s", GlobalConfig.BackendWriteFailedPolicy)
	}
	return nil

}

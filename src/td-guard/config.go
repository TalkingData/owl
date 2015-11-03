package main

import "github.com/Unknwon/goconfig"

type Config struct {
	TCPBIND   string
	GUARDBIND string

	LOG_DIR         string //日志保存目录
	LOG_LEVEL       int    //日志级别
	LOG_EXPIRE_DAYS int    //日志保留天数
}

func load_config(config_file string) (*Config, error) {
	cfg, err := goconfig.LoadConfigFile(config_file)
	if err != nil {
		return nil, err
	}
	return &Config{
		TCPBIND:         cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", "127.0.0.1:60010"),
		GUARDBIND:       cfg.MustValue(goconfig.DEFAULT_SECTION, "guard_bind", "127.0.0.1:60011"),
		LOG_DIR:         cfg.MustValue(goconfig.DEFAULT_SECTION, "log_dir", "./logs"),
		LOG_LEVEL:       cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", 4),
		LOG_EXPIRE_DAYS: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", 7),
	}, nil

}

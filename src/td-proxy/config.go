package main

import "github.com/Unknwon/goconfig"

type Config struct {
	TCPBIND  string
	HTTPBIND string

	//SERVER CONFIG
	TCPSERVER          string //tcp监听地址和端口
	HTTPSERVER         string
	BUFFER_SIZE        int //数据缓存buffer大小
	RECONNECTINTERVAL  int //服务器重连间隔,单位分钟
	HOSTUPDATEINTERVAL int //主机配置更新间隔，单位分钟
	MAX_PACKET_SIZE    int

	//LOG CONFIG
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
		TCPBIND:         cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", "0.0.0.0:8888"),
		HTTPBIND:        cfg.MustValue(goconfig.DEFAULT_SECTION, "http_bind", "0.0.0.0:8887"),
		TCPSERVER:       cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_server", ""),
		HTTPSERVER:      cfg.MustValue(goconfig.DEFAULT_SECTION, "http_server", ""),
		BUFFER_SIZE:     cfg.MustInt(goconfig.DEFAULT_SECTION, "buffer_size", 10000),
		MAX_PACKET_SIZE: cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", 4096),

		RECONNECTINTERVAL:  cfg.MustInt(goconfig.DEFAULT_SECTION, "reconnect_interval", 5),
		HOSTUPDATEINTERVAL: cfg.MustInt(goconfig.DEFAULT_SECTION, "host_update_interval", 5),
		LOG_DIR:            cfg.MustValue(goconfig.DEFAULT_SECTION, "log_dir", "./logs"),
		LOG_LEVEL:          cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", 4),
		LOG_EXPIRE_DAYS:    cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", 7),
	}, nil

}

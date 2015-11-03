package main

import "github.com/Unknwon/goconfig"

type Config struct {
	TCPBIND            string
	GUARDBIND          string
	TCPSERVER          string //tcp监听地址和端口
	HTTPSERVER         string
	BUFFER_SIZE        int //数据缓存buffer大小
	RECONNECTINTERVAL  int //服务器重连间隔,单位分钟
	HOSTUPDATEINTERVAL int //主机配置更新间隔，单位分钟

	ASSECTPLUGIN   string
	ASSESTINTERVAL int

	PORTPLUGIN   string
	PORTINTERVAL int

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
		TCPBIND:            cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", "127.0.0.1:65530"),
		GUARDBIND:          cfg.MustValue(goconfig.DEFAULT_SECTION, "guard_bind", "127.0.0.1:65531"),
		TCPSERVER:          cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_server", "127.0.0.1:8888"),
		HTTPSERVER:         cfg.MustValue(goconfig.DEFAULT_SECTION, "http_server", "127.0.0.1:8887"),
		BUFFER_SIZE:        cfg.MustInt(goconfig.DEFAULT_SECTION, "buffer_size", 10000),
		RECONNECTINTERVAL:  cfg.MustInt(goconfig.DEFAULT_SECTION, "reconnect_interval", 5),
		HOSTUPDATEINTERVAL: cfg.MustInt(goconfig.DEFAULT_SECTION, "host_update_interval", 5),

		LOG_DIR:         cfg.MustValue(goconfig.DEFAULT_SECTION, "log_dir", "./logs"),
		LOG_EXPIRE_DAYS: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", 7),
		LOG_LEVEL:       cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", 4),
		ASSECTPLUGIN:    cfg.MustValue(goconfig.DEFAULT_SECTION, "assest_collect_plugin", "get_assest_info.py"),
		ASSESTINTERVAL:  cfg.MustInt(goconfig.DEFAULT_SECTION, "assest_collect_interval", 15),
		PORTPLUGIN:      cfg.MustValue(goconfig.DEFAULT_SECTION, "port_collect_plugin", "get_port_listen.sh"),
		PORTINTERVAL:    cfg.MustInt(goconfig.DEFAULT_SECTION, "port_collect_interval", 10),
	}, nil

}

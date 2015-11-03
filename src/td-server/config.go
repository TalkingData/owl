package main

import "github.com/Unknwon/goconfig"

type Config struct {
	//MYSQL CONFIG
	MYSQL_BIND          string //mysql ip地址
	MYSQL_USER          string //mysql 登陆用户名
	MYSQL_PASSWORD      string //mysql 登陆密码
	MYSQL_DBNAME        string //mysql 数据库名称
	MYSQL_MAX_IDLE_CONN int    //mysql 最大空闲连接数
	MYSQL_MAX_CONN      int    //mysql 最大连接数

	//SERVER CONFIG
	TCPBIND     string //tcp监听地址和端口
	HTTPBIND    string //http监听地址和端口
	BUFFER_SIZE int    //数据缓存buffer大小

	//LOG CONFIG
	LOG_DIR         string //日志保存目录
	LOG_LEVEL       int    //日志级别
	LOG_EXPIRE_DAYS int    //日志保留天数
	//LOG_LEVEL string

	MAX_PACKET_SIZE int

	//TSDB CONFIG
	OPENTSDB_ADDR string //opentsdb 地址
	ENABLE_REDIS  bool
	REDIS_ADDR    string
	REDIS_KEY     string
}

func load_config(config_file string) (*Config, error) {
	cfg, err := goconfig.LoadConfigFile(config_file)
	if err != nil {
		return nil, err
	}
	return &Config{
		TCPBIND:     cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", "0.0.0.0:8888"),
		HTTPBIND:    cfg.MustValue(goconfig.DEFAULT_SECTION, "http_bind", "0.0.0.0:8887"),
		BUFFER_SIZE: cfg.MustInt(goconfig.DEFAULT_SECTION, "buffer_size", 10000),

		MYSQL_BIND:          cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_bind", "127.0.0.1:3306"),
		MYSQL_USER:          cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", "root"),
		MYSQL_DBNAME:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", "td_monitor"),
		MYSQL_PASSWORD:      cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", ""),
		MYSQL_MAX_CONN:      cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", 100),
		MYSQL_MAX_IDLE_CONN: cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", 30),
		MAX_PACKET_SIZE:     cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", 4096),
		LOG_DIR:             cfg.MustValue(goconfig.DEFAULT_SECTION, "log_dir", "./logs"),
		LOG_EXPIRE_DAYS:     cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", 7),
		LOG_LEVEL:           cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", 2),

		OPENTSDB_ADDR: cfg.MustValue(goconfig.DEFAULT_SECTION, "opentsdb_addr", "127.0.0.1:4248"),
		ENABLE_REDIS:  cfg.MustBool(goconfig.DEFAULT_SECTION, "enable_redis", false),
		REDIS_ADDR:    cfg.MustValue(goconfig.DEFAULT_SECTION, "redis_addr", "127.0.0.1:6379"),
		REDIS_KEY:     cfg.MustValue(goconfig.DEFAULT_SECTION, "redis_key", "tdmonitor"),
	}, nil

}

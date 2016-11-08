package main

import "github.com/Unknwon/goconfig"

const (
	ConfigFilePath              = "./conf/cfc.conf"
	DEFAULT_TCP_BIND            = "0.0.0.0:10020"
	DEFAULT_MYSQL_ADDR          = "127.0.0.1:3306"
	DEFAULT_MYSQL_DBNAME        = "owl"
	DEFAULT_MYSQL_USER          = "owl"
	DEFAULT_MYSQL_PASSWORD      = ""
	DEFAULT_MYSQL_MAX_CONN      = 20
	DEFAULT_MYSQL_MAX_IDLE_CONN = 5
	DEFAULT_MAX_PACKET_SIZE     = 4096
	DEFAULT_LOG_FILE            = "./logs/cfc.log"
	DEFAULT_LOG_EXPIRE_DAYS     = 7
	DEFAULT_LOG_LEVEL           = 3
)

var GlobalConfig *Config

type Config struct {
	//MYSQL CONFIG
	MYSQL_ADDR          string //mysql ip地址
	MYSQL_USER          string //mysql 登陆用户名
	MYSQL_PASSWORD      string //mysql 登陆密码
	MYSQL_DBNAME        string //mysql 数据库名称
	MYSQL_MAX_IDLE_CONN int    //mysql 最大空闲连接数
	MYSQL_MAX_CONN      int    //mysql 最大连接数

	//SERVER CONFIG
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
		TCP_BIND:            cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DEFAULT_TCP_BIND),
		MYSQL_ADDR:          cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_addr", DEFAULT_MYSQL_ADDR),
		MYSQL_USER:          cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", DEFAULT_MYSQL_USER),
		MYSQL_DBNAME:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", DEFAULT_MYSQL_DBNAME),
		MYSQL_PASSWORD:      cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", DEFAULT_MYSQL_PASSWORD),
		MYSQL_MAX_CONN:      cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", DEFAULT_MYSQL_MAX_CONN),
		MYSQL_MAX_IDLE_CONN: cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", DEFAULT_MYSQL_MAX_IDLE_CONN),
		MAX_PACKET_SIZE:     cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		LOG_FILE:            cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_FILE),
		LOG_EXPIRE_DAYS:     cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:           cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
	}
	return nil

}

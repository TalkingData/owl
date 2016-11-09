package main

import (
	"github.com/Unknwon/goconfig"
)

const (
	CONFIG_FILE_PATH                    = "./conf/api.conf"
	DEFAULT_HTTP_BIND                   = ":10060"
	DEFAULT_MYSQL_ADDR                  = "127.0.0.1:3306"
	DEFAULT_MYSQL_USER                  = "root"
	DEFAULT_MYSQL_DBNAME                = "owl"
	DEFAULT_MYSQL_PASSWORD              = ""
	DEFAULT_MAX_CONN                    = 20
	DEFAULT_MAX_IDLE_CONN               = 5
	DEFAULT_SECRET_KEY                  = "b690dfddeb13156b6b88946708210a90a3df1d285576e843c8870a2090226329"
	DEFAULT_OPENTSDB_ADDR               = "127.0.0.1:4242"
	DEFAULT_OPENTSDB_TIMEOUT            = 5
	DEFAULT_AUTO_BUILD_METRIC_TAG_INDEX = false
	DEFAULT_AUTO_BUILD_INTERVAL         = 10 //minute
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
	HTTP_BIND string //http 监听地址和端口
	//AUTH CONFIG
	SECRET_KEY string //秘钥

	OPENTSDB_ADDR    string
	OPENTSDB_TIMEOUT int

	AUTO_BUILD_METRIC_TAG_INDEX bool
	AUTO_BUILD_INTERVAL         int
}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATH)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		HTTP_BIND:                   cfg.MustValue(goconfig.DEFAULT_SECTION, "http_bind", DEFAULT_HTTP_BIND),
		MYSQL_ADDR:                  cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_addr", DEFAULT_MYSQL_ADDR),
		MYSQL_USER:                  cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", DEFAULT_MYSQL_USER),
		MYSQL_DBNAME:                cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", DEFAULT_MYSQL_DBNAME),
		MYSQL_PASSWORD:              cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", DEFAULT_MYSQL_PASSWORD),
		MYSQL_MAX_CONN:              cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", DEFAULT_MAX_CONN),
		MYSQL_MAX_IDLE_CONN:         cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", DEFAULT_MAX_IDLE_CONN),
		SECRET_KEY:                  cfg.MustValue(goconfig.DEFAULT_SECTION, "secret_key", DEFAULT_SECRET_KEY),
		OPENTSDB_ADDR:               cfg.MustValue(goconfig.DEFAULT_SECTION, "opentsdb_addr", DEFAULT_OPENTSDB_ADDR),
		OPENTSDB_TIMEOUT:            cfg.MustInt(goconfig.DEFAULT_SECTION, "opentsdb_timeout", DEFAULT_OPENTSDB_TIMEOUT),
		AUTO_BUILD_METRIC_TAG_INDEX: cfg.MustBool(goconfig.DEFAULT_SECTION, "auto_build_metric_tag_index", DEFAULT_AUTO_BUILD_METRIC_TAG_INDEX),
		AUTO_BUILD_INTERVAL:         cfg.MustInt(goconfig.DEFAULT_SECTION, "auto_build_interval", DEFAULT_AUTO_BUILD_INTERVAL),
	}
	return nil
}

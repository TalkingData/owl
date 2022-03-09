package main

import (
	"fmt"

	"github.com/Unknwon/goconfig"
)

const (

	//配置文件路径
	configFilePath        = "./conf/api.conf"
	defaultHTTPBind       = "0.0.0.0:10060"
	defaultMetricBind     = "0.0.0.0:10061"
	defaultAuthType       = "mysql"
	defaultPublicKeyFile  = "./certs/owl-api.key.pub"
	defaultPrivateKeyFile = "./certs/owl-api.key"
	defaultTokenExpired   = 7

	defaultMySQLAddr        = "127.0.0.1:3306"
	defaultMySQLDBName      = "owl"
	defaultMySQLUser        = "root"
	defaultMySQLPassword    = ""
	defaultMySQLMaxConn     = 20
	defaultMySQLMaxIdleConn = 5

	defaultLogFile       = "./logs/api.log"
	defaultLogExpireDays = 7
	defaultLogLevel      = 7

	defaultTimeSeriesStorage   = "opentsdb"
	defaultKairosdbAddr        = "127.0.0.1:8080"
	defaultOpentsdbAddr        = "127.0.0.1:4242"
	defaultOpentsdbReadTimeout = 10

	defaultAlarmHealthCheckURL = "http://127.0.0.1:10051"
)

var config *Config

// Config 配置
type Config struct {
	HTTPBind   string //api 服务局监听的本地地址和端口
	MetricBind string

	MySQLAddr        string //mysql ip地址
	MySQLUser        string //mysql 登陆用户名
	MySQLPassword    string //mysql 登陆密码
	MySQLDBName      string //mysql 数据库名称
	MySQLMaxIdleConn int    //mysql 最大空闲连接数
	MySQLMaxConn     int    //mysql 最大连接数

	TimeSeriesStorage string //目前仅支持 opentsdb, kairosdb

	OpentsdbAddr         string //opentsdb rest api 地址
	OpenttsdbReadTimeout int    //opentsdb 请求超时时间,单位秒

	KairosdbAddr string //kairosdb rest api 接口地址

	LogFile       string //日志保存目录
	LogLevel      int    //日志级别
	LogExpireDays int    //日志保留天数

	AuthType string //目前支持mysql, iam
	IamURL   string //iam 地址
	AppID    string //iam 中注册的 app id
	AppKey   string //iam 中注册的 app key

	PublicKeyFile  string //公钥文件路径
	PrivateKeyFile string //私钥文件路径
	TokenExpired   int    //token 超时时间,仅在 AuthType 为 mysql 时有效

	AlarmHealthCheckURL string //报警服务状态检查路径
}

func (c *Config) validate() error {
	switch c.AuthType {
	case "mysql":
	case "iam":
		if c.IamURL == "" {
			return fmt.Errorf("iam url must specified")
		}
		if c.AppID == "" {
			return fmt.Errorf("app id must specified")
		}
		if c.AppKey == "" {
			return fmt.Errorf("app key must specified")
		}
	default:
		return fmt.Errorf("auth type %s not support", c.AuthType)
	}

	switch c.TimeSeriesStorage {
	case "opentsdb", "kairosdb":
	default:
		return fmt.Errorf("timeseries storage %s not support", c.TimeSeriesStorage)
	}
	return nil
}

// InitGlobalConfig 初始化配置
func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(configFilePath)
	if err != nil {
		return err
	}
	config = &Config{
		HTTPBind:       cfg.MustValue(goconfig.DEFAULT_SECTION, "http_bind", defaultHTTPBind),
		MetricBind:     cfg.MustValue(goconfig.DEFAULT_SECTION, "metric_bind", defaultMetricBind),
		PublicKeyFile:  cfg.MustValue(goconfig.DEFAULT_SECTION, "public_key", defaultPublicKeyFile),
		PrivateKeyFile: cfg.MustValue(goconfig.DEFAULT_SECTION, "private_key", defaultPrivateKeyFile),

		AuthType:     cfg.MustValue(goconfig.DEFAULT_SECTION, "auth_type", defaultAuthType),
		IamURL:       cfg.MustValue(goconfig.DEFAULT_SECTION, "iam_url", ""),
		AppID:        cfg.MustValue(goconfig.DEFAULT_SECTION, "app_id", ""),
		AppKey:       cfg.MustValue(goconfig.DEFAULT_SECTION, "api_key", ""),
		TokenExpired: cfg.MustInt(goconfig.DEFAULT_SECTION, "token_expired", defaultTokenExpired),

		TimeSeriesStorage:    cfg.MustValue(goconfig.DEFAULT_SECTION, "timeseirs_storage", defaultTimeSeriesStorage),
		KairosdbAddr:         cfg.MustValue(goconfig.DEFAULT_SECTION, "kairosdb_addr", defaultKairosdbAddr),
		OpentsdbAddr:         cfg.MustValue(goconfig.DEFAULT_SECTION, "opentsdb_addr", defaultOpentsdbAddr),
		OpenttsdbReadTimeout: cfg.MustInt(goconfig.DEFAULT_SECTION, "opentsdb_timeout", defaultOpentsdbReadTimeout),

		MySQLAddr:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_addr", defaultMySQLAddr),
		MySQLUser:        cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", defaultMySQLUser),
		MySQLDBName:      cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", defaultMySQLDBName),
		MySQLPassword:    cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", defaultMySQLPassword),
		MySQLMaxConn:     cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", defaultMySQLMaxConn),
		MySQLMaxIdleConn: cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", defaultMySQLMaxIdleConn),

		LogFile:       cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", defaultLogFile),
		LogExpireDays: cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", defaultLogExpireDays),
		LogLevel:      cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", defaultLogLevel),

		AlarmHealthCheckURL: cfg.MustValue(goconfig.DEFAULT_SECTION, "alarm_health_check_url", defaultAlarmHealthCheckURL),
	}
	return config.validate()
}

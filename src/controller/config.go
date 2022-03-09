package main

import (
	"github.com/Unknwon/goconfig"
)

const (
	CONFIG_FILE_PATH                 = "./conf/controller.conf"
	DEFAULT_TCP_BIND                 = ":10050"
	DEFAULT_MYSQL_ADDR               = "127.0.0.1:3306"
	DEFAULT_MYSQL_USER               = "root"
	DEFAULT_MYSQL_DBNAME             = "owl"
	DEFAULT_MYSQL_PASSWORD           = ""
	DEFAULT_MAX_CONN                 = 20
	DEFAULT_MAX_IDLE_CONN            = 5
	DEFAULT_LOG_FILE                 = "./logs/controller.log"
	DEFAULT_LOG_EXPIRE_DAYS          = 7
	DEFAULT_LOG_LEVEL                = 3
	DEFAULT_MAX_PACKET_SIZE          = 409600
	DEFAULT_LOAD_STRATEGIES_INTERVAL = 300 //seconds
	DEFAULT_TASK_POOL_SIZE           = 409600
	DEFAULT_RESULT_POOL_SIZE         = 409600
	DEFAULT_EVENT_POOL_SIZE          = 1024
	DEFAULT_HTTP_SERVER              = ":10051"
	DEFAULT_TASK_SIZE                = 100
	DEFAULT_WORKER_COUNT             = 20
	DEFAULT_ACTION_TIMEOUT           = 60 //seconds
	DEFAULT_SEND_MAX                 = 100
	DEFAULT_MAX_INTERVAL_WAIT_TIME   = 5 //seconds
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

	MAX_PACKET_SIZE int //tcp 消息报最大size

	LOAD_STRATEGIES_INTERVAL int //获取策略时间间隔 单位秒

	TASK_POOL_SIZE int //任务池的缓冲大小

	RESULT_POOL_SIZE int //结果池的缓冲大小

	EVENT_POOL_SIZE int

	HTTP_SERVER string //Http服务的地址

	TASK_SIZE int //单次获取任务数

	WORKER_COUNT int //处理结果池的协程数

	ACTION_TIMEOUT int //报警动作超时时间

	SEND_MAX int //队列报警大于此值减速

	MAX_INTERVAL_WAIT_TIME int //队列报警发送间隔最大上限值

}

func InitGlobalConfig() error {
	cfg, err := goconfig.LoadConfigFile(CONFIG_FILE_PATH)
	if err != nil {
		return err
	}
	GlobalConfig = &Config{
		TCP_BIND:                 cfg.MustValue(goconfig.DEFAULT_SECTION, "tcp_bind", DEFAULT_TCP_BIND),
		MYSQL_ADDR:               cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_addr", DEFAULT_MYSQL_ADDR),
		MYSQL_USER:               cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_user", DEFAULT_MYSQL_USER),
		MYSQL_DBNAME:             cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_dbname", DEFAULT_MYSQL_DBNAME),
		MYSQL_PASSWORD:           cfg.MustValue(goconfig.DEFAULT_SECTION, "mysql_password", DEFAULT_MYSQL_PASSWORD),
		MYSQL_MAX_CONN:           cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_conn", DEFAULT_MAX_CONN),
		MYSQL_MAX_IDLE_CONN:      cfg.MustInt(goconfig.DEFAULT_SECTION, "mysql_max_idle_conn", DEFAULT_MAX_IDLE_CONN),
		LOG_FILE:                 cfg.MustValue(goconfig.DEFAULT_SECTION, "log_file", DEFAULT_LOG_FILE),
		LOG_EXPIRE_DAYS:          cfg.MustInt(goconfig.DEFAULT_SECTION, "log_expire_days", DEFAULT_LOG_EXPIRE_DAYS),
		LOG_LEVEL:                cfg.MustInt(goconfig.DEFAULT_SECTION, "log_level", DEFAULT_LOG_LEVEL),
		MAX_PACKET_SIZE:          cfg.MustInt(goconfig.DEFAULT_SECTION, "max_packet_size", DEFAULT_MAX_PACKET_SIZE),
		LOAD_STRATEGIES_INTERVAL: cfg.MustInt(goconfig.DEFAULT_SECTION, "load_strategies_interval", DEFAULT_LOAD_STRATEGIES_INTERVAL),
		TASK_POOL_SIZE:           cfg.MustInt(goconfig.DEFAULT_SECTION, "task_pool_size", DEFAULT_TASK_POOL_SIZE),
		RESULT_POOL_SIZE:         cfg.MustInt(goconfig.DEFAULT_SECTION, "result_pool_size", DEFAULT_RESULT_POOL_SIZE),
		EVENT_POOL_SIZE:          cfg.MustInt(goconfig.DEFAULT_SECTION, "event_pool_size", DEFAULT_EVENT_POOL_SIZE),
		HTTP_SERVER:              cfg.MustValue(goconfig.DEFAULT_SECTION, "http_server", DEFAULT_HTTP_SERVER),
		TASK_SIZE:                cfg.MustInt(goconfig.DEFAULT_SECTION, "task_size", DEFAULT_TASK_SIZE),
		WORKER_COUNT:             cfg.MustInt(goconfig.DEFAULT_SECTION, "worker_count", DEFAULT_WORKER_COUNT),
		ACTION_TIMEOUT:           cfg.MustInt(goconfig.DEFAULT_SECTION, "action_timeout", DEFAULT_ACTION_TIMEOUT),
		SEND_MAX:                 cfg.MustInt(goconfig.DEFAULT_SECTION, "send_max", DEFAULT_SEND_MAX),
		MAX_INTERVAL_WAIT_TIME:   cfg.MustInt(goconfig.DEFAULT_SECTION, "max_interval_wait_time", DEFAULT_MAX_INTERVAL_WAIT_TIME),
	}
	return nil
}

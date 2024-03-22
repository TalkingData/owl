package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Option func(d *gorm.DB)

// MysqlDefaultLogMode 设置Mysql默认LogMode
func MysqlDefaultLogMode(level logger.LogLevel) Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		d.Logger = logger.Default.LogMode(level)
	}
}

// MysqlMaxIdleConns 设置Mysql最大空闲连接数量
func MysqlMaxIdleConns(count int) Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxIdleConns(count)
	}
}

// MysqlMaxOpenConns 设置Mysql最大打开连接数量
func MysqlMaxOpenConns(count int) Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxOpenConns(count)
	}
}

// MysqlConnMaxLifetime 设置Mysql连接可复用的最大时间
func MysqlConnMaxLifetime(secs int) Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(secs))
	}
}

// SkipDefaultTransaction 是否忽略默认事务
func SkipDefaultTransaction(skip bool) Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		d.SkipDefaultTransaction = skip
	}
}

// UsingOpentracingPlugin 使用OpentracingPlugin
func UsingOpentracingPlugin() Option {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		_ = d.Use(new(OpentracingPlugin))
	}
}

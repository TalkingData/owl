package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// NewMysqlGorm 新建MysqlGorm
func NewMysqlGorm(address, user, password, dbName string, opts ...OrmOption) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		address,
		dbName,
	)

	m := mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        200,
		DisableDatetimePrecision: true,
	})

	d, err := gorm.Open(m, &gorm.Config{})
	if err != nil {
		panic(err.Error())
		return nil
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

// MysqlDefaultLogMode 设置Mysql默认LogMode
func MysqlDefaultLogMode(level logger.LogLevel) OrmOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		d.Logger = logger.Default.LogMode(level)
	}
}

// MysqlMaxIdleConns 设置Mysql最大空闲连接数量
func MysqlMaxIdleConns(count int) OrmOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxIdleConns(count)
	}
}

// MysqlMaxOpenConns 设置Mysql最大打开连接数量
func MysqlMaxOpenConns(count int) OrmOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetMaxOpenConns(count)
	}
}

// MysqlConnMaxLifetime 设置Mysql连接可复用的最大时间
func MysqlConnMaxLifetime(secs int) OrmOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		sqlDB, _ := d.DB()
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(secs))
	}
}

// SkipDefaultTransaction 是否忽略默认事务
func SkipDefaultTransaction(skip bool) OrmOption {
	return func(d *gorm.DB) {
		if d == nil {
			return
		}
		d.SkipDefaultTransaction = skip
	}
}

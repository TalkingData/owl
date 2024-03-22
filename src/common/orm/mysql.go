package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysqlGorm 新建MysqlGorm
func NewMysqlGorm(address, user, password, dbName string, opts ...Option) *gorm.DB {
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

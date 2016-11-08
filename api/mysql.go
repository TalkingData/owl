package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var mydb *gorm.DB

func InitMysqlConnPool() error {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
			GlobalConfig.MYSQL_USER,
			GlobalConfig.MYSQL_PASSWORD,
			GlobalConfig.MYSQL_ADDR,
			GlobalConfig.MYSQL_DBNAME))
	if err != nil {
		return err
	}

	db.DB().SetMaxIdleConns(GlobalConfig.MYSQL_MAX_IDLE_CONN)
	db.DB().SetMaxOpenConns(GlobalConfig.MYSQL_MAX_CONN)
	db.LogMode(false)

	mydb = db
	return nil
}

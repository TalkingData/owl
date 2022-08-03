package dao

import (
	"gorm.io/gorm"
	"owl/common/logger"
)

type Dao struct {
	db *gorm.DB
	lg *logger.Logger
}

func NewDao(d *gorm.DB, lg *logger.Logger) *Dao {
	return &Dao{
		d, lg,
	}
}

func (d *Dao) Close() {
	if d == nil {
		return
	}

	db, _ := d.db.DB()
	if db != nil {
		_ = db.Close()
	}
}

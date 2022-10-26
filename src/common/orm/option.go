package orm

import (
	"gorm.io/gorm"
)

type Option func(d *gorm.DB)

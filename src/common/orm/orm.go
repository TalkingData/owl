package orm

import (
	"gorm.io/gorm"
)

type OrmOption func(d *gorm.DB)

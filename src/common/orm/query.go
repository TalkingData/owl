package orm

import "gorm.io/gorm"

type Query map[string]interface{}

func (q Query) Where(d *gorm.DB) *gorm.DB {
	for k, v := range q {
		d = d.Where(k, v)
	}

	return d
}

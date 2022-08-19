package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetScript(query orm.Query) (s *model.Script, err error) {
	res := query.Where(d.db).Limit(1).Find(&s)
	return s, res.Error
}

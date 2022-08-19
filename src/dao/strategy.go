package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListStrategies(query orm.Query) (ss []*model.Strategy, err error) {
	res := query.Where(d.db).Find(&ss)
	return ss, res.Error
}

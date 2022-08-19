package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListActions(query orm.Query) (as []*model.Action, err error) {
	res := query.Where(d.db).Find(&as)
	return as, res.Error
}

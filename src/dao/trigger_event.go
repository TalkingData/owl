package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListTriggerEvents(query orm.Query) (te []*model.TriggerEvent, err error) {
	res := query.Where(d.db).Find(&te)
	return te, res.Error
}

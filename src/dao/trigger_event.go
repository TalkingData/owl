package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListTriggerEvents(ctx context.Context, q orm.Query) (te []*model.TriggerEvent, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&te)
	return te, res.Error
}

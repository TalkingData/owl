package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListActions(ctx context.Context, q orm.Query) (as []*model.Action, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&as)
	return as, res.Error
}

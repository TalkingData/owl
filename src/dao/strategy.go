package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListStrategies(ctx context.Context, q orm.Query) (ss []*model.Strategy, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&ss)
	return ss, res.Error
}

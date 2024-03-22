package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetScript(ctx context.Context, q orm.Query) (s *model.Script, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&s)
	return s, res.Error
}

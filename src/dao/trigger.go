package dao

import (
	"context"
	"owl/model"
)

func (d *Dao) ListTriggersByStrategyId(ctx context.Context, strategyId uint64) (tri []*model.Trigger, err error) {
	res := d.getDbWithCtx(ctx).Where("strategy_id=?", strategyId).
		Find(&tri)
	return tri, res.Error
}

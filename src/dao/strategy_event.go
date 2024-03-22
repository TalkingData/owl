package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetStrategyEvent(
	ctx context.Context,
	strategyId uint64, hostId string, status uint32,
) (se *model.StrategyEvent, err error) {
	q := orm.Query{
		"strategy_id": strategyId,
		"host_id":     hostId,
		"status":      status,
	}
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&se)
	return se, res.Error
}

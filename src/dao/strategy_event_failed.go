package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewStrategyEventFailed(
	ctx context.Context,
	strategyId uint64, status int32, hostId, message string,
) (*model.StrategyEventFailed, error) {
	p := model.StrategyEventFailed{
		StrategyId: strategyId,
		Status:     status,
		HostId:     hostId,
		Message:    message,
	}

	res := d.getDbWithCtx(ctx).Create(&p)
	return &p, res.Error
}

func (d *Dao) RemoveStrategyEventFailed(ctx context.Context, strategyId uint64, hostId string) error {
	q := orm.Query{
		"strategy_id=?": strategyId,
		"host_id=?":     hostId,
	}

	res := q.Where(d.getDbWithCtx(ctx)).Delete(model.StrategyEventFailed{})
	return res.Error
}

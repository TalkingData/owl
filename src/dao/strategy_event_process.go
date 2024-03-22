package dao

import (
	"context"
	"owl/model"
)

func (d *Dao) NewStrategyEventProcess(
	ctx context.Context,
	strategyEventId uint64, strategyEventStatus int32, processUser, processComments string,
) (*model.StrategyEventProcess, error) {
	sep := model.StrategyEventProcess{
		StrategyEventId: strategyEventId,
		ProcessStatus:   strategyEventStatus,
		ProcessUser:     processUser,
		ProcessComments: processComments,
	}

	res := d.getDbWithCtx(ctx).Create(&sep)
	return &sep, res.Error
}

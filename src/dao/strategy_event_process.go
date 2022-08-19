package dao

import "owl/model"

func (d *Dao) NewStrategyEventProcess(
	strategyEventId uint64, strategyEventStatus int, processUser, processComments string,
) (*model.StrategyEventProcess, error) {
	sep := model.StrategyEventProcess{
		StrategyEventId: strategyEventId,
		ProcessStatus:   strategyEventStatus,
		ProcessUser:     processUser,
		ProcessComments: processComments,
	}

	res := d.db.Create(&sep)
	return &sep, res.Error
}

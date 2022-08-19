package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewStrategyEventFailed(strategyId uint64, status int, hostId, message string) (*model.StrategyEventFailed, error) {
	p := model.StrategyEventFailed{
		StrategyId: strategyId,
		Status:     status,
		HostId:     hostId,
		Message:    message,
	}

	res := d.db.Create(&p)
	return &p, res.Error
}

func (d *Dao) RemoveStrategyEventFailed(strategyId uint64, hostId string) error {
	query := orm.Query{
		"strategy_id=?": strategyId,
		"host_id=?":     hostId,
	}

	res := query.Where(d.db).Delete(model.StrategyEventFailed{})
	return res.Error
}

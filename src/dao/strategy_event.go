package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetStrategyEvent(strategyId uint64, hostId string, status uint) (se *model.StrategyEvent, err error) {
	query := orm.Query{
		"strategy_id": strategyId,
		"host_id":     hostId,
		"status":      status,
	}
	res := query.Where(d.db).Limit(1).Find(&se)
	return se, res.Error
}

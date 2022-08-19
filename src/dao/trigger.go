package dao

import "owl/model"

func (d *Dao) ListTriggersByStrategyId(strategyId uint64) (tri []*model.Trigger, err error) {
	res := d.db.Where("strategy_id=?", strategyId).
		Find(&tri)
	return tri, res.Error
}

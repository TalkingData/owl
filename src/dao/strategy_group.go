package dao

import (
	"owl/model"
)

// ListHostGroupsByStrategyId 根据StrategyId列出对应所有HostGroup
func (d *Dao) ListHostGroupsByStrategyId(strategyId uint64) (hostGroups []*model.HostGroup, err error) {
	subQuery := d.db.Model(&model.StrategyGroup{}).
		Select("group_id").
		Where("strategy_id=?", strategyId)

	res := d.db.Where("id IN (?)", subQuery).
		Find(&hostGroups)

	return hostGroups, res.Error
}

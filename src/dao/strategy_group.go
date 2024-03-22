package dao

import (
	"context"
	"owl/model"
)

// ListHostGroupsByStrategyId 根据StrategyId列出对应所有HostGroup
func (d *Dao) ListHostGroupsByStrategyId(
	ctx context.Context,
	strategyId uint64,
) (hostGroups []*model.HostGroup, err error) {
	subQ := d.getDbWithCtx(ctx).Model(&model.StrategyGroup{}).
		Select("group_id").
		Where("strategy_id=?", strategyId)

	res := d.getDbWithCtx(ctx).Where("id IN (?)", subQ).
		Find(&hostGroups)

	return hostGroups, res.Error
}

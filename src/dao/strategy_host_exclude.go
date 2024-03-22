package dao

import (
	"context"
	"owl/model"
)

// ListExcludedHostsIdsByStrategyId 根据StrategyId列出被对应Strategy排除在外的所有HostsId
func (d *Dao) ListExcludedHostsIdsByStrategyId(ctx context.Context, strategyId uint64) (hostIds []string, err error) {
	exHosts := make([]*model.StrategyHostExclude, 0)
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{
		"strategy_id": strategyId,
	}).Find(&exHosts)

	for _, h := range exHosts {
		hostIds = append(hostIds, h.HostId)
	}

	return hostIds, res.Error
}

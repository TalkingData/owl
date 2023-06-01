package biz

import (
	"context"
	"owl/common/logger"
)

func (b *Biz) RefreshHostStatus(ctx context.Context) {
	b.logger.Info("Biz.RefreshHostStatus called.")
	defer b.logger.Info("Biz.RefreshHostStatus end.")

	// 主机的最近更新时间小于HostDownStatusThresholdSecs的，则将其状态变更为OK
	err := b.dao.SetHostStatus2OkByThresholdSecs(
		ctx,
		b.conf.HostDownStatusThresholdSecs,
		b.conf.Const.ExecSqlBatchLimit,
		b.conf.Const.ExecSqlBatchIntervalMs,
	)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "Biz.refreshHostStatus set host to 'OK' status failed.")
		return
	}

	// 主机的最近更新时间大于HostDownStatusThresholdSecs的，则将其状态变更为DOWN
	err = b.dao.SetHostStatus2DownByThresholdSecs(
		ctx,
		b.conf.HostDownStatusThresholdSecs,
		b.conf.Const.ExecSqlBatchLimit,
		b.conf.Const.ExecSqlBatchIntervalMs,
	)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "Biz.refreshHostStatus set host to 'DOWN' status failed.")
	}
}

package biz

import (
	"context"
	"owl/common/logger"
)

func (b *Biz) CleanExpiredMetric(ctx context.Context) {
	b.logger.Info("Biz.CleanExpiredMetric called.")
	defer b.logger.Info("Biz.CleanExpiredMetric end.")

	err := b.dao.CleanExpiredMetric(
		ctx,
		b.conf.CleanExpiredMetricCycleExpiredRatio,
		b.conf.Const.ExecSqlBatchLimit,
		b.conf.Const.ExecSqlBatchIntervalMs,
	)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling dao.CleanExpiredMetric.")
	}
}

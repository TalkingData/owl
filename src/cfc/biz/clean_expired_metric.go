package biz

import (
	"context"
	"owl/common/logger"
	"time"
)

func (b *Biz) CleanExpiredMetric(ctx context.Context) {
	tk := time.Tick(b.conf.CleanExpiredMetricIntervalSecs)
	if tk == nil {
		b.logger.Info("Biz.CleanExpiredMetric not enabled or conf.CleanExpiredMetricIntervalSecs is 0.")
		return
	}
	for {
		select {
		case <-tk:
			if err := b.dao.CleanExpiredMetric(
				b.conf.CleanExpiredMetricCycleExpiredRatio,
				b.conf.Const.CleanExpiredMetricBatchLimit,
				b.conf.Const.CleanExpiredMetricBatchIntervalMs,
			); err != nil {
				b.logger.WarnWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while dao.CleanExpiredMetric.")
			}
		case <-ctx.Done():
			b.logger.InfoWithFields(logger.Fields{
				"context_error": ctx.Err(),
			}, "Biz.CleanExpiredMetric exited by context done.")
			return
		}
	}
}

package biz

import (
	"context"
	"owl/common/logger"
	"owl/common/orm"
	"owl/model"
	"time"
)

func (b *Biz) RefreshHostStatus(ctx context.Context) {
	tk := time.Tick(b.conf.RefreshHostStatusIntervalSecs)
	if tk == nil {
		b.logger.Info("Biz.RefreshHostStatus not enabled or conf.RefreshHostStatusIntervalSecs is 0.")
		return
	}
	for {
		select {
		case <-tk:
			b.refreshHostStatus()
		case <-ctx.Done():
			b.logger.InfoWithFields(logger.Fields{
				"context_error": ctx.Err(),
			}, "Biz.RefreshHostStatus exited by context done.")
			return
		}
	}
}

func (b *Biz) refreshHostStatus() {
	b.refreshHostStatus2Ok()
	b.refreshHostStatus2Down()
}

func (b *Biz) refreshHostStatus2Ok() {
	// 主机的最近更新时间小于HostDownStatusThresholdSecs的，则将其状态变更为OK
	ct, err := b.dao.SetHostStatusByQuery(orm.Query{
		"UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at) <= ?": b.conf.HostDownStatusThresholdSecs,
		"status=?": model.HostStatusDown,
	}, model.HostStatusOk)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"rows_affected": ct,
			"error":         err,
		}, "Biz.refreshHostStatus set host to 'OK' status failed.")
		return
	}
	b.logger.InfoWithFields(logger.Fields{
		"rows_affected": ct,
	}, "Biz.refreshHostStatus set host to 'OK' status done.")
}

func (b *Biz) refreshHostStatus2Down() {
	// 主机的最近更新时间大于HostDownStatusThresholdSecs的，说明主机已经很久没上报心跳了，则将其状态变更为DOWN
	ct, err := b.dao.SetHostStatusByQuery(orm.Query{
		"UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at) > ?": b.conf.HostDownStatusThresholdSecs,
		"status=?": model.HostStatusOk,
	}, model.HostStatusDown)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"rows_affected": ct,
			"error":         err,
		}, "Biz.refreshHostStatus set host to 'DOWN' status failed.")
		return
	}
	b.logger.InfoWithFields(logger.Fields{
		"rows_affected": ct,
	}, "Biz.refreshHostStatus set host to 'DOWN' status done.")
}

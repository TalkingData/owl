package biz

import (
	"context"
	"owl/common/logger"
	"owl/common/orm"
	"owl/common/utils"
)

func (b *Biz) ReceiveAgentMetric(
	ctx context.Context,
	hostId, metric, dataType string,
	cycle int32,
	tags map[string]string,
) {
	// 如果主机id是空，将尝试从ts data获取主机名，根据主机名从数据库查找主机id
	if hostId == "" {
		hostname := tags["host"]
		hostObj, err := b.dao.GetHost(ctx, orm.Query{"hostname": hostname})
		if err != nil {
			b.logger.WarnWithFields(logger.Fields{
				"hostname": hostname,
				"error":    err,
			}, "An error occurred while calling dao.GetHost.")
			return
		}
		if hostObj == nil || hostObj.Id == "" {
			b.logger.WarnWithFields(logger.Fields{
				"hostname": hostname,
				"error":    err,
			}, "Host object not found, Skipped receive agent metric.")
			return
		}

		hostId = hostObj.Id
	}

	// 删除标记主机唯一标识的tags，才能使得用户在界面中查询时可聚合
	delete(tags, "host")
	delete(tags, "uuid")

	b.logger.DebugWithFields(logger.Fields{
		"host_id":   hostId,
		"metric":    metric,
		"data_type": dataType,
		"cycle":     cycle,
		"tags":      tags,
	}, "Biz.ReceiveAgentMetric prepare execute dao.SetOrNetMetric.")
	_, err := b.dao.SetOrNetMetric(ctx, hostId, metric, utils.Tags2String(tags), dataType, cycle)
	if err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"host_id":   hostId,
			"metric":    metric,
			"data_type": dataType,
			"cycle":     cycle,
			"tags":      tags,
			"error":     err,
		}, "An error occurred while calling dao.SetOrNetMetric.")
		return
	}
}

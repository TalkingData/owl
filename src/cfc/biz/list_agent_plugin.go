package biz

import (
	"owl/common/logger"
	"owl/model"
)

func (b *Biz) ListAgentPlugins(hostId string) ([]*model.Plugin, error) {
	finalPlugins := make([]*model.Plugin, 0)
	idMap := make(map[string]struct{})

	// 查找主机的插件
	hPlugins, err := b.dao.ListHostsPluginsByHostId(hostId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id": hostId,
			"error":   err,
		}, "An error occurred while dao.ListHostsPlugins.")
		return nil, err
	}
	for _, p := range hPlugins {
		uniqueKey := p.GenUniqueKey()
		if _, ok := idMap[uniqueKey]; ok {
			b.logger.Warn(logger.Fields{
				"unique_key": uniqueKey,
				"plugin":     p,
			}, "Duplicate found, skipped.")
			continue
		}
		idMap[uniqueKey] = struct{}{}
		finalPlugins = append(finalPlugins, p)
	}

	// 查找主机所在组的插件
	hgPlugins, err := b.dao.ListHostGroupsPluginsByHostId(hostId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"host_id": hostId,
			"error":   err,
		}, "An error occurred while dao.ListHostGroupsPluginsByHostId.")
		return nil, err
	}
	for _, p := range hgPlugins {
		uniqueKey := p.GenUniqueKey()
		if _, ok := idMap[uniqueKey]; ok {
			b.logger.Warn(logger.Fields{
				"unique_key": uniqueKey,
				"plugin":     p,
			}, "Duplicate found, skipped.")
			continue
		}
		idMap[uniqueKey] = struct{}{}
		finalPlugins = append(finalPlugins, p)
	}

	return finalPlugins, nil
}
package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListHostPlugins(ctx context.Context, query orm.Query) (hPlugins []*model.HostPlugin, err error) {
	res := query.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&hPlugins)
	return hPlugins, res.Error
}

// ListHostsPluginsByHostId 根据HostId列出主机的所有插件
func (d *Dao) ListHostsPluginsByHostId(ctx context.Context, hostId string) (plugins []*model.Plugin, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.HostPlugin{}).
		Select("host_plugin.plugin_id AS id, "+
			"p.name AS name, "+
			"p.path AS path, "+
			"host_plugin.args AS args, "+
			"p.checksum AS `checksum`, "+
			"host_plugin.`interval` AS `interval`, "+
			"host_plugin.timeout AS timeout").
		Joins("LEFT JOIN plugin AS p ON p.id=host_plugin.plugin_id").
		Where("host_plugin.host_id=?", hostId).
		Find(&plugins)

	return plugins, res.Error
}

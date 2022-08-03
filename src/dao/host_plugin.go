package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) ListHostPlugins(query orm.Query) (hPlugins []*model.HostPlugin, err error) {
	db := query.Where(d.db)
	res := db.Limit(1).Find(&hPlugins)
	return hPlugins, res.Error
}

func (d *Dao) ListHostsPluginsByHostId(hostId string) (plugins []*model.Plugin, err error) {
	res := d.db.Model(&model.HostPlugin{}).
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

package dao

import (
	"owl/model"
)

func (d *Dao) ListHostGroupsPluginsByHostId(hostId string) (plugins []*model.Plugin, err error) {
	subQuery := d.db.Raw("SELECT host_group_id FROM host_group_host AS hgh WHERE hgh.host_id=?", hostId)
	res := d.db.Model(model.HostGroupPlugin{}).
		Select("host_group_plugin.plugin_id AS id, "+
			"p.name AS name, "+
			"p.path AS path, "+
			"host_group_plugin.args AS args, "+
			"p.checksum AS `checksum`, "+
			"host_group_plugin.interval AS `interval`, "+
			"host_group_plugin.timeout AS timeout").
		Joins("LEFT JOIN plugin AS p ON p.id=host_group_plugin.plugin_id").
		Where("host_group_plugin.group_id IN (?)", subQuery).
		Find(&plugins)
	return plugins, res.Error
}

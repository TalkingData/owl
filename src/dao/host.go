package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) SetOrNewHost(id, ip, hostname, agentVer string, uptime, idlePct float64) (obj *model.Host, err error) {
	res := d.db.Where(map[string]interface{}{
		"id":       id,
		"ip":       ip,
		"hostname": hostname,
	}).Attrs(&model.Host{
		Id: id,
	}).Assign(&model.Host{
		Name:         "",
		Ip:           ip,
		Hostname:     hostname,
		Uptime:       uptime,
		IdlePct:      idlePct,
		AgentVersion: agentVer,
	}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewHost(id, ip, hostname, agentVer string, uptime, idlePct float64) (*model.Host, error) {
	h := model.Host{
		Id:           id,
		Name:         "",
		Ip:           ip,
		Hostname:     hostname,
		Uptime:       uptime,
		IdlePct:      idlePct,
		AgentVersion: agentVer,
	}

	res := d.db.Create(&h)
	return &h, res.Error
}

func (d *Dao) SetHost(id, ip, hostname, agentVer string, uptime, idlePct float64) (h *model.Host, err error) {
	res := d.db.Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"ip":            ip,
			"hostname":      hostname,
			"agent_version": agentVer,
			"uptime":        uptime,
			"idle_pct":      idlePct,
		}).Find(&h)
	return h, res.Error
}

func (d *Dao) SetHostStatusByQuery(query orm.Query, status string) (int64, error) {
	// 这里务必使用d.db.Table("host")，这是为了避开gorm自动修改update_at字段
	res := query.Where(d.db.Table("host")).
		Updates(map[string]interface{}{
			"status": status,
		})
	return res.RowsAffected, res.Error
}

func (d *Dao) SetHostStatus(id, status string) (int64, error) {
	res := d.db.Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status": status,
		})
	return res.RowsAffected, res.Error
}

func (d *Dao) GetHost(query orm.Query) (h *model.Host, err error) {
	res := query.Where(d.db).Limit(1).Find(&h)
	return h, res.Error
}

func (d *Dao) ListHosts(query orm.Query) (hs []*model.Host, err error) {
	res := query.Where(d.db).Find(&hs)
	return hs, res.Error
}

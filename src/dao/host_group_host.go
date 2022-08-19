package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewHostGroupHost(hostGroupId uint, hostId string) (*model.HostGroupHost, error) {
	hgh := model.HostGroupHost{
		HostGroupId: hostGroupId,
		HostId:      hostId,
	}

	res := d.db.Create(&hgh)
	return &hgh, res.Error
}

func (d *Dao) GetHostGroupHost(hostGroupId uint, hostId string) (ph *model.HostGroupHost, err error) {
	res := d.db.Where(map[string]interface{}{
		"host_group_id": hostGroupId,
		"host_id":       hostId,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetHostGroupHostCount(hostGroupId uint, hostId string) (count int64, err error) {
	query := orm.Query{
		"host_group_id": hostGroupId,
		"host_id":       hostId,
	}
	res := query.Where(d.db.Model(&model.HostGroupHost{})).Count(&count)
	return count, res.Error
}

// ListHostsByHostGroupId 根据HostGroupId列出所有Host
func (d *Dao) ListHostsByHostGroupId(hostGroupId uint) (hosts []*model.Host, err error) {
	subQuery := d.db.Model(&model.HostGroupHost{}).
		Select("host_id").
		Where("host_group_id=?", hostGroupId)

	res := d.db.Where("id IN (?)", subQuery).
		Find(&hosts)

	return hosts, res.Error
}

func (d *Dao) IsHostInHostGroup(hostGroupId uint, hostId string) (exist bool, err error) {
	count, err := d.GetHostGroupHostCount(hostGroupId, hostId)
	return count > 0, err
}

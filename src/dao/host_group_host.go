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

func (d *Dao) GetHostGroupHost(hostGroupId uint, hostID string) (ph *model.HostGroupHost, err error) {
	res := d.db.Where(map[string]interface{}{
		"host_group_id": hostGroupId,
		"host_id":       hostID,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetHostGroupHostCount(hostGroupId uint, hostID string) (count int64, err error) {
	query := orm.Query{
		"host_group_id": hostGroupId,
		"host_id":       hostID,
	}
	res := query.Where(d.db.Model(&model.HostGroupHost{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsHostInHostGroup(hostGroupId uint, hostID string) (exist bool, err error) {
	count, err := d.GetHostGroupHostCount(hostGroupId, hostID)
	return count > 0, err
}

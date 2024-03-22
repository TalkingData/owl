package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewHostGroupHost(ctx context.Context, hostGroupId uint32, hostId string) (*model.HostGroupHost, error) {
	hgh := model.HostGroupHost{
		HostGroupId: hostGroupId,
		HostId:      hostId,
	}

	res := d.getDbWithCtx(ctx).Create(&hgh)
	return &hgh, res.Error
}

func (d *Dao) GetHostGroupHost(
	ctx context.Context,
	hostGroupId uint32, hostId string,
) (ph *model.HostGroupHost, err error) {
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{
		"host_group_id": hostGroupId,
		"host_id":       hostId,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetHostGroupHostCount(ctx context.Context, hostGroupId uint32, hostId string) (count int64, err error) {
	query := orm.Query{
		"host_group_id": hostGroupId,
		"host_id":       hostId,
	}
	res := query.Where(d.getDbWithCtx(ctx).Model(&model.HostGroupHost{})).Count(&count)
	return count, res.Error
}

// ListHostsByHostGroupId 根据HostGroupId列出所有Host
func (d *Dao) ListHostsByHostGroupId(ctx context.Context, hostGroupId uint32) (hosts []*model.Host, err error) {
	subQuery := d.getDbWithCtx(ctx).Model(&model.HostGroupHost{}).
		Select("host_id").
		Where("host_group_id=?", hostGroupId)

	res := d.getDbWithCtx(ctx).Where("id IN (?)", subQuery).
		Find(&hosts)

	return hosts, res.Error
}

func (d *Dao) IsHostInHostGroup(ctx context.Context, hostGroupId uint32, hostId string) (exist bool, err error) {
	count, err := d.GetHostGroupHostCount(ctx, hostGroupId, hostId)
	return count > 0, err
}

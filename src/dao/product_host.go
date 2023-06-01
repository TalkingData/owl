package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewProductHost(ctx context.Context, productId uint32, hostId string) (*model.ProductHost, error) {
	ph := model.ProductHost{
		ProductId: productId,
		HostId:    hostId,
	}

	res := d.getDbWithCtx(ctx).Create(&ph)
	return &ph, res.Error
}

func (d *Dao) GetProductHost(ctx context.Context, productId uint32, hostId string) (ph *model.ProductHost, err error) {
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{
		"product_id": productId,
		"host_id":    hostId,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetProductHostCount(ctx context.Context, productId uint32, hostId string) (count int64, err error) {
	query := orm.Query{
		"product_id": productId,
		"host_id":    hostId,
	}
	res := query.Where(d.getDbWithCtx(ctx).Model(&model.ProductHost{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsHostInProduct(ctx context.Context, productId uint32, hostId string) (exist bool, err error) {
	count, err := d.GetProductHostCount(ctx, productId, hostId)
	return count > 0, err
}

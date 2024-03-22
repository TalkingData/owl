package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetOrNewHostGroup(
	ctx context.Context,
	productId uint32, name, description, creator string,
) (obj *model.HostGroup, err error) {
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{
		"name":       name,
		"product_id": productId,
	}).Attrs(&model.HostGroup{
		Name:        name,
		Description: description,
		ProductId:   productId,
		Creator:     creator,
	}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewHostGroup(
	ctx context.Context,
	productId uint32, name, description, creator string,
) (*model.HostGroup, error) {
	hg := model.HostGroup{
		Name:        name,
		Description: description,
		ProductId:   productId,
		Creator:     creator,
	}

	res := d.getDbWithCtx(ctx).Create(&hg)
	return &hg, res.Error
}

func (d *Dao) GetHostGroup(ctx context.Context, q orm.Query) (hg *model.HostGroup, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&hg)
	return hg, res.Error
}

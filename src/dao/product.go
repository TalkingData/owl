package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetOrNewProduct(ctx context.Context, name, description, creator string) (obj *model.Product, err error) {
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{"name": name}).
		Attrs(&model.Product{
			Name:        name,
			Description: description,
			Creator:     creator,
			IsDelete:    false,
		}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewProduct(ctx context.Context, name, description, creator string) (*model.Product, error) {
	p := model.Product{
		Name:        name,
		Description: description,
		Creator:     creator,
		IsDelete:    false,
	}

	res := d.getDbWithCtx(ctx).Create(&p)
	return &p, res.Error
}

func (d *Dao) SetProduct(ctx context.Context, id uint32, name, description string) (p *model.Product, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).Find(&p)
	return p, res.Error
}

func (d *Dao) GetProduct(ctx context.Context, q orm.Query) (p *model.Product, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&p)
	return p, res.Error
}

func (d *Dao) GetProductCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Product{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsProductExist(ctx context.Context, q orm.Query) (exist bool, err error) {
	count, err := d.GetProductCount(ctx, q)
	return count > 0, err
}

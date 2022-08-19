package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetOrNewHostGroup(productId uint, name, description, creator string) (obj *model.HostGroup, err error) {
	res := d.db.Where(map[string]interface{}{
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

func (d *Dao) NewHostGroup(productId uint, name, description, creator string) (*model.HostGroup, error) {
	hg := model.HostGroup{
		Name:        name,
		Description: description,
		ProductId:   productId,
		Creator:     creator,
	}

	res := d.db.Create(&hg)
	return &hg, res.Error
}

func (d *Dao) GetHostGroup(query orm.Query) (hg *model.HostGroup, err error) {
	res := query.Where(d.db).Limit(1).Find(&hg)
	return hg, res.Error
}

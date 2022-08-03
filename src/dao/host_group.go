package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetOrNewHostGroup(productID uint, name, description, creator string) (obj *model.HostGroup, err error) {
	res := d.db.Where(map[string]interface{}{
		"name":       name,
		"product_id": productID,
	}).Attrs(&model.HostGroup{
		Name:        name,
		Description: description,
		ProductId:   productID,
		Creator:     creator,
	}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewHostGroup(productID uint, name, description, creator string) (*model.HostGroup, error) {
	hg := model.HostGroup{
		Name:        name,
		Description: description,
		ProductId:   productID,
		Creator:     creator,
	}

	res := d.db.Create(&hg)
	return &hg, res.Error
}

func (d *Dao) GetHostGroup(query orm.Query) (hg *model.HostGroup, err error) {
	db := query.Where(d.db)
	res := db.Limit(1).Find(&hg)
	return hg, res.Error
}

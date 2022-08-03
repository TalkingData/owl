package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetOrNewProduct(name, description, creator string) (obj *model.Product, err error) {
	res := d.db.Where(map[string]interface{}{"name": name}).
		Attrs(&model.Product{
			Name:        name,
			Description: description,
			Creator:     creator,
			IsDelete:    false,
		}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewProduct(name, description, creator string) (*model.Product, error) {
	p := model.Product{
		Name:        name,
		Description: description,
		Creator:     creator,
		IsDelete:    false,
	}

	res := d.db.Create(&p)
	return &p, res.Error
}

func (d *Dao) SetProduct(id uint, name, description string) (p *model.Product, err error) {
	res := d.db.Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).Find(&p)
	return p, res.Error
}

func (d *Dao) GetProduct(query orm.Query) (p *model.Product, err error) {
	db := query.Where(d.db)
	res := db.Limit(1).Find(&p)
	return p, res.Error
}

func (d *Dao) GetProductCount(query orm.Query) (count int64, err error) {
	db := query.Where(d.db.Model(&model.Product{}))
	res := db.Count(&count)
	return count, res.Error
}

func (d *Dao) IsProductExist(query orm.Query) (exist bool, err error) {
	count, err := d.GetProductCount(query)
	return count > 0, err
}

package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewProductHost(productId uint, hostId string) (*model.ProductHost, error) {
	ph := model.ProductHost{
		ProductId: productId,
		HostId:    hostId,
	}

	res := d.db.Create(&ph)
	return &ph, res.Error
}

func (d *Dao) GetProductHost(productId uint, hostId string) (ph *model.ProductHost, err error) {
	res := d.db.Where(map[string]interface{}{
		"product_id": productId,
		"host_id":    hostId,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetProductHostCount(productId uint, hostId string) (count int64, err error) {
	query := orm.Query{
		"product_id": productId,
		"host_id":    hostId,
	}
	res := query.Where(d.db.Model(&model.ProductHost{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsHostInProduct(productId uint, hostId string) (exist bool, err error) {
	count, err := d.GetProductHostCount(productId, hostId)
	return count > 0, err
}

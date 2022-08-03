package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) NewProductHost(productID uint, hostID string) (*model.ProductHost, error) {
	ph := model.ProductHost{
		ProductId: productID,
		HostId:    hostID,
	}

	res := d.db.Create(&ph)
	return &ph, res.Error
}

func (d *Dao) GetProductHost(productID uint, hostID string) (ph *model.ProductHost, err error) {
	res := d.db.Where(map[string]interface{}{
		"product_id": productID,
		"host_id":    hostID,
	}).Limit(1).Find(&ph)
	return ph, res.Error
}

func (d *Dao) GetProductHostCount(productID uint, hostID string) (count int64, err error) {
	query := orm.Query{
		"product_id": productID,
		"host_id":    hostID,
	}
	res := query.Where(d.db.Model(&model.ProductHost{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsHostInProduct(productID uint, hostID string) (exist bool, err error) {
	count, err := d.GetProductHostCount(productID, hostID)
	return count > 0, err
}

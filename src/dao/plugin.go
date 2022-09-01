package dao

import (
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetPlugin(query orm.Query) (p *model.Plugin, err error) {
	res := query.Where(d.db).Limit(1).Find(&p)
	return p, res.Error
}

func (d *Dao) GetPluginCount(query orm.Query) (count int64, err error) {
	res := query.Where(d.db.Model(&model.Plugin{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsPluginExist(query orm.Query) (exist bool, err error) {
	count, err := d.GetPluginCount(query)
	return count > 0, err
}
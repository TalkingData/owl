package dao

import (
	"context"
	"owl/common/orm"
	"owl/model"
)

func (d *Dao) GetPlugin(ctx context.Context, q orm.Query) (p *model.Plugin, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&p)
	return p, res.Error
}

func (d *Dao) GetPluginCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Plugin{})).Count(&count)
	return count, res.Error
}

func (d *Dao) IsPluginExist(ctx context.Context, q orm.Query) (exist bool, err error) {
	count, err := d.GetPluginCount(ctx, q)
	return count > 0, err
}

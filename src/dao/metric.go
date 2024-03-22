package dao

import (
	"context"
	"errors"
	"fmt"
	"owl/common/orm"
	"owl/model"
	"time"
)

func (d *Dao) CleanExpiredMetric(ctx context.Context, cycleExpiredRatio, limit, batchIntervalMs int) error {
	if limit < 1 {
		return errors.New("limit must greater than 0")
	}

	m := model.Metric{}
	sqlTemp := fmt.Sprintf("DELETE FROM %s "+
		"WHERE UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at) > (cycle * ?) "+
		"LIMIT ?", m.TableName())

	for {
		res := d.getDbWithCtx(ctx).Exec(sqlTemp, cycleExpiredRatio, limit)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected < int64(limit) {
			return nil
		}

		time.Sleep(time.Millisecond * time.Duration(batchIntervalMs))
	}
}

func (d *Dao) SetOrNetMetric(
	ctx context.Context,
	hostId, metric, tags, dt string,
	cycle int32,
) (obj *model.Metric, err error) {
	res := d.getDbWithCtx(ctx).Where(map[string]interface{}{
		"host_id": hostId,
		"metric":  metric,
		"tags":    tags,
	}).Assign(&model.Metric{
		HostId: hostId,
		Metric: metric,
		Tags:   tags,
		Dt:     dt,
		Cycle:  cycle,
	}).FirstOrCreate(&obj)

	return obj, res.Error
}

func (d *Dao) NewMetric(ctx context.Context, hostId, metric, tags, dt string, cycle int32) (*model.Metric, error) {
	m := model.Metric{
		HostId: hostId,
		Metric: metric,
		Tags:   tags,
		Dt:     dt,
		Cycle:  cycle,
	}

	res := d.getDbWithCtx(ctx).Create(&m)
	return &m, res.Error
}

func (d *Dao) SetMetric(
	ctx context.Context,
	id uint64,
	hostId, metric, tags, dt string,
	cycle int32,
) (m *model.Metric, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Metric{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"host_id": hostId,
			"metric":  metric,
			"tags":    tags,
			"dt":      dt,
			"cycle":   cycle,
		}).Find(&m)
	return m, res.Error
}

func (d *Dao) GetMetric(ctx context.Context, q orm.Query) (m *model.Metric, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&m)
	return m, res.Error
}

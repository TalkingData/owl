package dao

import (
	"errors"
	"owl/common/orm"
	"owl/model"
	"time"
)

func (d *Dao) CleanExpiredMetric(cycleExpiredRatio, limit, batchIntervalMs int) error {
	if limit < 1 {
		return errors.New("limit must greater than 0")
	}

	for {
		res := d.db.Raw("DELETE "+
			"FROM metric "+
			"WHERE UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at) > (cycle * ?) "+
			"LIMIT ?", cycleExpiredRatio, limit)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected < int64(limit) {
			return nil
		}

		time.Sleep(time.Millisecond * time.Duration(batchIntervalMs))
	}
}

func (d *Dao) SetOrNetMetric(hostId, metric, tags, dt string, cycle int) (obj *model.Metric, err error) {
	res := d.db.Where(map[string]interface{}{
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

func (d *Dao) NewMetric(hostId, metric, tags, dt string, cycle int) (*model.Metric, error) {
	m := model.Metric{
		HostId: hostId,
		Metric: metric,
		Tags:   tags,
		Dt:     dt,
		Cycle:  cycle,
	}

	res := d.db.Create(&m)
	return &m, res.Error
}

func (d *Dao) SetMetric(id uint64, hostId, metric, tags, dt string, cycle int) (m *model.Metric, err error) {
	res := d.db.Model(&model.Metric{}).
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

func (d *Dao) GetMetric(query orm.Query) (m *model.Metric, err error) {
	db := query.Where(d.db)
	res := db.Limit(1).Find(&m)
	return m, res.Error
}

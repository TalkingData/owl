package dao

import (
	"context"
	"fmt"
	"owl/common/logger"
	"owl/common/orm"
	"owl/model"
	"time"
)

func (d *Dao) SetOrNewHostById(
	ctx context.Context,
	id, ip, hostname, agentVer string,
	uptime, idlePct float64,
) (obj *model.Host, err error) {
	// 按照Id查找主机对象
	_o, e := d.GetHost(ctx, orm.Query{"id": id})
	if e != nil {
		return nil, e
	}
	// 找不到主机对象则创建并返回
	if _o == nil || len(_o.Id) < 1 {
		return d.NewHost(ctx, id, ip, hostname, agentVer, uptime, idlePct)
	}

	// 找到的主机对象如果主机对象的IP或主机名与传入参数不一致，则先记录日志
	if _o.Ip != ip || _o.Hostname != hostname {
		d.lg.WarnWithFields(logger.Fields{
			"id":                   id,
			"host_object_ip":       _o.Ip,
			"new_ip":               ip,
			"host_object_hostname": _o.Hostname,
			"new_hostname":         hostname,
		}, "Host object with the same ID has undergone changes in hostname or IP, and the corresponding fields have been updated.")
	}
	// 更新主机对象并返回
	return d.SetHost(ctx, id, ip, hostname, agentVer, uptime, idlePct)
}

func (d *Dao) NewHost(
	ctx context.Context,
	id, ip, hostname, agentVer string,
	uptime, idlePct float64,
) (*model.Host, error) {
	h := model.Host{
		Id:           id,
		Name:         "",
		Ip:           ip,
		Hostname:     hostname,
		Uptime:       uptime,
		IdlePct:      idlePct,
		AgentVersion: agentVer,
	}

	res := d.getDbWithCtx(ctx).Create(&h)
	return &h, res.Error
}

func (d *Dao) SetHost(
	ctx context.Context,
	id, ip, hostname, agentVer string,
	uptime, idlePct float64,
) (h *model.Host, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"ip":            ip,
			"hostname":      hostname,
			"agent_version": agentVer,
			"uptime":        uptime,
			"idle_pct":      idlePct,
		}).Find(&h)
	return h, res.Error
}

func (d *Dao) SetHostStatus2OkByThresholdSecs(
	ctx context.Context,
	thresholdSecs, limit, batchIntervalMs int,
) error {
	h := model.Host{}
	sqlTemp := fmt.Sprintf("UPDATE %s "+
		"SET status=? "+
		"WHERE status<>? AND (UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at)) <= ? "+
		"LIMIT ?", h.TableName())

	for {
		res := d.getDbWithCtx(ctx).Exec(sqlTemp, model.HostStatusOk, model.HostStatusOk, thresholdSecs, limit)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected < int64(limit) {
			return nil
		}

		time.Sleep(time.Millisecond * time.Duration(batchIntervalMs))
	}
}

func (d *Dao) SetHostStatus2DownByThresholdSecs(
	ctx context.Context,
	thresholdSecs, limit, batchIntervalMs int,
) error {
	h := model.Host{}
	sqlTemp := fmt.Sprintf("UPDATE %s "+
		"SET status=? "+
		"WHERE status=? AND (UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at)) > ? "+
		"LIMIT ?", h.TableName())

	for {
		res := d.getDbWithCtx(ctx).Exec(sqlTemp, model.HostStatusDown, model.HostStatusOk, thresholdSecs, limit)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected < int64(limit) {
			return nil
		}

		time.Sleep(time.Millisecond * time.Duration(batchIntervalMs))
	}
}

func (d *Dao) SetHostStatus(ctx context.Context, id, status string) (int64, error) {
	res := d.getDbWithCtx(ctx).Model(&model.Host{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status": status,
		})
	return res.RowsAffected, res.Error
}

func (d *Dao) GetHost(ctx context.Context, q orm.Query) (h *model.Host, err error) {
	res := q.Where(d.db).Limit(1).Find(&h)
	return h, res.Error
}

func (d *Dao) ListHosts(ctx context.Context, q orm.Query) (hs []*model.Host, err error) {
	res := q.Where(d.db).Find(&hs)
	return hs, res.Error
}

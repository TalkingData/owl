package biz

import (
	"context"
	"owl/common/logger"
)

func (b *Biz) ReceiveAgentHeartbeat(
	ctx context.Context,
	hostId, ip, hostname, version string,
	uptime, idlePct float64,
) error {
	// 准备更新或创建主机
	b.logger.DebugWithFields(logger.Fields{
		"agent_host_id":  hostId,
		"agent_ip":       ip,
		"agent_hostname": hostname,
		"agent_version":  version,
		"agent_uptime":   uptime,
		"agent_idle_pct": idlePct,
	}, "Biz.ReceiveAgentHeartbeat prepare execute dao.SetOrNewHostById.")
	_, err := b.dao.SetOrNewHostById(ctx, hostId, ip, hostname, version, uptime, idlePct)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  hostId,
			"agent_ip":       ip,
			"agent_hostname": hostname,
			"error":          err,
		}, "An error occurred while calling dao.SetOrNewHostById.")
		return err
	}

	return nil
}

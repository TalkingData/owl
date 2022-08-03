package biz

import "owl/common/logger"

func (b *Biz) ReceiveAgentHeartbeat(hostId, ip, hostname, version string, uptime, idlePct float64) error {
	// 准备更新或创建主机
	b.logger.InfoWithFields(logger.Fields{
		"agent_host_id":  hostId,
		"agent_ip":       ip,
		"agent_hostname": hostname,
		"agent_version":  version,
		"agent_uptime":   uptime,
		"agent_idle_pct": idlePct,
	}, "Biz.ReceiveAgentHeartbeat prepare execute dao.SetOrNewHost.")
	_, err := b.dao.SetOrNewHost(hostId, ip, hostname, version, uptime, idlePct)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"agent_host_id":  hostId,
			"agent_ip":       ip,
			"agent_hostname": hostname,
			"error":          err,
		}, "An error occurred while dao.SetOrNewHost.")
		return err
	}

	return nil
}

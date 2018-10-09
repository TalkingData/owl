package main

import (
	"database/sql"
	"fmt"
	"owl/common/types"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mydb *Storage

type Storage struct {
	*sqlx.DB
}

func InitMysqlConnPool() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		GlobalConfig.MySQLUser, GlobalConfig.MySQLPassword, GlobalConfig.MySQLAddr, GlobalConfig.MySQLDBName)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(GlobalConfig.MySQLMaxIdleConn)
	db.SetMaxOpenConns(GlobalConfig.MySQLMaxConn)
	mydb = &Storage{db}
	return nil
}

// func (s *Storage) createHost(host *types.Host) error {
// 	now := time.Now().Format(timeFomart)
// 	sqlString := fmt.Sprintf("insert into `host`(`id`, `ip`, `hostname`, `uptime`, `idle_pct`, `agent_version`, `create_at`, `update_at`) "+
// 		" values('%s', '%s', '%s', %0.2f, %0.2f, '%s','%s','%s')", host.ID, host.IP, host.Hostname, host.Uptime, host.IdlePct, host.AgentVersion, now, now)
// 	lg.Debug("create host:%s", sqlString)
// 	_, err := s.Exec(sqlString)
// 	return err
// }

// func (s *Storage) updateHost(host *types.Host) error {
// 	sqlString := fmt.Sprintf("update `host` set `ip`='%s', `uptime`=%0.2f, `idle_pct`=%0.2f, `hostname`='%s', `agent_version`='%s', `update_at`='%s' where id='%s'",
// 		host.IP, host.Uptime, host.IdlePct, host.Hostname, host.AgentVersion, time.Now().Format(timeFomart), host.ID)
// 	lg.Debug("update host:%s", sqlString)
// 	_, err := s.Exec(sqlString)
// 	return err
// }

func (s *Storage) createOrUpdateHost(host *types.Host) error {
	now := time.Now().Format(timeFomart)
	sqlString := fmt.Sprintf("insert into `host`(`id`, `ip`, `hostname`, `uptime`, `idle_pct`, `agent_version`, `create_at`, `update_at`) "+
		" values('%s', '%s', '%s', %0.2f, %0.2f, '%s','%s','%s') ON DUPLICATE key UPDATE ip=VALUES(ip), uptime=VALUES(uptime),"+
		"idle_pct=VALUES(idle_pct), hostname=VALUES(hostname), agent_version=VALUES(agent_version), update_at=VALUES(update_at)",
		host.ID, host.IP, host.Hostname, host.Uptime, host.IdlePct, host.AgentVersion, now, now)
	lg.Debug("create or update host:%s", sqlString)
	_, err := s.Exec(sqlString)
	return err
}

func (s *Storage) getHostByID(hostID string) (*types.Host, error) {
	host := &types.Host{}
	sqlString := fmt.Sprintf("select id, ip, hostname, agent_version,status,create_at, update_at  from `host` where id='%s'", hostID)
	lg.Debug("getHostByID:%s", sqlString)
	err := s.Get(host, sqlString)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return host, err
}
func (s *Storage) getHostByHostname(hostname string) (*types.Host, error) {
	host := &types.Host{}
	sqlString := fmt.Sprintf("select id, ip, hostname, agent_version,status,create_at, update_at  from `host` where hostname='%s'", hostname)
	lg.Debug("getHostByHostname:%s", sqlString)
	err := s.Get(host, sqlString)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return host, err
}

func (s *Storage) getAllHosts() []*types.Host {
	hosts := []*types.Host{}
	sqlString := fmt.Sprintf("select id, ip, hostname, agent_version,status,create_at, update_at  from `host`")
	if err := s.Select(&hosts, sqlString); err != nil {
		lg.Error("getNoMaintainHost %s", err)
		return nil
	}
	return hosts
}

func (s *Storage) setHostAlive(hostID string, status string) {
	sqlString := fmt.Sprintf("update `host` set `status` = '%s' where `id`='%s'", status, hostID)
	lg.Debug("setHostAlive:%s", sqlString)
	s.Exec(sqlString)
}

func (s *Storage) getHostPlugins(hostID string) ([]types.Plugin, error) {
	plugins := []types.Plugin{}
	idMap := make(map[string]types.Plugin)
	sqlString := fmt.Sprintf("select hp.id, p.name, p.path, hp.args, p.checksum, hp.interval, hp.timeout from host_plugin as hp "+
		" left join plugin as p on hp.plugin_id = p.id where host_id='%s'", hostID)
	lg.Debug("getHostPlugins:%s", sqlString)
	rows, err := s.Query(sqlString)
	if err != nil {
		lg.Error("getHostPlugins error:%v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		plugin := types.Plugin{}
		if err := rows.Scan(&plugin.ID, &plugin.Name, &plugin.Path, &plugin.Args, &plugin.Checksum, &plugin.Interval, &plugin.Timeout); err != nil {
			lg.Error("getHostPlugins error:%v", err)
			continue
		}
		plugins = append(plugins, plugin)
		idMap[plugin.UniqueKey()] = plugin
	}
	//获取主机组所有的插件
	sqlString = fmt.Sprintf("select hgp.id, p.name, p.path, hgp.args, p.checksum, hgp.interval, hgp.timeout from plugin as p "+
		" left join host_group_plugin as hgp on p.id = hgp.plugin_id where hgp.group_id in (select host_group_id from host_group_host where host_id='%s')", hostID)
	lg.Debug("getHostGroupPlugins:%s", sqlString)
	rows, err = s.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		plugin := types.Plugin{}
		if err := rows.Scan(&plugin.ID, &plugin.Name, &plugin.Path, &plugin.Args, &plugin.Checksum, &plugin.Interval, &plugin.Timeout); err != nil {
			lg.Error("getHostPlugins error:%v", err)
			continue
		}
		uniqueKey := plugin.UniqueKey()
		if p, ok := idMap[uniqueKey]; ok {
			lg.Warn("getHostPlugins: duplicate host group plugin (%v, %v)", plugin, p)
			continue
		}
		plugins = append(plugins, plugin)
		idMap[uniqueKey] = plugin
	}
	return plugins, nil
}

func (s *Storage) createOrUpdateMetric(hostID string, tsd types.TimeSeriesData) error {
	now := time.Now().Format(timeFomart)
	sqlString := fmt.Sprintf("insert into `metric` (`host_id`, `metric`, `tags`, `dt`, `cycle`, `create_at`, `update_at`) "+
		"values('%s', '%s', '%s', '%s', %d, '%s', '%s') ON DUPLICATE key UPDATE update_at = VALUES(update_at), cycle=VALUES(cycle)",
		hostID, tsd.Metric, tsd.Tags2String(), tsd.DataType, tsd.Cycle, now, now)
	lg.Debug("create or update metric:%s", sqlString)
	_, err := s.Exec(sqlString)
	return err
}

func (s *Storage) cleanupExpiredMetrics() error {
	rawSQL := fmt.Sprintf("delete from metric where UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(update_at) > (cycle * %d)",
		GlobalConfig.MetricExpiredCycle)
	lg.Info("cleanupExpiredMetrics:%s", rawSQL)
	_, err := s.Exec(rawSQL)
	return err
}

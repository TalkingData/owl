package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "owl/common/types"
	"time"
)

var mydb *db

type db struct {
	*sql.DB
}

func InitMysqlConnPool() error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		GlobalConfig.MYSQL_USER, GlobalConfig.MYSQL_PASSWORD, GlobalConfig.MYSQL_ADDR, GlobalConfig.MYSQL_DBNAME))
	if err != nil {
		return err
	}
	row, err := conn.Query("select 1")
	if err != nil {
		return err
	}
	defer row.Close()
	conn.SetMaxIdleConns(GlobalConfig.MYSQL_MAX_IDLE_CONN)
	conn.SetMaxOpenConns(GlobalConfig.MYSQL_MAX_CONN)
	mydb = &db{conn}
	return nil
}

func (this *db) CreateHost(id, sn, ip, hostname, agent_version string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := this.Exec("INSERT INTO `host`(`id`, `sn`, `ip`, `hostname`, `agent_version`,`create_at`, `update_at`) values(?,?,?,?,?,?,?)",
		id, sn, ip, hostname, agent_version, now, now)
	if err != nil {
		return err
	}
	return nil
}

func (this *db) UpdateHost(host *Host) {
	_, err := this.Exec("UPDATE `host` SET `ip`=?, `hostname`=?, `agent_version`=?, `update_at`=? WHERE id=?",
		host.IP, host.Hostname, host.AgentVersion, time.Now().Format("2006-01-02 15:04:05"), host.ID)
	if err != nil {
		lg.Error("update host error %s", err)
	}
}

func (this *db) DeleteHost(host *Host) error {
	_, err := this.Exec("DELETE FROM `host` WHERE ip=? and hostname=?", host.IP, host.Hostname)
	return err
}

func (this *db) GetHost(host_id string) *Host {
	host := &Host{ID: host_id}
	if err := this.QueryRow("SELECT `ip`, `hostname`, `agent_version`, `sn` FROM `host` WHERE id=?",
		host_id).Scan(&host.IP, &host.Hostname, &host.AgentVersion, &host.SN); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}
	return host
}

func (this *db) GetNoMaintainHost() []*Host {
	var (
		tm    time.Time
		err   error
		hosts []*Host
	)
	rows, err := this.Query("SELECT `id`, `ip`, `hostname`, `status`, `update_at` FROM `host` WHERE `status` != '2'")
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		h := &Host{}
		var update_at string
		if err := rows.Scan(&h.ID, &h.IP, &h.Hostname, &h.Status, &update_at); err != nil {
			continue
		}
		tm, err = time.ParseInLocation("2006-01-02 15:04:05", update_at, time.Local)
		if err != nil {
			continue
		}
		h.UpdateAt = tm
		hosts = append(hosts, h)
	}
	return hosts
}

func (this *db) SetHostAlive(id string, st string) {
	this.Exec("UPDATE `host` SET `status` = ? WHERE `id`=?", st, id)
}

func (this *db) GetPlugins(host_id string) []Plugin {
	plugins := []Plugin{}
	idMap := make(map[int]struct{})
	rows, err := this.Query("SELECT `id`, `name`, `args`, `interval`, `timeout` FROM `plugin` WHERE "+
		" id in(SELECT `plugin_id` FROM `host_plugin` WHERE `host_id`=?)", host_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		plugin := Plugin{}
		if err := rows.Scan(&plugin.ID, &plugin.Name, &plugin.Args, &plugin.Interval, &plugin.Timeout); err != nil {
			fmt.Println(err)
			continue
		}
		plugins = append(plugins, plugin)
		idMap[plugin.ID] = struct{}{}
	}
	//获取主机组所有的插件
	rows, err = this.Query("SELECT `id`, `name`, `args`, `interval`, `timeout` FROM `plugin` WHERE "+
		"id in (SELECT `plugin_id` FROM `group_plugin` WHERE group_id "+
		"in(SELECT `group_id` FROM `host_group` WHERE host_id=?))", host_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		plugin := Plugin{}
		if err := rows.Scan(&plugin.ID, &plugin.Name, &plugin.Args, &plugin.Interval, &plugin.Timeout); err != nil {
			fmt.Println(err)
			continue
		}
		if _, ok := idMap[plugin.ID]; ok {
			continue
		}
		plugins = append(plugins, plugin)
		idMap[plugin.ID] = struct{}{}
	}
	return plugins
}

func (this *db) MetricIsExists(host_id, metric string) bool {
	var id int
	err := this.QueryRow("SELECT `id` FROM `metric` WHERE `host_id`=? and `name`=? ", host_id, metric).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (this *db) CreateMetric(host_id string, tsd TimeSeriesData) error {
	_, err := this.Exec("INSERT INTO `metric` (`host_id`, `name`, `dt`, `cycle`) values(?,?,?,?)",
		host_id, tsd.GetMetric(), tsd.DataType, tsd.Cycle)
	return err
}

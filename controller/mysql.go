package main

import (
	"database/sql"
	"fmt"

	. "owl/common/types"

	_ "github.com/go-sql-driver/mysql"
)

var mydb *db

type db struct {
	*sql.DB
}

func InitMysqlConnPool() error {
	var err error
	var conn *sql.DB
	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		GlobalConfig.MYSQL_USER, GlobalConfig.MYSQL_PASSWORD, GlobalConfig.MYSQL_ADDR, GlobalConfig.MYSQL_DBNAME))
	if err != nil {
		return err
	}
	err = conn.Ping()
	if err != nil {
		return err
	}
	conn.SetMaxIdleConns(GlobalConfig.MYSQL_MAX_IDLE_CONN)
	conn.SetMaxOpenConns(GlobalConfig.MYSQL_MAX_CONN)
	mydb = &db{conn}
	return nil
}

func (this *db) GetStrategies() []*Strategy {
	strategies := []*Strategy{}
	var rows *sql.Rows
	var err error
	rows, err = this.Query("SELECT `id`, `name`, `priority`, `pid`, `alarm_count`, `cycle`, `expression`, `description`, `user_id`, `enable` FROM `strategy`")
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		strategy := Strategy{}
		if err := rows.Scan(&strategy.ID, &strategy.Name, &strategy.Priority, &strategy.Pid, &strategy.AlarmCount, &strategy.Cycle, &strategy.Expression, &strategy.Description, &strategy.UserID, &strategy.Enable); err != nil {
			lg.Error(err.Error())
			continue
		}
		strategies = append(strategies, &strategy)
	}
	return strategies
}

func (this *db) GetHostsByStrategyID(strategy_id int) []*Host {
	rows, err := this.Query("SELECT `id`, `name`, `ip`, `sn`, `hostname`, `status` FROM host WHERE id IN (SELECT host_id FROM strategy_host WHERE strategy_id=?)", strategy_id)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	hosts := []*Host{}
	for rows.Next() {
		host := Host{}
		if err := rows.Scan(&host.ID, &host.Name, &host.IP, &host.SN, &host.Hostname, &host.Status); err != nil {
			lg.Error(err.Error())
			continue
		}
		hosts = append(hosts, &host)
	}
	return hosts
}

func (this *db) GetGroupsByStrategyID(strategy_id int) []*Group {
	rows, err := this.Query("SELECT `id`, `name` FROM `group` WHERE `id` IN (SELECT `group_id` FROM `strategy_group` WHERE `strategy_id`=?)", strategy_id)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	groups := []*Group{}
	for rows.Next() {
		group := Group{}
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			lg.Error(err.Error())
			continue
		}
		groups = append(groups, &group)
	}
	return groups
}

func (this *db) GetHostsByGroupID(group_id int) []*Host {
	rows, err := this.Query("SELECT `id`, `name`, `ip`, `sn`, `hostname`, `status` FROM host WHERE id IN (SELECT host_id FROM host_group WHERE group_id = ?)", group_id)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	hosts := []*Host{}
	for rows.Next() {
		host := Host{}
		if err := rows.Scan(&host.ID, &host.Name, &host.IP, &host.SN, &host.Hostname, &host.Status); err != nil {
			lg.Error(err.Error())
			continue
		}
		hosts = append(hosts, &host)
	}
	return hosts
}

func (this *db) GetTriggersByStrategyID(strategy_id int) map[string]*Trigger {
	rows, err := this.Query("SELECT `id`, `strategy_id`, `metric`, `tags`, `number`, `index`, `name`, `method`, `symbol`, `threshold`, `description` FROM `trigger` WHERE `strategy_id` = ?", strategy_id)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	triggers := make(map[string]*Trigger)
	for rows.Next() {
		trigger := Trigger{}
		if err := rows.Scan(&trigger.ID, &trigger.StrategyID, &trigger.Metric, &trigger.Tags, &trigger.Number, &trigger.Index, &trigger.Name, &trigger.Method, &trigger.Symbol, &trigger.Threshold, &trigger.Description); err != nil {
			lg.Error(err.Error())
			continue
		}
		triggers[trigger.Index] = &trigger
	}
	return triggers
}

func (this *db) GetActions(strategy_id, action_type int) []*Action {
	rows, err := this.Query("SELECT `id`, `strategy_id`, `type`, `file_path`, `alarm_subject`, `restore_subject`, `alarm_template`, `restore_template`, `send_type` FROM `action` WHERE `strategy_id` = ? AND `type` = ?", strategy_id, action_type)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	actions := []*Action{}
	for rows.Next() {
		action := Action{}
		if err := rows.Scan(&action.ID, &action.StrategyID, &action.Type, &action.FilePath, &action.AlarmSubject, &action.RestoreSubject, &action.AlarmTemplate, &action.RestoreTemplate, &action.SendType); err != nil {
			lg.Error(err.Error())
			continue
		}
		actions = append(actions, &action)
	}
	return actions
}

func (this *db) CreateStrategyEvent(strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) (int64, error) {
	tx, err := this.Begin()
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}
	stmt, err := tx.Prepare("INSERT INTO `strategy_event` VALUES (0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	result, err := stmt.Exec(
		strategy_event.StrategyID,
		strategy_event.StrategyName,
		strategy_event.Priority,
		strategy_event.Cycle,
		strategy_event.AlarmCount,
		strategy_event.Expression,
		strategy_event.CreateTime,
		strategy_event.UpdateTime,
		strategy_event.Count,
		strategy_event.Status,
		strategy_event.HostID,
		strategy_event.HostCname,
		strategy_event.HostName,
		strategy_event.IP,
		strategy_event.SN,
		strategy_event.ProcessUser,
		strategy_event.ProcessComments,
		strategy_event.ProcessTime)
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	last_id, err := result.LastInsertId()
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	for _, trigger_event_set := range trigger_event_sets {
		for _, trigger_event := range trigger_event_set {
			stmt_insert, err := tx.Prepare("INSERT INTO `trigger_event` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				lg.Error(err.Error())
				return -1, err
			}
			_, err = stmt_insert.Exec(
				last_id,
				trigger_event.Index,
				trigger_event.Metric,
				trigger_event.Tags,
				trigger_event.Number,
				trigger_event.AggregateTags,
				trigger_event.CurrentThreshold,
				trigger_event.Method,
				trigger_event.Symbol,
				trigger_event.Threshold,
				trigger_event.Triggered)
			if err != nil {
				lg.Error(err.Error())
				return -1, err
			}
		}
	}
	defer tx.Rollback()

	tx.Commit()

	return last_id, nil
}

func (this *db) UpdateStrategyEvent(strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) error {
	tx, err := this.Begin()
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE `strategy_event` SET `update_time` = ?, `count` = ?, `status` = ?, `process_user` = ?, `process_comments` = ?, `process_time` = ? WHERE `id` = ?")
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(
		strategy_event.UpdateTime,
		strategy_event.Count,
		strategy_event.Status,
		strategy_event.ProcessUser,
		strategy_event.ProcessComments,
		strategy_event.ProcessTime,
		strategy_event.ID)
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	stmt_delete, err := tx.Prepare("DELETE FROM `trigger_event` WHERE `strategy_event_id` = ?")
	_, err = stmt_delete.Exec(strategy_event.ID)
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	for _, trigger_event_set := range trigger_event_sets {
		for _, trigger_event := range trigger_event_set {
			stmt_insert, err := tx.Prepare("INSERT INTO `trigger_event` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				lg.Error(err.Error())
				return err
			}
			_, err = stmt_insert.Exec(
				strategy_event.ID,
				trigger_event.Index,
				trigger_event.Metric,
				trigger_event.Tags,
				trigger_event.Number,
				trigger_event.AggregateTags,
				trigger_event.CurrentThreshold,
				trigger_event.Method,
				trigger_event.Symbol,
				trigger_event.Threshold,
				trigger_event.Triggered)
			if err != nil {
				lg.Error(err.Error())
				return err
			}
		}
	}
	defer tx.Rollback()

	tx.Commit()

	return nil
}

func (this *db) GetStrategyEvent(strategy_id, status int, host_id string) *StrategyEvent {
	strategy_event := &StrategyEvent{}
	if err := this.QueryRow("SELECT `id`, `strategy_id`, `strategy_name`, `priority`, `cycle`, `alarm_count`, `expression`, `create_time`, `update_time`, `count`, `status`, `host_id`, `host_cname`, `host_name`, `ip`, `sn`, `process_user`, `process_comments`, `process_time` FROM `strategy_event` WHERE `strategy_id` = ? AND `host_id` = ? AND `status` = ?",
		strategy_id,
		host_id,
		status).Scan(&strategy_event.ID,
		&strategy_event.StrategyID,
		&strategy_event.StrategyName,
		&strategy_event.Priority,
		&strategy_event.Cycle,
		&strategy_event.AlarmCount,
		&strategy_event.Expression,
		&strategy_event.CreateTime,
		&strategy_event.UpdateTime,
		&strategy_event.Count,
		&strategy_event.Status,
		&strategy_event.HostID,
		&strategy_event.HostCname,
		&strategy_event.HostName,
		&strategy_event.IP,
		&strategy_event.SN,
		&strategy_event.ProcessUser,
		&strategy_event.ProcessComments,
		&strategy_event.ProcessTime); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		lg.Error(err.Error())
		return nil
	}

	return strategy_event
}

func (this *db) GetUsersByGroups(action_id int) []*User {
	sql := fmt.Sprintf("SELECT `id`, `username`, `phone`, `mail`, `weixin` FROM `user` WHERE `id` IN (SELECT `user_id` FROM `user_user_group` WHERE `user_group_id` IN (SELECT `user_group_id` FROM action_user_group WHERE action_id = %d))", action_id)
	rows, err := this.Query(sql)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Phone, &user.Mail, &user.Weixin); err != nil {
			lg.Error(err.Error())
			continue
		}
		users = append(users, &user)
	}
	return users
}

func (this *db) GetUsers(action_id int) []*User {
	sql := fmt.Sprintf("SELECT `id`, `username`, `phone`, `mail`, `weixin` FROM `user` WHERE `id` IN (SELECT `user_id` FROM action_user WHERE action_id = %d)", action_id)
	rows, err := this.Query(sql)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	users_obj := []*User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Phone, &user.Mail, &user.Weixin); err != nil {
			lg.Error(err.Error())
			continue
		}
		users_obj = append(users_obj, &user)
	}
	return users_obj
}

func (this *db) CreateActionResult(action_result *ActionResult) error {
	stmt, err := this.Prepare("INSERT INTO `action_result` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	_, err = stmt.Exec(
		action_result.StrategyEventID,
		action_result.ActionID,
		action_result.ActionType,
		action_result.ActionSendType,
		action_result.UserID,
		action_result.Username,
		action_result.Phone,
		action_result.Mail,
		action_result.Weixin,
		action_result.Subject,
		action_result.Content,
		action_result.Success,
		action_result.Response)
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	return nil
}

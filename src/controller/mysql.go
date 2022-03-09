package main

import (
	"database/sql"
	"fmt"
	"owl/common/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mydb *db

type db struct {
	*sqlx.DB
}

func InitMysqlConnPool() error {
	var err error
	var conn *sqlx.DB
	conn, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		GlobalConfig.MYSQL_USER, GlobalConfig.MYSQL_PASSWORD, GlobalConfig.MYSQL_ADDR, GlobalConfig.MYSQL_DBNAME))
	if err != nil {
		return err
	}
	conn.SetMaxIdleConns(GlobalConfig.MYSQL_MAX_IDLE_CONN)
	conn.SetMaxOpenConns(GlobalConfig.MYSQL_MAX_CONN)
	mydb = &db{conn}
	return nil
}

func (this *db) GetProducts() []*types.Product {
	products := []*types.Product{}
	if err := this.Select(&products, "SELECT id, name FROM `product` WHERE is_delete=0"); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return products
}

func (this *db) GetDeletedProducts() []*types.Product {
	products := []*types.Product{}
	if err := this.Select(&products, "SELECT id, name FROM `product` WHERE is_delete=1"); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return products
}

func (this *db) GetStrategies(product_id int) []*types.Strategy {
	strategies := []*types.Strategy{}
	if err := this.Select(&strategies, "SELECT * FROM `strategy` WHERE product_id=?", product_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return strategies
}

func (this *db) GetHostsExByStrategyID(strategy_id int) map[string]string {
	ids := make([]string, 0)
	if err := this.Select(&ids, "SELECT host_id FROM strategy_host_exclude WHERE strategy_id=?", strategy_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	exHosts := make(map[string]string, 0)
	for _, id := range ids {
		exHosts[id] = id
	}
	return exHosts
}

func (this *db) GetGroupsByStrategyID(strategy_id int) []*types.Group {
	groups := []*types.Group{}
	if err := this.Select(&groups, "SELECT `id`, `name` FROM `host_group` WHERE `id` IN (SELECT `group_id` FROM `strategy_group` WHERE `strategy_id`=?)", strategy_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return groups
}

func (this *db) GetHostsByGroupID(group_id int) []*types.Host {
	hosts := []*types.Host{}
	if err := this.Select(&hosts, "SELECT `id`, `ip`, `hostname`, `status`, `mute_time` FROM host WHERE id IN (SELECT host_id FROM host_group_host WHERE host_group_id = ?)", group_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return hosts
}

func (this *db) GetTriggersByStrategyID(strategy_id int) map[string]*types.Trigger {
	rows, err := this.Query("SELECT `id`, `strategy_id`, `metric`, `tags`, `number`, `index`, `method`, `symbol`, `threshold`, `description` FROM `trigger` WHERE `strategy_id` = ?", strategy_id)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	defer rows.Close()
	triggers := make(map[string]*types.Trigger)
	for rows.Next() {
		trigger := types.Trigger{}
		if err := rows.Scan(&trigger.ID, &trigger.StrategyID, &trigger.Metric, &trigger.Tags, &trigger.Number, &trigger.Index, &trigger.Method, &trigger.Symbol, &trigger.Threshold, &trigger.Description); err != nil {
			lg.Error(err.Error())
			continue
		}
		triggers[trigger.Index] = &trigger
	}
	return triggers
}

func (this *db) GetActions(strategy_id, action_type int) []*types.Action {
	actions := []*types.Action{}
	if err := this.Select(&actions, "SELECT * FROM `action` WHERE `strategy_id` = ? AND `type` = ?", strategy_id, action_type); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return actions
}

func (this *db) GetAllActions(strategy_id int) []*types.Action {
	actions := []*types.Action{}
	if err := this.Select(&actions, "SELECT * FROM `action` WHERE `strategy_id` = ?", strategy_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return actions
}

func (this *db) CreateStrategyEvent(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) (int64, error) {
	tx, err := this.Begin()
	defer tx.Rollback()
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}
	stmt, err := tx.Prepare("INSERT INTO `strategy_event` VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?, DEFAULT, ?, ?, ?, ?, ?)")
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	result, err := stmt.Exec(
		strategy_event.ProductID,
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
		strategy_event.HostName,
		strategy_event.IP)
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	last_id, err := result.LastInsertId()
	if err != nil {
		lg.Error(err.Error())
		return -1, err
	}

	for _, trigger_event := range trigger_events {
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
	tx.Commit()

	return last_id, nil
}

func (this *db) CreateStrategyEventFailed(strategyID int, hostID string, status int, message string) error {
	_, err := this.Exec("REPLACE INTO strategy_event_failed VALUES (?, ?, ?, ?, DEFAULT, DEFAULT)", strategyID, hostID, status, message)
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	return err
}

func (this *db) UpdateStrategyEvent(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) error {
	tx, err := this.Begin()
	defer tx.Rollback()
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE `strategy_event` SET `strategy_name` = ?, `priority` = ?, `cycle` = ?, `alarm_count` = ?, `expression` = ?, `update_time` = ?, `count` = ?, `status` = ?, `host_id` = ?, `host_name` = ?, `ip` = ? WHERE `id` = ?")
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(
		strategy_event.StrategyName,
		strategy_event.Priority,
		strategy_event.Cycle,
		strategy_event.AlarmCount,
		strategy_event.Expression,
		strategy_event.UpdateTime,
		strategy_event.Count,
		strategy_event.Status,
		strategy_event.HostID,
		strategy_event.HostName,
		strategy_event.IP,
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
	for _, trigger_event := range trigger_events {
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

	tx.Commit()

	return nil
}

func (this *db) GetStrategyEvent(strategy_id, status int, host_id string) *types.StrategyEvent {
	strategy_event := types.StrategyEvent{}
	if err := this.Get(&strategy_event, "SELECT * FROM `strategy_event` WHERE `strategy_id` = ? AND `host_id` = ? AND `status` = ?", strategy_id, host_id, status); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		lg.Error(err.Error())
		return nil
	}
	return &strategy_event
}

func (this *db) GetTriggeredTriggerEvents(strategy_event_id int64) []*types.TriggerEvent {
	trigger_events := []*types.TriggerEvent{}
	if err := this.Select(&trigger_events, "SELECT * FROM trigger_event WHERE strategy_event_id = ? and triggered=1", strategy_event_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return trigger_events
}

func (this *db) GetUsersByGroups(action_id int) []*types.User {
	users := []*types.User{}
	if err := this.Select(&users, "SELECT `id`, `username`, `phone`, `mail`, `wechat` FROM `user` WHERE `id` IN (SELECT `user_id` FROM `user_group_user` WHERE `user_group_id` IN (SELECT `user_group_id` FROM action_user_group WHERE action_id = ?))", action_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return users
}

func (this *db) GetUsers(action_id int) []*types.User {
	users := []*types.User{}
	if err := this.Select(&users, "SELECT `id`, `username`, `phone`, `mail`, `wechat` FROM `user` WHERE `id` IN (SELECT `user_id` FROM action_user WHERE action_id = ?)", action_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return users
}

func (this *db) CreateActionResult(action_result *types.ActionResult) error {
	stmt, err := this.Prepare("INSERT INTO `action_result` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	_, err = stmt.Exec(
		action_result.StrategyEventID,
		action_result.Count,
		action_result.ActionID,
		action_result.ActionType,
		action_result.ActionKind,
		action_result.ScriptID,
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

func (this *db) CreateStrategyEventRecord(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) error {
	tx, err := this.Begin()
	defer tx.Rollback()
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO `strategy_event_record` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	_, err = stmt.Exec(
		strategy_event.ID,
		strategy_event.Count,
		strategy_event.StrategyID,
		strategy_event.StrategyName,
		strategy_event.Priority,
		strategy_event.Cycle,
		strategy_event.AlarmCount,
		strategy_event.Expression,
		strategy_event.CreateTime,
		strategy_event.UpdateTime,
		strategy_event.AwareEndTime,
		strategy_event.Status,
		strategy_event.HostID,
		strategy_event.HostName,
		strategy_event.IP)
	if err != nil {
		lg.Error(err.Error())
		return err
	}

	for _, trigger_event := range trigger_events {
		stmt_insert, err := tx.Prepare("INSERT INTO `trigger_event_record` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			lg.Error(err.Error())
			return err
		}
		_, err = stmt_insert.Exec(
			strategy_event.ID,
			strategy_event.Count,
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

	tx.Commit()

	return nil
}

func (this *db) GetScript(script_id int) *types.Script {
	script := types.Script{}
	if err := this.Get(&script, "SELECT * FROM `scripts` WHERE `id` = ?", script_id); err != nil {
		lg.Error(err.Error())
		return nil
	}
	return &script
}

func (this *db) DeleteStrategyFailed(strategy_id int, host_id string) error {
	_, err := this.Exec("DELETE FROM strategy_event_failed WHERE strategy_id = ? AND host_id = ?", strategy_id, host_id)
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	return err
}

func (this *db) CreateStrategyEventProcess(strategy_event_id int64, strategy_event_status int, process_user, process_comments, process_time string) error {
	_, err := this.Exec("INSERT INTO strategy_event_process VALUES(?, ?, ?, ?, DEFAULT)", strategy_event_id, strategy_event_status, process_user, process_comments)
	if err != nil {
		lg.Error(err.Error())
		return err
	}
	return err
}

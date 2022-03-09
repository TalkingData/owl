package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"owl/common/types"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mydb *db

const (
	dbDateFormat = "%Y-%m-%d %H:%i:%s"
	timeFormat   = "2006-01-02 15:04:05"
)

// InitMysqlConnPool 初始化数据库连接池
func InitMysqlConnPool() error {
	var err error
	var conn *sqlx.DB
	conn, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.MySQLUser, config.MySQLPassword, config.MySQLAddr, config.MySQLDBName))
	if err != nil {
		return err
	}
	conn.SetMaxIdleConns(config.MySQLMaxIdleConn)
	conn.SetMaxOpenConns(config.MySQLMaxConn)
	mydb = &db{conn}
	return nil
}

type db struct {
	*sqlx.DB
}

// ---------------------------------- Operation --------------------------------------------------
func (d *db) CreateOperation(operation *Operation) (err error) {
	rawSQL := fmt.Sprintf("INSERT INTO operations values('%s', '%s', '%s', '%s', '%s', %d, DEFAULT)",
		operation.IP, operation.Operator, operation.Method, operation.API, operation.Body, operation.Result)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
	}
	return
}

func (d *db) GetOperations(where, limit string) []*Operation {
	operations := []*Operation{}
	rawSQL := fmt.Sprintf("SELECT * FROM operations WHERE %s ORDER BY time desc LIMIT %s", where, limit)
	if err := d.Select(&operations, rawSQL); err != nil {
		log.Println(err)
	}
	return operations
}

func (d *db) GetOperationsCount(where string) int {
	var total int
	rawSQL := fmt.Sprintf("SELECT COUNT(*) FROM operations WHERE %s", where)
	if err := d.Get(&total, rawSQL); err != nil {
		log.Println(err)
	}
	return total
}

// ---------------------------------- Strategy --------------------------------------------------

// GetStrategies 获取策略列表
func (d *db) GetStrategies(where, order, limit string) []*StrategySummary {
	strategies := []*StrategySummary{}
	rawSQL := `SELECT s.*, IFNULL(user.username,'') user_name, IFNULL(se.alert_count, 0) alert_count, IFNULL(sef.nodata_count, 0) nodata_count, IFNULL(sef.unknow_count, 0) unknow_count
	FROM strategy s LEFT JOIN user ON s.user_id = user.id LEFT JOIN (select strategy_id, sum(case when status=4 then 1 else 0 end) as nodata_count, sum(case when status=5 then 1 else 0 end) as unknow_count
	from strategy_event_failed group by strategy_id) sef ON s.id = sef.strategy_id LEFT JOIN (select strategy_id, sum(case when status=1 then 1 else 0 end) as alert_count from strategy_event group by strategy_id) se ON s.id = se.strategy_id GROUP BY s.id`
	if len(where) > 0 {
		rawSQL = fmt.Sprintf("%s HAVING %s", rawSQL, where)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s ORDER BY %s", rawSQL, order)
	}
	if len(limit) > 0 {
		rawSQL = fmt.Sprintf("%s LIMIT %s", rawSQL, limit)
	}
	log.Println(rawSQL)
	if err := d.Select(&strategies, rawSQL); err != nil {
		log.Println(err)
	}
	return strategies
}

// GetStrategiesCount 获取策略数量
func (d *db) GetStrategiesCount(where string) int {
	var total int
	rawSQL := `SELECT s.*, IFNULL(user.username, '') user_name, IFNULL(se.alert_count, 0) alert_count, IFNULL(sef.nodata_count, 0) nodata_count, IFNULL(sef.unknow_count, 0) unknow_count
	FROM strategy s LEFT JOIN user ON s.user_id = user.id LEFT JOIN (select strategy_id, sum(case when status=4 then 1 else 0 end) as nodata_count, sum(case when status=5 then 1 else 0 end) as unknow_count
	from strategy_event_failed group by strategy_id) sef ON s.id = sef.strategy_id LEFT JOIN (select strategy_id, sum(case when status=1 then 1 else 0 end) as alert_count from strategy_event group by strategy_id) se ON s.id = se.strategy_id GROUP BY s.id`
	rawSQL = fmt.Sprintf("SELECT COUNT(*) FROM (%s HAVING %s) as s_se_sef", rawSQL, where)
	if err := d.Get(&total, rawSQL); err != nil {
		log.Println(err)
	}
	return total
}

// GetStrategiesByHostGroupID 获取主机组下的所有策略
func (d *db) GetStrategiesByHostGroupID(where, limit string) []*StrategySimple {
	strategies := []*StrategySimple{}
	rawSQL := `SELECT s.*, IFNULL(u.username, '') user_name FROM strategy_group sg LEFT JOIN strategy s ON sg.strategy_id = s.id LEFT JOIN user u ON s.user_id = u.id WHERE %s LIMIT %s`
	rawSQL = fmt.Sprintf(rawSQL, where, limit)
	if err := d.Select(&strategies, rawSQL); err != nil {
		log.Println(err)
	}
	return strategies
}

// GetStrategiesByHostGroupIDCount 获取主机组下的所有策略个数
func (d *db) GetStrategiesByHostGroupIDCount(where string) int {
	var total int
	rawSQL := `SELECT COUNT(*) FROM strategy_group sg LEFT JOIN strategy s ON sg.strategy_id = s.id LEFT JOIN user u ON s.user_id = u.id WHERE %s`
	rawSQL = fmt.Sprintf(rawSQL, where)
	if err := d.Get(&total, rawSQL); err != nil {
		log.Println(err)
	}
	return total
}

// GetStrategy 获取单个策略
func (d *db) GetStrategy(id int64, productID int) *Strategy {
	strategy := Strategy{}
	if err := d.Get(&strategy, "SELECT * FROM strategy WHERE id=? and product_id=?", id, productID); err != nil {
		log.Println(err)
	}
	return &strategy
}

// CreateStrategy 创建一条策略
func (d *db) CreateStrategy(strategy *StrategyDetail) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			log.Println(err)
		}
	}()
	tx := d.MustBegin()
	defer tx.Rollback()
	r := tx.MustExec("INSERT INTO strategy VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		strategy.ProductID, strategy.Name, strategy.Priority, strategy.AlarmCount, strategy.Cycle, strategy.Expression, strategy.Description, strategy.Enable, strategy.UserID)
	strategyID, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	for _, sg := range strategy.Groups {
		tx.MustExec("INSERT INTO strategy_group VALUES (DEFAULT, ?, ?)", strategyID, sg.GroupID)
	}
	for _, se := range strategy.ExcludeHosts {
		tx.MustExec("INSERT INTO strategy_host_exclude VALUES (DEFAULT, ?, ?)", strategyID, se.ID)
	}
	for _, tr := range strategy.Triggers {
		tx.MustExec("INSERT INTO `trigger` VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			strategyID, tr.Metric, tr.Tags, tr.Number, tr.Index, tr.Method, tr.Symbol, tr.Threshold, tr.Description)
	}
	for _, ac := range strategy.Actions {
		a := tx.MustExec("INSERT INTO action VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			strategyID, ac.Type, ac.Kind, ac.AlarmSubject, ac.AlarmTemplate, ac.RestoreSubject, ac.RestoreTemplate, ac.ScriptID, ac.BeginTime, ac.EndTime, ac.TimePeriod)
		actionID, err := a.LastInsertId()
		if err != nil {
			log.Println(err)
			return err
		}
		for _, g := range ac.UserGroups {
			tx.MustExec("INSERT INTO action_user_group VALUES (DEFAULT, ?, ?)", actionID, g.UserGroupID)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return
}

// DeleteStrategies 批量删除策略
func (d *db) DeleteStrategies(ids []string, productID int) (err error) {
	rawSQL := fmt.Sprintf("DELETE FROM strategy WHERE id in (%s) and product_id=%d", strings.Join(ids, ","), productID)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
	}
	return
}

// GetHostGroupsByStrategyID 获取策略下所监控的主机组
func (d *db) GetHostGroupsByStrategyID(id int64) []*types.Group {
	groups := []*types.Group{}
	if err := d.Select(&groups, "SELECT host_group.id, name FROM host_group JOIN strategy_group ON strategy_group.group_id = host_group.id WHERE strategy_id = ?", id); err != nil {
		log.Println(err)
	}
	return groups
}

// GetTriggersByStrategyID 获取策略下的逻辑表达式
func (d *db) GetTriggersByStrategyID(id int64) []*Trigger {
	triggers := []*Trigger{}
	if err := d.Select(&triggers, "SELECT * FROM `trigger` WHERE strategy_id = ? ORDER BY `index` ASC", id); err != nil {
		log.Println(err)
	}
	return triggers
}

// GetActionsByStrategyID 获取策略下的动作列表
func (d *db) GetActionsByStrategyID(id int64) []*ActionInfo {
	actions := []*ActionInfo{}
	if err := d.Select(&actions, "SELECT * FROM action WHERE strategy_id = ?", id); err != nil {
		log.Println(err)
	}
	return actions
}

// GetUserGroupsByActionID 获取动作下的用户组
func (d *db) GetUserGroupsByActionID(id int) []*types.UserGroup {
	groups := []*types.UserGroup{}
	if err := d.Select(&groups, "SELECT user_group.id, user_group.name FROM user_group JOIN action_user_group ON user_group.ID = action_user_group.user_group_id WHERE action_id = ?", id); err != nil {
		log.Println(err)
	}
	return groups
}

// UpdateStrategy 修改策略
func (d *db) UpdateStrategy(strategy *StrategyDetail) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			log.Println(err)
		}
	}()
	tx := d.MustBegin()
	defer tx.Rollback()
	tx.MustExec("UPDATE strategy SET name=?, priority=?, alarm_count=?, cycle=?, expression=?, description=?, enable=? WHERE id=?",
		strategy.Name, strategy.Priority, strategy.AlarmCount, strategy.Cycle, strategy.Expression, strategy.Description, strategy.Enable, strategy.ID)
	tx.MustExec("DELETE FROM strategy_group WHERE strategy_id = ?", strategy.ID)
	for _, sg := range strategy.Groups {
		tx.MustExec("INSERT INTO strategy_group VALUES (0, ?, ?)", strategy.ID, sg.GroupID)
	}
	tx.MustExec("DELETE FROM strategy_host_exclude WHERE strategy_id = ?", strategy.ID)
	for _, se := range strategy.ExcludeHosts {
		tx.MustExec("INSERT INTO strategy_host_exclude VALUES (0, ?, ?)", strategy.ID, se.ID)
	}
	tx.MustExec("DELETE FROM `trigger` WHERE strategy_id = ?", strategy.ID)
	for _, tr := range strategy.Triggers {
		tx.MustExec("INSERT INTO `trigger` VALUES (0, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			strategy.ID, tr.Metric, tr.Tags, tr.Number, tr.Index, tr.Method, tr.Symbol, tr.Threshold, tr.Description)
	}
	tx.MustExec("DELETE FROM action WHERE strategy_id = ?", strategy.ID)
	for _, ac := range strategy.Actions {
		a := tx.MustExec("INSERT INTO action VALUES (0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			strategy.ID, ac.Type, ac.Kind, ac.AlarmSubject, ac.AlarmTemplate, ac.RestoreSubject, ac.RestoreTemplate, ac.ScriptID, ac.BeginTime, ac.EndTime, ac.TimePeriod)
		actionID, err := a.LastInsertId()
		if err != nil {
			return err
		}
		for _, g := range ac.UserGroups {
			tx.MustExec("INSERT INTO action_user_group VALUES (0, ?, ?)", actionID, g.UserGroupID)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return
}

// UpdateStrategiesStatus 修改策略状态
func (d *db) UpdateStrategiesStatus(ids []string, productID int, enable string) (err error) {
	rawSQL := fmt.Sprintf("UPDATE strategy SET enable=%s WHERE id in (%s) and product_id=%d", enable, strings.Join(ids, ","), productID)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return
}

// ---------------------------------- Strategy Event --------------------------------------------------

// GetStrategyEvents 获取报警事件列表
func (d *db) GetStrategyEvents(where, order, limit string) []*StrategyEvent {
	events := []*StrategyEvent{}
	rawSQL := "SELECT * FROM strategy_event"
	if limit != "" && where != "" {
		rawSQL = fmt.Sprintf("%s WHERE %s ORDER BY %s LIMIT %s", rawSQL, where, order, limit)
	} else if where != "" {
		rawSQL = fmt.Sprintf("%s WHERE %s ORDER BY %s", rawSQL, where, order)
	} else {
		rawSQL = fmt.Sprintf("%s ORDER BY %s", rawSQL, order)
	}
	if err := d.Select(&events, rawSQL); err != nil {
		log.Println(err)
	}
	return events
}

// GetStrategiesEventsCount 获取报警事件个数
func (d *db) GetStrategiesEventsCount(where string) int {
	var total int
	rawSQL := "SELECT count(*) FROM strategy_event"
	rawSQL = fmt.Sprintf("%s WHERE %s", rawSQL, where)
	if err := d.Get(&total, rawSQL); err != nil {
		log.Println(err)
	}
	return total
}

// UpdateStrategyEventsStatus 修改报警事件状态
func (d *db) UpdateStrategyEventsStatus(ids []string, awareEndTime string, productID, status int) (err error) {
	rawSQL := fmt.Sprintf("UPDATE strategy_event SET status=%d, aware_end_time='%s' WHERE id in (%s) AND product_id=%d AND status != %d", status, awareEndTime, strings.Join(ids, ","), productID, types.EVENT_CLOSED)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return
}

// CreateStrategyEventProcesses 创建报警事件处理记录
func (d *db) CreateStrategyEventProcesses(strategyEventIDs []string, strategyEventStatus int, processUser, processComments string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			fmt.Println(err)
		}
	}()
	tx := d.MustBegin()
	defer tx.Rollback()
	for _, eventID := range strategyEventIDs {
		tx.MustExec("INSERT INTO strategy_event_process VALUES(?, ?, ?, ?, DEFAULT)", eventID, strategyEventStatus, processUser, processComments)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return
}

func (d *db) CleanupHostEvents(hostID string) {
	var err error
	if err = d.cleanupHostEvent(hostID); err != nil {
		log.Println("clean-up host event failed, error:", err.Error())
	}
	if err = d.cleanupHostEventFailed(hostID); err != nil {
		log.Println("clean-up host failed event failed, error:", err.Error())
	}
}

func (d *db) cleanupHostEventFailed(hostID string) error {
	rawSQL := fmt.Sprintf("delete from strategy_event_failed where host_id='%s'", hostID)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

func (d *db) cleanupHostEvent(hostID string) error {
	rawSQL := fmt.Sprintf("delete from strategy_event where host_id='%s'", hostID)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

// GetStrategyEventsFailed 获取失败的报警事件
func (d *db) GetStrategyEventsFailed(where, order, limit string) []*StrategyEventFailed {
	eventsFailed := []*StrategyEventFailed{}
	rawSQL := "SELECT sef.status, sef.update_time, sef.message, h.hostname, h.ip FROM strategy_event_failed sef LEFT JOIN host h ON sef.host_id = h.id"
	if limit != "" {
		rawSQL = fmt.Sprintf("%s WHERE %s ORDER BY %s LIMIT %s", rawSQL, where, order, limit)
	} else {
		rawSQL = fmt.Sprintf("%s WHERE %s ORDER BY %s", rawSQL, where, order)
	}
	log.Println(rawSQL)
	if err := d.Select(&eventsFailed, rawSQL); err != nil {
		log.Println(err)
	}
	return eventsFailed
}

// GetStrategyEventsFailedCount 获取失败的报警事件总数
func (d *db) GetStrategyEventsFailedCount(where string) int {
	var total int
	rawSQL := "SELECT count(*) FROM strategy_event_failed sef LEFT JOIN host h ON sef.host_id = h.id"
	rawSQL = fmt.Sprintf("%s WHERE %s", rawSQL, where)
	log.Println(rawSQL)
	if err := d.Get(&total, rawSQL); err != nil {
		log.Println(err)
	}
	return total
}

// GetStrategyEventProcessRecord 获取处理报警事件的处理记录
func (d *db) GetStrategyEventProcessRecord(eventID int64) []*StrategyEventProcess {
	strategyEventProcess := []*StrategyEventProcess{}
	rawSQL := "SELECT process_status, process_user, process_comments, process_time FROM strategy_event_process WHERE strategy_event_id = ? ORDER BY process_time DESC"
	if err := d.Select(&strategyEventProcess, rawSQL, eventID); err != nil {
		log.Println(err)
	}
	return strategyEventProcess
}

// GetAlarmRecord 获取报警历史记录
func (d *db) GetAlarmRecords(eventID int64, order, limit string) (records []*AlarmRecord, total int) {
	records = make([]*AlarmRecord, 0)
	events := []*StrategyEventRecord{}
	rawSQL := fmt.Sprintf("SELECT count(*) FROM strategy_event_record WHERE strategy_event_id = ?")
	if err := d.Get(&total, rawSQL, eventID); err != nil {
		log.Println(err)
		return
	}
	rawSQL = fmt.Sprintf("SELECT * FROM strategy_event_record WHERE strategy_event_id = ? ORDER BY %s LIMIT %s", order, limit)
	if err := d.Select(&events, rawSQL, eventID); err != nil {
		log.Println(err)
		return
	}
	for _, event := range events {
		record := new(AlarmRecord)
		record.StrategyEvent = event
		record.TriggerEvents = d.GetTriggersRecords(event.StrategyEventID, event.Count)
		record.ActionResults = d.GetActionResults(event.StrategyEventID, event.Count)
		records = append(records, record)
	}
	return
}

// GetTriggersRecords 获取报警事件下的表达式组
func (d *db) GetTriggersRecords(eventID int64, count int) []*TriggerEventRecord {
	triggers := []*TriggerEventRecord{}
	rawSQL := "SELECT * FROM trigger_event_record WHERE strategy_event_id = ? AND count = ?"
	if err := d.Select(&triggers, rawSQL, eventID, count); err != nil {
		log.Println(err)
		return nil
	}
	return triggers
}

// GetActionResults 获取报警事件下的报警结果组
func (d *db) GetActionResults(eventID int64, count int) []*ActionResult {
	results := []*ActionResult{}
	rawSQL := "SELECT action_result.*, scripts.name as script_name, scripts.file_path as file_path FROM action_result LEFT JOIN scripts ON action_result.script_id = scripts.id WHERE strategy_event_id = ? AND count = ?"
	if err := d.Select(&results, rawSQL, eventID, count); err != nil {
		log.Println(err)
		return nil
	}
	return results
}

// GetHostsExByStrategyID 获取排除的主机列表
func (d *db) GetHostsExByStrategyID(strategyID int64) []*AlarmHost {
	hosts := []*AlarmHost{}
	if err := d.Select(&hosts, "SELECT h.id, h.ip, h.hostname FROM strategy_host_exclude she LEFT JOIN host h ON she.host_id = h.id WHERE strategy_id=?;", strategyID); err != nil {
		log.Println(err)
	}
	return hosts
}

// ---------------------------------- Script --------------------------------------------------

// GetScriptByScriptID 获取单个动作脚本
func (d *db) GetScriptByScriptID(id int) *Script {
	script := Script{}
	if err := d.Get(&script, "SELECT id, name, file_path FROM scripts WHERE id=?", id); err != nil {
		log.Println(err)
	}
	return &script
}

// GetScripts 获取脚本列表
func (d *db) GetScripts() []*Script {
	scripts := []*Script{}
	if err := d.Select(&scripts, "SELECT id, name, file_path FROM scripts"); err != nil {
		log.Println(err)
	}
	return scripts
}

// CreateScript 创建脚本
func (d *db) CreateScript(script *Script) (err error) {
	rawSQL := fmt.Sprintf("INSERT INTO scripts VALUES(DEFAULT, '%s', '%s')", script.Name, script.FilePath)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return
}

// UpdateScript 修改脚本
func (d *db) UpdateScript(script *Script) (err error) {
	rawSQL := fmt.Sprintf("UPDATE scripts set name='%s', file_path='%s' WHERE id='%d'", script.Name, script.FilePath, script.ID)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return

}

// DeleteScript 删除脚本
func (d *db) DeleteScript(scriptIDs []string) (err error) {
	rawSQL := fmt.Sprintf("DELETE FROM scripts WHERE id in (%s)", strings.Join(scriptIDs, ","))
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return
}

// ---------------------------------- StrategyTemplate --------------------------------------------------

// GetStrategyTemplateDetail 获取单个策略模板的详细
func (d *db) GetStrategyTemplateDetail(strategyTemplateID int) *StrategyTemplateDetail {
	strategy := d.GetStrategyTemplate(strategyTemplateID)
	triggers := d.GetTriggerTemplates(strategyTemplateID)
	std := new(StrategyTemplateDetail)
	std.ID = strategy.ID
	std.Name = strategy.Name
	std.AlarmCount = strategy.AlarmCount
	std.Cycle = strategy.Cycle
	std.Expression = strategy.Expression
	std.Description = strategy.Description
	std.TriggerTemplates = triggers
	return std
}

// GetStrategyTemplate 获取单个策略模板
func (d *db) GetStrategyTemplate(strategyTemplateID int) *StrategyTemplate {
	strategyTemplate := StrategyTemplate{}
	if err := d.Get(&strategyTemplate, "SELECT * FROM strategy_template WHERE id=?", strategyTemplateID); err != nil {
		log.Println(err)
	}
	return &strategyTemplate
}

// GetTriggerTemplates 获取表达式模板
func (d *db) GetTriggerTemplates(strategyTemplateID int) []*TriggerTemplate {
	triggerTemplates := []*TriggerTemplate{}
	if err := d.Select(&triggerTemplates, "SELECT * FROM trigger_template WHERE strategy_template_id=?", strategyTemplateID); err != nil {
		log.Println(err)
	}
	return triggerTemplates
}

// GetStrategyTemplate 获取多个策略模板
func (d *db) GetStrategyTemplates() []*StrategyTemplate {
	strategyTemplates := []*StrategyTemplate{}
	if err := d.Select(&strategyTemplates, "SELECT * FROM strategy_template"); err != nil {
		log.Println(err)
	}
	return strategyTemplates
}

// CreateStrategyTemplate 创建一条策略模板
func (d *db) CreateStrategyTemplate(strategyTemplate *StrategyTemplateDetail) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			log.Println(err)
		}
	}()
	tx := d.MustBegin()
	defer tx.Rollback()
	r := tx.MustExec("INSERT INTO strategy_template VALUES (DEFAULT, ?, ?, ?, ?, ?)",
		strategyTemplate.Name, strategyTemplate.AlarmCount, strategyTemplate.Cycle, strategyTemplate.Expression, strategyTemplate.Description)
	strategyID, err := r.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, tt := range strategyTemplate.TriggerTemplates {
		tx.MustExec("INSERT INTO trigger_template VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			strategyID, tt.Metric, tt.Tags, tt.Number, tt.Index, tt.Method, tt.Symbol, tt.Threshold, tt.Description)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return
}

// UpdateStrategyTemplate 修改策略模板
func (d *db) UpdateStrategyTemplate(strategyTemplate *StrategyTemplateDetail) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			log.Println(err)
		}
	}()
	tx := d.MustBegin()
	defer tx.Rollback()
	tx.MustExec("REPLACE INTO strategy_template VALUES (?, ?, ?, ?, ?, ?)",
		strategyTemplate.ID, strategyTemplate.Name, strategyTemplate.AlarmCount, strategyTemplate.Cycle, strategyTemplate.Expression, strategyTemplate.Description)
	tx.MustExec("DELETE FROM trigger_template WHERE strategy_template_id = ?", strategyTemplate.ID)
	for _, tt := range strategyTemplate.TriggerTemplates {
		tx.MustExec("INSERT INTO trigger_template VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			tt.ID, tt.StrategyTemplateID, tt.Metric, tt.Tags, tt.Number, tt.Index, tt.Method, tt.Symbol, tt.Threshold, tt.Description)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return
}

// DeleteStrategyTemplates 批量删除策略模板
func (d *db) DeleteStrategyTemplates(ids []string) (err error) {
	rawSQL := fmt.Sprintf("DELETE FROM strategy_template WHERE id in (%s)", strings.Join(ids, ","))
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return
}

// ------------------------------------ new --------------------------------------------

func (d *db) getProducts(paging bool, query string, order string, isDelete string, offset, limit int) (int, []*Product) {
	var (
		products = []*Product{}
		err      error
		cnt      int
	)
	rawSQL := fmt.Sprintf("select id, name, description, creator from product")
	cntSQL := fmt.Sprintf("select count(*) from product")
	switch isDelete {
	case "true":
		rawSQL = fmt.Sprintf("%s where is_delete = 1", rawSQL)
		cntSQL = fmt.Sprintf("%s where is_delete = 1", cntSQL)
	default:
		rawSQL = fmt.Sprintf("%s where is_delete = 0", rawSQL)
		cntSQL = fmt.Sprintf("%s where is_delete = 0", cntSQL)
	}
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and name like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&products, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, products
}

//获取用户所在的产品线列表
func (d *db) getUserProducts(user *User) []Product {
	products := []Product{}
	rawSQL := fmt.Sprintf("select p.id, p.name, p.description, p.creator from product p where p.id in "+
		"(select product_id from product_user where p.is_delete = 0 and user_id =%d)", user.ID)
	log.Println(rawSQL)
	if err := d.Select(&products, rawSQL); err != nil {
		log.Println(err)
	}
	return products
}

// 根据名称获取产品线
func (d *db) getProductByName(productName string) *Product {
	product := &Product{}
	rawSQL := fmt.Sprintf("select id, name, description, creator from product where name='%s'", productName)
	log.Println(rawSQL)
	if err := d.Get(product, rawSQL); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
	}
	return product
}

func (d *db) getAllUsers(paging bool, query string, order string, offset, limit int) (int, []User) {
	var (
		users = make([]User, 0)
		cnt   int
		err   error
	)
	rawSQL := fmt.Sprintf("select id, username, display_name, role, wechat, type, mail, phone, status, create_at, update_at from user")
	cntSQL := fmt.Sprintf("select count(*) from user")
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s where username like '%%%s%%' or display_name like '%%%s%%' or phone like '%%%s%%'", rawSQL, query, query, query)
		cntSQL = fmt.Sprintf("%s where username like '%%%s%%' or display_name like '%%%s%%' or phone like '%%%s%%'", cntSQL, query, query, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&users, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, users
}

// get user profile
func (d *db) getUserProfile(username string) *User {
	user := &User{}
	rawSQL := fmt.Sprintf("select id, username, password, display_name, phone, wechat, type, mail, role, "+
		"DATE_FORMAT(create_at,'%s') as create_at, DATE_FORMAT(update_at,'%s') "+
		"as update_at from user where username='%s'",
		timeFormat, timeFormat, username)
	log.Println(rawSQL)
	if err := d.Get(user, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		return nil
	}
	return user
}

// 创建用户
func (d *db) createUser(user *User) (*User, error) {
	rawSQL := fmt.Sprintf("insert into user(username, password, display_name, mail, phone, wechat, type, role, status) "+
		"values('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d)",
		user.Username, user.Password, user.DisplayName, user.EmailAddress, user.PhoneNum, user.Wechat, user.Type, user.Role, user.Status)
	log.Println(rawSQL)
	res, err := d.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	user.ID = int(id)
	return user, nil
}

//根据用户id删除用户
func (d *db) deleteUser(userID int) error {
	rawSQL := fmt.Sprintf("delete from user where id=%d", userID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//更新用户配置信息
func (d *db) updateUserProfile(user *User) error {
	rawSQL := fmt.Sprintf("update user set display_name='%s', password='%s', phone='%s', wechat='%s', update_at='%s' where id=%d",
		user.DisplayName, user.Password, user.PhoneNum, user.Wechat, time.Now().Format(timeFormat), user.ID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//设置用户角色
func (d *db) setUserRole(userID int, roleNum int) error {
	rawSQL := fmt.Sprintf("update user set role=%d where id=%d", roleNum, userID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//创建产品
func (d *db) createProduct(product *Product) (*Product, error) {
	rawSQL := fmt.Sprintf("insert into product(name, description, creator, create_at) values('%s', '%s', '%s', '%s')",
		product.Name, product.Description, product.Creator, time.Now().Format(timeFormat))
	log.Println(rawSQL)
	res, err := d.Exec(rawSQL)
	if err != nil {
		log.Printf("createProduct error %s", err)
		return nil, err
	}

	id, _ := res.LastInsertId()
	product.ID = int(id)
	return product, nil
}

// 更新产品线
func (d *db) updateProduct(product *Product) error {
	rawSQL := fmt.Sprintf("update product set name='%s', description='%s' where id=%d", product.Name, product.Description, product.ID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//软删除产品线
func (d *db) deleteProduct(productID int) error {
	rawSQL := fmt.Sprintf("update product set is_delete=1 where id=%d", productID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ----------------------------------  Product User --------------------------------------------------

// 获取产品下的用户
func (d *db) getProductUsers(productID int, paging bool, query string, order string, offset, limit int) (int, []types.User) {
	var (
		users = make([]types.User, 0)
		err   error
	)
	rawSQL := fmt.Sprintf("select u.id, u.username, u.role, u.phone, u.mail, u.wechat,u.status "+
		"from user u inner join product_user pu on u.id = pu.user_id where pu.product_id = %d", productID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", rawSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	if err = d.Select(&users, rawSQL); err != nil {
		log.Println(err)
	}
	return d.getProductUsersCnt(productID, query), users
}

//获取产品线用户数
func (d *db) getProductUsersCnt(productID int, query string) int {
	var cnt int
	cntSQL := fmt.Sprintf("select count(*) from user u inner join product_user pu on u.id = pu.user_id where pu.product_id = %d", productID)
	if len(query) > 0 {
		cntSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", cntSQL, query)
	}
	log.Println(cntSQL)
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt
}

//获取不在产品线的用户
func (d *db) getNotInProductUsers(productID int, paging bool, query string, order string, offset, limit int) (int, []types.User) {
	var (
		users = make([]types.User, 0)
		cnt   int
		err   error
	)
	rawSQL := fmt.Sprintf("select u.id, u.username, u.role, u.phone, u.mail, u.wechat,u.status "+
		"from user u where u.id not in (select user_id from product_user where product_id = %d)", productID)
	cntSQL := fmt.Sprintf("select count(*) from user u where u.id not in (select user_id from product_user where product_id = %d)", productID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&users, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, users
}

// 添加用户到产品线
func (d *db) addUsers2Product(productID int, userids []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range userids {
		rawSQL := fmt.Sprintf("insert ignore into product_user(product_id, user_id) values(%d, %d)",
			productID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// 从产品线中移除用户
func (d *db) removeUsersFromProduct(productID int, userids []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range userids {
		rawSQL := fmt.Sprintf("delete from product_user where product_id = %d and user_id = %d",
			productID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// --------------------------------------- UserGroup --------------------------------------

// 获取产品线下的用户组
func (d *db) getProductUserGroups(productID int, paging bool, query string, order string, offset, limit int) (int, []UserGroup) {
	var (
		userGroups = make([]UserGroup, 0)
		err        error
		cnt        int
	)
	rawSQL := fmt.Sprintf("select id, name, description from user_group")
	cntSQL := fmt.Sprintf("select count(*) from user_group")
	if productID != 0 {
		rawSQL = fmt.Sprintf("%s where product_id = %d", rawSQL, productID)
		cntSQL = fmt.Sprintf("%s where product_id = %d", cntSQL, productID)
	}
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and name like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&userGroups, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, userGroups
}

//获取指定名称的用户组
func (d *db) findProductUserGroup(productID int, groupName string) *UserGroup {
	userGroup := &UserGroup{}
	rawSQL := fmt.Sprintf("select id, name, description from user_group where product_id = %d and name = '%s'",
		productID, groupName)
	log.Println(rawSQL)
	if err := d.Get(userGroup, rawSQL); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
	}
	return userGroup
}

//创建用户组
func (d *db) createProductUserGroup(productID int, group UserGroup) (UserGroup, error) {
	rawSQL := fmt.Sprintf("insert into user_group(product_id, name, description) values(%d, '%s', '%s')",
		productID, group.Name, group.Desc)
	res, err := d.Exec(rawSQL)
	if err != nil {
		return group, err
	}
	id, _ := res.LastInsertId()
	group.ID = int(id)
	return group, nil
}

//更新用户组
func (d *db) updateProductUserGroup(productID int, group UserGroup) error {
	rawSQL := fmt.Sprintf("update user_group set name = '%s', description='%s' where product_id = %d and id = %d",
		group.Name, group.Desc, productID, group.ID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//删除产品线用户组
func (d *db) deleteProductUserGroup(productID, userGroupID int) error {
	rawSQL := fmt.Sprintf("delete from user_group where id =%d and product_id = %d", userGroupID, productID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//获取用户组下的用户
func (d *db) getProductUserGroupUsers(userGroupID int, paging bool, query string, order string, offset, limit int) (int, []types.User) {
	var (
		users = make([]types.User, 0)
		err   error
		cnt   int
	)
	rawSQL := fmt.Sprintf("select u.id, u.username, u.role, u.phone, u.mail, u.wechat,u.status "+
		"from user u where id in (select user_id from user_group_user where user_group_id = %d)", userGroupID)
	cntSQL := fmt.Sprintf("select count(*) from user u where id in"+
		" (select user_id from user_group_user where user_group_id = %d)", userGroupID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&users, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, users
}

//获取不在产品线用户组内的用户
func (d *db) getNotInProductUserGroupUsers(productID, userGroupID int, paging bool, query string, order string, offset, limit int) (int, []types.User) {
	var (
		users = make([]types.User, 0)
		err   error
		cnt   int
	)
	rawSQL := fmt.Sprintf("select u.id, u.username, u.role, u.phone, u.mail, u.wechat,u.status "+
		"from user u inner join product_user as pu on u.id = pu.user_id where pu.product_id = %d and u.id not "+
		"in (select user_id from user_group_user where user_group_id = %d)", productID, userGroupID)
	cntSQL := fmt.Sprintf("select count(*) from user u inner join product_user as pu "+
		" on u.id = pu.user_id where pu.product_id = %d and u.id not "+
		"in (select user_id from user_group_user where user_group_id = %d)", productID, userGroupID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and (u.username like '%%%s%%')", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&users, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, users
}

//添加用户到用户组
func (d *db) addUsers2UserGroup(groupID int, userids []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range userids {
		rawSQL := fmt.Sprintf("insert into user_group_user(user_group_id, user_id) values(%d, %d)",
			groupID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//从用户组移除用户
func (d *db) removeUsersFromUserGroup(groupID int, userids []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range userids {
		rawSQL := fmt.Sprintf("delete from user_group_user where user_group_id =%d and user_id = %d",
			groupID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//获取所有主机
func (d *db) getAllHosts(paging bool, noProduct bool, query string, order string, offset, limit int) (int, []WarpHost) {
	var (
		hosts = make([]WarpHost, 0)
		cnt   int
		err   error
	)
	rawSQL := fmt.Sprintf("select host.id, host.ip, host.name, host.hostname, host.agent_version, host.status, "+
		"DATE_FORMAT(host.create_at,'%s') as create_at, DATE_FORMAT(host.update_at,'%s') as update_at,"+
		"DATE_FORMAT(host.mute_time, '%s') as mute_time, host.uptime, host.idle_pct, count(host_plugin.id) as plugin_cnt, "+
		"IFNULL(product.name,'-') as products from host left join host_plugin on host.id = host_plugin.host_id left join "+
		"(select group_concat(p.name) as name, ph.host_id from product_host as ph left join product as p on ph.product_id = p.id "+
		"group by ph.host_id) as product on host.id = product.host_id ",
		dbDateFormat, dbDateFormat, dbDateFormat)
	cntSQL := fmt.Sprintf("select count(host.id) from host left join " +
		"(select group_concat(p.name) as name, ph.host_id from product_host as ph left join product as p on ph.product_id = p.id " +
		"group by ph.host_id) as product on host.id = product.host_id ")

	if noProduct {
		rawSQL = fmt.Sprintf("%s where product.name is null ", rawSQL)
		cntSQL = fmt.Sprintf("%s where product.name is null ", cntSQL)
	}

	if len(query) > 0 {
		if noProduct {
			rawSQL = fmt.Sprintf("%s and (ip like '%%%s%%' or hostname like '%%%s%%')", rawSQL, query, query)
			cntSQL = fmt.Sprintf("%s and (ip like '%%%s%%' or hostname like '%%%s%%')", cntSQL, query, query)
		} else {
			rawSQL = fmt.Sprintf("%s where ip like '%%%s%%' or hostname like '%%%s%%'", rawSQL, query, query)
			cntSQL = fmt.Sprintf("%s where ip like '%%%s%%' or hostname like '%%%s%%'", cntSQL, query, query)
		}
	}
	rawSQL = fmt.Sprintf("%s group by host.id", rawSQL)
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&hosts, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hosts
}

func (d *db) getHostByID(hostID string) *Host {
	host := &Host{}
	rawSQL := fmt.Sprintf("select id, ip, name, hostname, agent_version, status, "+
		"DATE_FORMAT(create_at,'%s') as create_at, DATE_FORMAT(update_at,'%s') as update_at,"+
		"DATE_FORMAT(mute_time, '%s') as mute_time, uptime, idle_pct from host where id='%s'",
		dbDateFormat, dbDateFormat, dbDateFormat, hostID)
	log.Println(rawSQL)
	if err := d.Get(host, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}
	return host
}

func (d *db) getHostByIP(hostIP string) *Host {
	host := &Host{}
	rawSQL := fmt.Sprintf("select id, ip, name, hostname, agent_version, status, "+
		"DATE_FORMAT(create_at,'%s') as create_at, DATE_FORMAT(update_at,'%s') as update_at,"+
		"DATE_FORMAT(mute_time, '%s') as mute_time, uptime, idle_pct from host where ip='%s'",
		dbDateFormat, dbDateFormat, dbDateFormat, hostIP)
	log.Println(rawSQL)
	if err := d.Get(host, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}
	return host
}

//删除主机
func (d *db) deleteHost(hostID string) error {
	rawSQL := fmt.Sprintf("delete from host where id='%s'", hostID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *db) muteHost(hostID string, muteTime string) error {
	rawSQL := fmt.Sprintf("update host set mute_time='%s' where id='%s'", muteTime, hostID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//获取主机metrics

func (d *db) getHostMetrics(hostID string, paging bool, prefix string, query string, order string, offset, limit int) (int, []*MetricSummary) {
	metrics := make([]*MetricSummary, 0)
	cnt := 0
	rawSQL := fmt.Sprintf("select id, metric, tags, dt, cycle, DATE_FORMAT(update_at,'%s') as update_at "+
		" from metric where host_id='%s'", dbDateFormat, hostID)
	cntSQL := fmt.Sprintf("select count(*) from metric where host_id='%s'", hostID)

	if len(prefix) != 0 {
		rawSQL = fmt.Sprintf("%s and metric like '%s%%'", rawSQL, prefix)
		cntSQL = fmt.Sprintf("%s and metric like '%s%%'", cntSQL, prefix)
	}
	if len(query) != 0 {
		queryArr := strings.Split(query, "/")
		if len(queryArr) > 1 {
			rawSQL = fmt.Sprintf("%s and (metric like '%%%s%%' and tags like '%%%s%%')", rawSQL, queryArr[0], queryArr[1])
			cntSQL = fmt.Sprintf("%s and (metric like '%%%s%%' and tags like '%%%s%%')", cntSQL, queryArr[0], queryArr[1])
		} else {
			rawSQL = fmt.Sprintf("%s and metric like '%%%s%%' ", rawSQL, queryArr[0])
			cntSQL = fmt.Sprintf("%s and metric like '%%%s%%' ", cntSQL, queryArr[0])
		}
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&metrics, rawSQL); err != nil {
		log.Println(err)
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, metrics
}

func (d *db) deleteHostMetrics(hostID string, ids []int) error {

	tx := d.MustBegin()
	for _, id := range ids {
		rawSQL := fmt.Sprintf("delete from metric where host_id ='%s' and id = %d", hostID, id)
		log.Println(rawSQL)
		if _, err := tx.Exec(rawSQL); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (d *db) getHostAppNames(hostID string) []string {
	metrics := []string{}
	appNames := []string{}
	appMap := make(map[string]struct{})
	rawSQL := fmt.Sprintf("select distinct metric from metric where host_id='%s'", hostID)
	log.Println(rawSQL)
	if err := d.Select(&metrics, rawSQL); err != nil {
		log.Println(err)
		return nil
	}
	for _, metric := range metrics {
		appName := strings.Split(metric, ".")[0]
		if _, ok := appMap[appName]; ok {
			continue
		}
		appMap[appName] = struct{}{}
		appNames = append(appNames, appName)
	}
	sort.Strings(appNames)
	return appNames
}

func (d *db) getHostHostGroups(productID int, hostID string) []HostGroup {
	var (
		hostGroups = make([]HostGroup, 0)
		err        error
		rawSQL     string
	)
	rawSQL = fmt.Sprintf("select id, name, description, creator, DATE_FORMAT(create_at,'%s') as create_at,"+
		"DATE_FORMAT(update_at,'%s') as update_at  from host_group where id in (select host_group_id from host_group_host where host_id='%s')",
		dbDateFormat, dbDateFormat, hostID)
	if productID != 0 {
		rawSQL = fmt.Sprintf("select id, name, description, creator, DATE_FORMAT(create_at,'%s') as create_at,"+
			"DATE_FORMAT(update_at,'%s') as update_at  from host_group where product_id = %d and id in (select host_group_id from host_group_host where host_id='%s')",
			dbDateFormat, dbDateFormat, productID, hostID)
	}
	log.Println(rawSQL)
	if err = d.Select(&hostGroups, rawSQL); err != nil {
		log.Println(err)
	}
	return hostGroups
}

func (d *db) getHostProducts(hostID string) []Product {
	var (
		products = make([]Product, 0)
		err      error
	)
	rawSQL := fmt.Sprintf("select id, name, description, creator from product where id in (select product_id from product_host where host_id='%s')", hostID)
	log.Print(rawSQL)
	if err = d.Select(&products, rawSQL); err != nil {
		log.Println(err)
	}
	return products
}

func (d *db) getHostPlugins(hostID string, paging bool, query string, order string, offset, limit int) (int, []*types.Plugin) {
	plugins := make([]*types.Plugin, 0)
	cnt := 0
	rawSQL := fmt.Sprintf("select hp.id, p.name, p.path, p.checksum, hp.args, hp.interval, hp.timeout, hp.comment from host_plugin as hp"+
		" left join plugin as p on p.id = hp.plugin_id where hp.host_id='%s'", hostID)
	cntSQL := fmt.Sprintf("select count(*) from host_plugin as hp left join plugin as p on p.id = hp.plugin_id where host_id='%s'", hostID)

	if len(query) != 0 {
		rawSQL = fmt.Sprintf("%s and (p.name like '%%%s%%' or p.path like '%%%s%%')", rawSQL, query, query)
		cntSQL = fmt.Sprintf("%s and (p.name like '%%%s%%' and p.path like '%%%s%%')", cntSQL, query, query)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&plugins, rawSQL); err != nil {
		log.Println(err)
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, plugins
}

func (d *db) createHostPlugin(hostID string, plugin *types.Plugin) (*types.Plugin, error) {
	rawSQL := fmt.Sprintf("insert into host_plugin(`host_id`, `plugin_id`, `args`, `interval`, `timeout`, `comment`)"+
		" values('%s', %d, '%s', %d, %d, '%s')",
		hostID, plugin.ID, plugin.Args, plugin.Interval, plugin.Timeout, plugin.Comment)
	log.Println(rawSQL)
	res, err := d.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	plugin.ID = int(id)
	return plugin, nil
}

func (d *db) updateHostPlugin(hostID string, plugin *types.Plugin) error {
	rawSQL := fmt.Sprintf("update host_plugin set `args`='%s', `interval`=%d, `timeout`=%d, `comment`='%s' "+
		" where host_id='%s' and id=%d",
		plugin.Args, plugin.Interval, plugin.Timeout, plugin.Comment, hostID, plugin.ID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *db) deleteHostPlugin(hostID string, pluginID int) error {
	rawSQL := fmt.Sprintf("delete from host_plugin where host_id='%s' and id=%d", hostID, pluginID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 获取产品下的主机
func (d *db) getProductHosts(productID int, noGroup bool, paging bool, query string, order string, offset, limit int) (int, []WarpHost) {
	var (
		hosts = make([]WarpHost, 0)
		err   error
	)

	rawSQL := fmt.Sprintf("select host.id, host.ip, host.name, host.hostname, host.agent_version, host.status,"+
		"DATE_FORMAT(host.create_at,'%s') as create_at, DATE_FORMAT(host.update_at,'%s') as update_at,"+
		"host.mute_time, host.uptime, host.idle_pct, count(host_plugin.id) as plugin_cnt, "+
		"IFNULL(groups.name,'-') as groups from host right join product_host on product_host.host_id = host.id "+
		"left join host_plugin on host.id = host_plugin.host_id left join (select group_concat(hg.name) as name,"+
		"hgh.host_id from host_group_host as hgh left join host_group as hg on hgh.host_group_id = hg.id where hg.product_id=%d "+
		"group by hgh.host_id) as groups on host.id = groups.host_id where product_host.product_id=%d",
		dbDateFormat, dbDateFormat, productID, productID)
	if noGroup {
		rawSQL = fmt.Sprintf("%s and groups.name is null", rawSQL)
	}
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (ip like '%%%s%%' or hostname like '%%%s%%')", rawSQL, query, query)
	}
	rawSQL = fmt.Sprintf("%s group by host.id", rawSQL)
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d, %d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	if err = d.Select(&hosts, rawSQL); err != nil {
		log.Println(err)
	}
	return d.getProductHostsCnt(productID, noGroup, query), hosts
}

//获取产品线主机数
func (d *db) getProductHostsCnt(productID int, noGroup bool, query string) int {
	var cnt int
	cntSQL := fmt.Sprintf("select count(*) from host h inner join product_host ph on h.id = ph.host_id where ph.product_id = %d", productID)
	if noGroup {
		cntSQL = fmt.Sprintf("%s and h.id not in (select host_id from host_group_host where host_group_id in (select id from host_group where product_id= %d))", cntSQL, productID)
	}
	if len(query) > 0 {
		cntSQL = fmt.Sprintf("%s and (h.ip like '%%%s%%' or h.hostname like '%%%s%%')", cntSQL, query, query)
	}
	log.Println(cntSQL)
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt
}

//获取不在产品线的主机
func (d *db) getNotInProductHosts(productID int, paging bool, query string, order string, offset, limit int) (int, []Host) {
	hosts := []Host{}
	cnt := 0
	rawSQL := fmt.Sprintf("select h.id, h.ip, h.name, h.hostname, h.agent_version, h.status, "+
		"DATE_FORMAT(h.create_at,'%s') as create_at, DATE_FORMAT(h.update_at,'%s') as update_at,"+
		" h.mute_time, h.uptime, h.idle_pct from host as h where h.id not in (select host_id from product_host where product_id = %d)",
		dbDateFormat, dbDateFormat, productID)
	cntSQL := fmt.Sprintf("select count(*) from host as h where h.id not in(select host_id from product_host where product_id = %d)", productID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (h.ip like '%%%s%%' or h.hostname like '%%%s%%')",
			rawSQL, query, query)
		cntSQL = fmt.Sprintf("%s and (h.ip like '%%%s%%' or h.hostname like '%%%s%%')",
			cntSQL, query, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&hosts, rawSQL); err != nil {
		log.Println(err)
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hosts
}

//添加主机到产品线
func (d *db) addHosts2Product(productID int, ids []string) error {
	tx := d.MustBegin()
	var err error
	for _, id := range ids {
		rawSQL := fmt.Sprintf("insert ignore into product_host(product_id, host_id) values(%d,'%s')", productID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//从产品线移除主机
func (d *db) removeHostsFromProduct(productID int, ids []string) (err error) {
	tx := d.MustBegin()
	for idx, id := range ids {
		ids[idx] = fmt.Sprintf("'%s'", id)
	}
	hostIDS := strings.Join(ids, ",")
	rawSQL := fmt.Sprintf("delete from product_host where product_id = %d and host_id in (%s)", productID, hostIDS)
	log.Println(rawSQL)
	if _, err = tx.Exec(rawSQL); err != nil {
		tx.Rollback()
		return
	}
	rawSQL = fmt.Sprintf("delete from host_group_host where host_id in (%s) "+
		"and host_group_id in (select id from host_group where product_id=%d)", hostIDS, productID)
	log.Println(rawSQL)
	if _, err = tx.Exec(rawSQL); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

//获取产品线下的主机组
func (d *db) getProductHostGroups(productID int, paging bool, query string, user string, order string, offset, limit int) (int, []WarpHostGroup) {
	var (
		groups = make([]WarpHostGroup, 0)
		err    error
		cnt    int
		rawSQL string
	)
	rawSQL = fmt.Sprintf("select hg.id, hg.name, hg.description, hg.creator, DATE_FORMAT(hg.create_at,'%s') as create_at,"+
		"DATE_FORMAT(hg.update_at,'%s') as update_at, count(distinct host_group_plugin.id) as plugin_cnt, "+
		"count(distinct host_group_host.id) as host_cnt, count(distinct strategy_group.id) as strategy_cnt "+
		" from host_group as hg left join host_group_plugin on hg.id = host_group_plugin.group_id left join host_group_host "+
		" on hg.id = host_group_host.host_group_id left join strategy_group on hg.id=strategy_group.group_id "+
		" where hg.product_id=%d",
		dbDateFormat, dbDateFormat, productID)
	cntSQL := fmt.Sprintf("select count(*) from host_group where product_id = %d", productID)
	if len(user) > 0 {
		rawSQL = fmt.Sprintf("%s and hg.creator='%s'", rawSQL, user)
		cntSQL = fmt.Sprintf("%s and creator='%s'", cntSQL, user)
	}
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and hg.name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and name like '%%%s%%'", cntSQL, query)
	}
	rawSQL = fmt.Sprintf("%s group by hg.id", rawSQL)
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&groups, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, groups
}

func (d *db) getHostGroupPlugins(groupID int, paging bool, query string, offset, limit int) (int, []types.Plugin) {
	var (
		plugins = make([]types.Plugin, 0)
		err     error
		cnt     int
	)
	rawSQL := fmt.Sprintf("select hgp.id, p.name, p.path, p.checksum, hgp.args, hgp.interval, hgp.timeout, hgp.comment from host_group_plugin as hgp "+
		" left join plugin as p on p.id = hgp.plugin_id where hgp.group_id=%d", groupID)
	cntSQL := fmt.Sprintf("select count(*) from host_group_plugin as hgp left join plugin as p on p.id = hgp.plugin_id where hgp.group_id=%d", groupID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and p.name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and p.name like '%%%s%%'", cntSQL, query)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&plugins, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, plugins
}

func (d *db) updateHostGroupPlugin(groupID int, plugin *types.Plugin) error {
	rawSQL := fmt.Sprintf("update host_group_plugin set `args`='%s', `interval`=%d, `timeout`=%d, `comment`='%s' "+
		" where group_id=%d and id=%d",
		plugin.Args, plugin.Interval, plugin.Timeout, plugin.Comment, groupID, plugin.ID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *db) createHostGroupPlugin(groupID int, plugin *types.Plugin) (*types.Plugin, error) {
	rawSQL := fmt.Sprintf("insert into host_group_plugin(`group_id`, `plugin_id`, `args`, `interval`, `timeout`, `comment`)"+
		" values(%d, %d, '%s', %d, %d, '%s')",
		groupID, plugin.ID, plugin.Args, plugin.Interval, plugin.Timeout, plugin.Comment)
	log.Println(rawSQL)
	res, err := d.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	plugin.ID = int(id)
	return plugin, nil
}

func (d *db) deleteHostGroupPlugin(groupID int, pluginID int) error {
	rawSQL := fmt.Sprintf("delete from host_group_plugin where group_id=%d and id=%d", groupID, pluginID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//获取主机组下的主机
func (d *db) getProductHostGroupHosts(productID, groupID int, paging bool, query string, order string, offset, limit int) (int, []WarpHost) {
	var (
		hosts = make([]WarpHost, 0)
		err   error
		cnt   int
	)
	rawSQL := fmt.Sprintf("select host.id, host.ip, host.name, host.hostname, host.agent_version, host.status,"+
		"DATE_FORMAT(host.create_at,'%s') as create_at, DATE_FORMAT(host.update_at,'%s') as update_at,"+
		"host.mute_time, host.uptime, host.idle_pct, count(host_plugin.id) as plugin_cnt, "+
		"IFNULL(groups.name,'-') as groups from host left join host_plugin on host.id = host_plugin.host_id left join "+
		"(select group_concat(hg.name) as name, hgh.host_id from host_group_host as hgh left join host_group as hg"+
		" on hgh.host_group_id = hg.id where hg.product_id=%d group by hgh.host_id) as groups on host.id = groups.host_id "+
		" left join host_group_host as hgh on host.id = hgh.host_id where hgh.host_group_id=%d",
		dbDateFormat, dbDateFormat, productID, groupID)

	cntSQL := fmt.Sprintf("select count(*) from host where id in (select host_id from host_group_host where host_group_id = %d)",
		groupID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (host.ip like '%%%s%%' or host.hostname like '%%%s%%')", rawSQL, query, query)
		cntSQL = fmt.Sprintf("%s and (host.ip like '%%%s%%' or host.hostname like '%%%s%%')", cntSQL, query, query)
	}
	rawSQL = fmt.Sprintf("%s group by host.id", rawSQL)
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&hosts, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hosts
}

//获取不在产品线主机组内的主机
func (d *db) getNotInProductHostGroupHosts(productID int, groupID int, paging bool, query string, order string, offset, limit int) (int, []Host) {
	var (
		hosts = make([]Host, 0)
		err   error
		cnt   int
	)
	rawSQL := fmt.Sprintf("select h.id, h.ip, h.name, h.hostname, h.agent_version,"+
		" h.status, DATE_FORMAT(h.create_at,'%s') as create_at, DATE_FORMAT(h.update_at,'%s') as update_at,"+
		" h.mute_time, h.uptime, h.idle_pct from host h inner join product_host as ph on h.id = ph.host_id "+
		"where ph.product_id = %d and h.id not in"+
		"(select host_id from host_group_host where host_group_id= %d)",
		dbDateFormat, dbDateFormat, productID, groupID)
	cntSQL := fmt.Sprintf("select count(*) from host h inner join product_host as ph on h.id = ph.host_id "+
		"where ph.product_id = %d and h.id not in"+
		"(select host_id from host_group_host where host_group_id= %d)",
		productID, groupID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and (h.ip like '%%%s%%' or h.hostname like '%%%s%%')", rawSQL, query, query)
		cntSQL = fmt.Sprintf("%s and (h.ip like '%%%s%%' or h.hostname like '%%%s%%')", cntSQL, query, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err = d.Select(&hosts, rawSQL); err != nil {
		log.Println(err)
	}
	if err = d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hosts
}

//添加主机到主机组
func (d *db) addHost2HostGroup(groupID int, ids []string) error {
	tx := d.MustBegin()
	var err error
	for _, id := range ids {
		rawSQL := fmt.Sprintf("insert ignore into host_group_host(host_id, host_group_id) values('%s', %d)", id, groupID)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//从主机组内移除主机
func (d *db) removeHostFromHostGroup(groupID int, ids []string) error {
	tx := d.MustBegin()
	var err error
	for _, id := range ids {
		rawSQL := fmt.Sprintf("delete from host_group_host where host_id = '%s' and  host_group_id =%d", id, groupID)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//根据组id获取产品线主机组
func (d *db) getProductHostGroupByID(productID int, groupID int) HostGroup {
	rawSQL := fmt.Sprintf("select id, name, description, creator, DATE_FORMAT(create_at,'%s') as create_at,"+
		"DATE_FORMAT(update_at,'%s') as update_at from host_group where product_id =%d and id = %d",
		dbDateFormat, dbDateFormat, productID, groupID)
	hostGroup := HostGroup{}
	d.Get(&hostGroup, rawSQL)
	return hostGroup
}

//获取产品线下指定名称的主机组
func (d *db) getProductHostGroupByName(productID int, groupName string) HostGroup {
	rawSQL := fmt.Sprintf("select id, name, description, creator, DATE_FORMAT(create_at,'%s') as create_at,"+
		"DATE_FORMAT(update_at,'%s') as update_at from host_group where product_id =%d and name = '%s'",
		dbDateFormat, dbDateFormat, productID, groupName)
	log.Println(rawSQL)
	hostGroup := HostGroup{}
	d.Get(&hostGroup, rawSQL)
	return hostGroup
}

//创建产品线机组
func (d *db) createProductHostGroup(productID int, group HostGroup) error {
	timeString := time.Now().Format(timeFormat)
	rawSQL := fmt.Sprintf("insert into host_group(name, description, creator, product_id, create_at, update_at) "+
		"values('%s', '%s', '%s', %d, '%s', '%s')",
		group.Name, group.Description, group.Creator, productID, timeString, timeString)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

//更新主机组
func (d *db) updateHostGroup(hostGroup HostGroup) error {
	rawSQL := fmt.Sprintf("update host_group set name ='%s', description='%s', update_at='%s' where id = %d",
		hostGroup.Name, hostGroup.Description, time.Now().Format(timeFormat), hostGroup.ID)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

//删除产品线主机组
func (d *db) deleteProductHostGroup(productID int, hostGroupID int) error {
	rawSQL := fmt.Sprintf("delete from host_group where id =%d and product_id = %d", hostGroupID, productID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//获取插件列表
func (d *db) getPlugins(query string, paging bool, offset, limit int) (int, []types.Plugin) {
	plugins := []types.Plugin{}
	cnt := 0
	rawSQL := fmt.Sprintf("select `id`, `name`, `path`, `args`, `interval`, `timeout`, `checksum`, `comment` from plugin")
	cntSQL := fmt.Sprintf("select count(*) from plugin")
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s where name like '%%%s%%' or path like '%%%s%%'", rawSQL, query, query)
		cntSQL = fmt.Sprintf("%s where name like '%%%s%%' or path like '%%%s%%'", cntSQL, query, query)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&plugins, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, plugins
}

//创建插件
func (d *db) createPlugin(plugin *types.Plugin, creator string) (*types.Plugin, error) {
	now := time.Now().Format(timeFormat)
	rawSQL := fmt.Sprintf("insert into plugin(`name`, `args`, `path`, `checksum`,"+
		" `interval`, `create_at`,`update_at`,`creator`, `comment`) "+
		"values('%s', '%s', '%s', '%s', %d, '%s', '%s', '%s', '%s')",
		plugin.Name, plugin.Args, plugin.Path, plugin.Checksum, plugin.Interval, now, now, creator, plugin.Comment)
	log.Println(rawSQL)
	res, err := d.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	plugin.ID = int(id)
	return plugin, nil
}

//更新插件
func (d *db) updatePlugin(plugin types.Plugin) error {
	rawSQL := fmt.Sprintf("update plugin set `name` ='%s', `args`='%s', `path`='%s', `checksum`='%s', `interval`=%d "+
		", `comment`='%s' where id = %d",
		plugin.Name, plugin.Args, plugin.Path, plugin.Checksum, plugin.Interval, plugin.Comment, plugin.ID)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

//删除插件
func (d *db) deletePlugin(pluginID int) error {
	rawSQL := fmt.Sprintf("delete from plugin where id = %d", pluginID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//增加主机组和插件关联关系
func (d *db) addHostGroups2Plugin(pluginID int, hostGroupIds []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range hostGroupIds {
		rawSQL := fmt.Sprintf("insert into host_group_plugin(plugin_id, group_id) values(%d, %d)",
			pluginID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//移除插件和主机组的关联
func (d *db) removeHostGroupsFromPlugin(pluginID int, hostGroupIds []int) error {
	tx := d.MustBegin()
	var err error
	for _, id := range hostGroupIds {
		rawSQL := fmt.Sprintf("delete from host_group_plugin where plugin_id=%d and group_id=%d",
			pluginID, id)
		log.Println(rawSQL)
		if _, err = tx.Exec(rawSQL); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//获取插件关联的主机组
func (d *db) getPluginHostGroups(pluginID int, paging bool, query string, order string, offset, limit int) (int, []PluginHostGroup) {
	hostGroups := []PluginHostGroup{}
	cnt := 0
	rawSQL := fmt.Sprintf("select p.name as product, g.id, g.name, g.description, g.creator, "+
		"DATE_FORMAT(g.create_at,'%s') as create_at, DATE_FORMAT(g.update_at,'%s') as update_at "+
		"from host_group as g inner join product as p on p.id = g.product_id where g.id in "+
		"(select group_id from host_group_plugin where plugin_id = %d)", timeFormat, timeFormat, pluginID)
	cntSQL := fmt.Sprintf("select count(*) from host_group where id in "+
		"(select group_id from host_group_plugin where plugin_id = %d)", pluginID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and g.name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and g.name like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	if err := d.Select(&hostGroups, rawSQL); err != nil {
		log.Println(err)
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hostGroups
}

//获取未关联插件的主机组
func (d *db) getNotInPluginHostGroups(pluginID int, paging bool, query string, order string, offset, limit int) (int, []PluginHostGroup) {
	hostGroups := []PluginHostGroup{}
	cnt := 0
	rawSQL := fmt.Sprintf("select p.name as product, g.id, g.name, g.description, g.creator, "+
		" DATE_FORMAT(g.create_at,'%s') as create_at, DATE_FORMAT(g.update_at,'%s') as update_at "+
		" from host_group as g inner join product as p on p.id = g.product_id where g.id not in "+
		"(select group_id from host_group_plugin where plugin_id = %d)", timeFormat, timeFormat, pluginID)
	cntSQL := fmt.Sprintf("select count(*) from host_group where id not in "+
		"(select group_id from host_group_plugin where plugin_id = %d)", pluginID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and g.name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and g.name like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	if err := d.Select(&hostGroups, rawSQL); err != nil {
		log.Println(err)
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, hostGroups
}

//获取metric列表
func (d *db) suggestMetrics(productID int) []string {
	rawSQL := fmt.Sprintf("select distinct metric from metric")
	// if productID != 0 {
	// rawSQL = fmt.Sprintf("%s where host_id in (select host_id from product_host where product_id = %d)",
	// rawSQL, productID)
	// }
	log.Println(rawSQL)
	metrics := []string{}
	if err := d.Select(&metrics, rawSQL); err != nil {
		log.Println(err)
	}
	return metrics
}

//获取metric对应的tag列表
func (d *db) suggestMetricTagSet(productID int, metric string) map[string][]string {
	rawSQL := fmt.Sprintf("select tags from metric where metric='%s'", metric)
	hostGroupSQL := "select name from host_group"
	if productID != 0 {
		rawSQL = fmt.Sprintf("%s and host_id in (select host_id from product_host where product_id=%d)", rawSQL, productID)
		hostGroupSQL = fmt.Sprintf("%s where product_id = %d", hostGroupSQL, productID)
	}
	log.Println(rawSQL)
	log.Println(hostGroupSQL)
	tagSet := make(map[string][]string)
	tagSetString := []string{}
	groupnameSet := []string{}
	m := map[string]struct{}{}
	if err := d.Select(&tagSetString, rawSQL); err != nil {
		log.Println(err)
		return tagSet
	}
	if err := d.Select(&groupnameSet, hostGroupSQL); err != nil {
		log.Println(err)
		return tagSet
	}
	for _, set := range tagSetString {
		for _, tagkv := range strings.Split(set, ",") {
			if _, ok := m[tagkv]; ok {
				continue
			}
			fields := strings.Split(tagkv, "=")
			if len(fields) != 2 {
				continue
			}
			tagk, tagv := fields[0], fields[1]
			m[tagkv] = struct{}{}
			tagSet[tagk] = append(tagSet[tagk], tagv)
		}
	}
	tagSet["host_group"] = groupnameSet
	return tagSet
}

//获取产品线下看板
func (d *db) getProductPanels(productID int, query string, order string, paging bool, offset, limit int) (int, []Panel) {
	panels := []Panel{}
	cnt := 0
	rawSQL := fmt.Sprintf("select `id`, `name`, `creator` from panel where product_id=%d", productID)
	cntSQL := fmt.Sprintf("select count(*) from panel where product_id=%d", productID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and name like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and name like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&panels, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	return cnt, panels
}

//创建产品线看板
func (d *db) createProductPanel(productID int, panel *Panel) (*Panel, error) {
	rawSQL := fmt.Sprintf("insert into panel(`name`, `product_id`, `creator`)"+
		"values('%s', '%d', '%s')", panel.Name, productID, panel.Creator)
	res, err := d.Exec(rawSQL)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	panel.ID = int(id)
	return panel, nil
}

//更新看板
func (d *db) updatePanel(panel Panel) error {
	rawSQL := fmt.Sprintf("update panel set `name` ='%s' where id = %d", panel.Name, panel.ID)
	log.Println(rawSQL)
	_, err := d.Exec(rawSQL)
	return err
}

//删除看板
func (d *db) deletePanel(panelID int) error {
	rawSQL := fmt.Sprintf("delete from panel where id = %d", panelID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 获取看板下的图表
func (d *db) getPanelCharts(panelID int, query string, order string, paging bool, offset, limit int) (int, []*Chart) {
	charts := []*Chart{}
	cnt := 0
	rawSQL := fmt.Sprintf("select `id`, `title`, `creator`, `span`, `height`, `type` from chart where panel_id=%d", panelID)
	cntSQL := fmt.Sprintf("select count(*) from chart where panel_id=%d", panelID)
	if len(query) > 0 {
		rawSQL = fmt.Sprintf("%s and title like '%%%s%%'", rawSQL, query)
		cntSQL = fmt.Sprintf("%s and title like '%%%s%%'", cntSQL, query)
	}
	if len(order) > 0 {
		rawSQL = fmt.Sprintf("%s order by %s", rawSQL, order)
	}
	if paging {
		rawSQL = fmt.Sprintf("%s limit %d,%d", rawSQL, offset, limit)
	}
	log.Println(rawSQL)
	log.Println(cntSQL)
	if err := d.Select(&charts, rawSQL); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}
	if err := d.Get(&cnt, cntSQL); err != nil {
		log.Println(err)
	}
	for _, chart := range charts {
		chart.Elements = d.getChartElements(chart.ID)
	}
	return cnt, charts
}

//获取图表元素信息
func (d *db) getChartElements(chartID int) []*ChartElement {
	rawSQL := fmt.Sprintf("select metric, tags from chart_element where chart_id = %d", chartID)
	log.Println(rawSQL)
	chartElements := []*ChartElement{}
	if err := d.Select(&chartElements, rawSQL); err != nil {
		log.Println(err)
	}
	return chartElements
}

//删除图表元素
func (d *db) deleteChartElements(chartID int) error {
	rawSQL := fmt.Sprintf("delete from chart_element where chart_id=%d", chartID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		return err
	}
	return nil
}

//创建看板图表
func (d *db) createPanelChart(panelID int, chart *Chart) (*Chart, error) {
	rawSQL := fmt.Sprintf("insert into chart(`title`, `span`, `height`, `creator`, `create_at`, `type`, `panel_id`)"+
		"values('%s', %d, %d, '%s', '%s', %d, %d)",
		chart.Title, chart.Span, chart.Height, chart.Creator, time.Now().Format(timeFormat), chart.Type, panelID)
	log.Println(rawSQL)
	tx := d.MustBegin()
	res, err := tx.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}
	id, _ := res.LastInsertId()
	chart.ID = int(id)
	for _, ele := range chart.Elements {
		rawSQL = fmt.Sprintf("insert into chart_element(`metric`, `tags`, `chart_id`) values('%s', '%s', %d)",
			ele.Metric, ele.Tags, chart.ID)
		log.Println(rawSQL)
		_, err = tx.Exec(rawSQL)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return chart, nil
}

//更新看板图表
func (d *db) updatePanelChart(panelID int, chart *Chart) error {
	rawSQL := fmt.Sprintf("update chart set `title` = '%s', `span` = %d, `height` = %d, `type` = %d where id = %d and panel_id = %d",
		chart.Title, chart.Span, chart.Height, chart.Type, chart.ID, panelID)
	log.Println(rawSQL)
	tx := d.MustBegin()
	_, err := tx.Exec(rawSQL)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	d.deleteChartElements(chart.ID)
	for _, ele := range chart.Elements {
		rawSQL = fmt.Sprintf("insert into chart_element(`metric`, `tags`, `chart_id`) values('%s', '%s', %d)",
			ele.Metric, ele.Tags, chart.ID)
		log.Println(rawSQL)
		_, err = tx.Exec(rawSQL)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

//删除看板图表
func (d *db) deletePanelChart(panelID, chartID int) error {
	rawSQL := fmt.Sprintf("delete from chart where id = %d and panel_id = %d", chartID, panelID)
	log.Println(rawSQL)
	if _, err := d.Exec(rawSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

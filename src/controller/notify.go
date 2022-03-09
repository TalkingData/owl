package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"owl/common/types"
	"text/template"
	"time"
)

var STRATEGY_STATUS_MAPPING = map[int]string{1: "活跃报警", 2: "已知悉报警", 3: "已关闭报警"}
var STRATEGY_PRIORITY_MAPPING = map[int]string{1: "严重", 2: "较严重", 3: "注意"}

type templateStrategy struct {
	ID                int64
	NAME              string
	TYPE              string
	CYCLE             int
	PRIORITY          string
	STATUS            string
	ALARM_COUNT       int
	COUNT             int
	UPDATE_TIME       string
	EXPRESSION        string
	EXPRESSION_DETAIL string
}

type templateHost struct {
	NAME   string
	IP     string
	STATUS string
}

type notifyTemplate struct {
	STRATEGY templateStrategy
	HOST     templateHost
}

func (c *Controller) processStrategyEventForever() {
	c.eventQueuesMutex.RLock()
	defer c.eventQueuesMutex.RUnlock()
	for _, queue := range c.eventQueues {
		lg.Info("process queue event %s", queue.name)
		go processSingleQueue(queue)
	}
}

//TODO: fix when product delete, goroutine leak
func processSingleQueue(queue *EventPool) {
	duration := time.Millisecond * 100
	for {
		if queue.len() > GlobalConfig.SEND_MAX {
			duration = time.Microsecond * time.Duration(queue.len())
			if duration.Seconds() > float64(GlobalConfig.MAX_INTERVAL_WAIT_TIME) {
				duration = time.Second * time.Duration(GlobalConfig.MAX_INTERVAL_WAIT_TIME)
			}
		}
		if !queue.mute {
			event := queue.getQueueEvent()
			if queue.mute {
				queue.putQueueEvent(event)
			} else {
				go processSingleEvent(event)
			}
		}
		time.Sleep(duration)
	}
}

func processSingleEvent(event *QueueEvent) {
	switch event.status {
	case NEW_ALARM:
		lastID, _ := mydb.CreateStrategyEvent(event.strategy_event, event.trigger_events)
		event.strategy_event.ID = lastID
		go doAlarmAction(event.strategy_event, event.trigger_events)
	case OLD_ALARM:
		mydb.UpdateStrategyEvent(event.strategy_event, event.trigger_events)
		go doAlarmAction(event.strategy_event, event.trigger_events)
	case RESTORE_ALARM:
		mydb.UpdateStrategyEvent(event.strategy_event, event.trigger_events)
		mydb.CreateStrategyEventProcess(event.strategy_event.ID, event.strategy_event.Status, "系统", "报警恢复", time.Now().Format("2006-01-02 15:04:05"))
		go doRestoreAction(event.strategy_event, event.trigger_events)
		event.strategy_event.Count++
	default:
		lg.Error(fmt.Sprintf("unknow event type %d from event queue.", event.status))
	}
	mydb.CreateStrategyEventRecord(event.strategy_event, event.trigger_events)
}

//发送报警
func doAlarmAction(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) {
	strategy_event.Status = types.EVENT_NEW
	actions := mydb.GetAllActions(strategy_event.StrategyID)
	for _, action := range actionFilter(actions, strategy_event) {
		lg.Info("do action:%v", action)
		switch action.Kind {
		case types.ACTION_NOTIFY:
			subject := action.AlarmSubject
			content := fillTemplate(action.AlarmTemplate, generateTemplateObj(strategy_event, triggerEventFilter(trigger_events, types.ACTION_ALARM)))
			go broadcast(strategy_event, subject, content, action)
		case types.ACTION_RUN:
			go doRun(strategy_event, action)
		}
	}
}

//发送恢复通知
func doRestoreAction(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) {
	strategy_event.Status = types.EVENT_CLOSED
	actions := mydb.GetActions(strategy_event.StrategyID, 1)
	for _, action := range actions {
		switch action.Kind {
		case types.ACTION_NOTIFY:
			subject := action.RestoreSubject
			content := fillTemplate(action.RestoreTemplate, generateTemplateObj(strategy_event, triggerEventFilter(trigger_events, types.ACTION_RESTORE)))
			go broadcast(strategy_event, subject, content, action)
		case types.ACTION_RUN:
			go doRun(strategy_event, action)

		}
	}
}

func triggerEventFilter(trigger_events map[string]*types.TriggerEvent, action_type int) map[string]*types.TriggerEvent {
	new_trigger_events := make(map[string]*types.TriggerEvent)
	switch action_type {
	case types.ACTION_ALARM:
		for index, trigger_event := range trigger_events {
			if trigger_event.Triggered == true {
				new_trigger_events[index] = trigger_event
			}
		}
	case types.ACTION_RESTORE:
		for index, trigger_event := range trigger_events {
			if trigger_event.TriggerChanged == true {
				new_trigger_events[index] = trigger_event
			}
		}
	}
	return new_trigger_events
}

func actionFilter(actions []*types.Action, strategy_event *types.StrategyEvent) []*types.Action {
	filtered := make([]*types.Action, 0)
	for _, action := range actions {
		now := time.Now()
		if !(now.After(generateTime(now, action.BeginTime)) && now.Before(generateTime(now, action.EndTime))) {
			continue
		}
		if !(now.Sub(strategy_event.CreateTime).Minutes() >= float64(action.TimePeriod)) {
			continue
		}
		filtered = append(filtered, action)
	}
	return filtered
}

func generateTime(nowDateTime time.Time, onlyTime string) time.Time {
	newTime, err := time.Parse("15:04:05", onlyTime)
	if err != nil {
		lg.Error(err.Error())
	}
	return time.Date(nowDateTime.Year(), nowDateTime.Month(), nowDateTime.Day(), newTime.Hour(), newTime.Minute(), newTime.Second(), 0, time.Local)
}

func broadcast(event *types.StrategyEvent, subject, content string, action *types.Action) {
	users := make(map[int]*types.User)
	users_obj_from_group := mydb.GetUsersByGroups(action.ID)
	for _, user := range users_obj_from_group {
		users[user.ID] = user
	}
	for _, user := range users {
		go func(user *types.User) {
			defer func() {
				if r := recover(); r != nil {
					lg.Error("when send notify to user %s occur error %s", user.Username, r)
				}
			}()
			params := make([]string, 0)
			params = append(params, subject)
			params = append(params, content)
			user_info, err := json.Marshal(user)
			if err != nil {
				lg.Error(err.Error())
			}
			params = append(params, string(user_info))
			script := mydb.GetScript(action.ScriptID)
			var success bool
			result, err := runScript(script.FilePath, params)
			if err != nil {
				success = false
				lg.Error("run script %s | %s | %s | %s", script.Name, script.FilePath, params, result)
			} else {
				success = true
				lg.Info("run script %s | %s | %s | %s", script.Name, script.FilePath, params, result)
			}
			action_result := types.NewActionResult(event.ID,
				event.Count,
				action.ID,
				action.Type,
				action.Kind,
				action.ScriptID,
				user.ID,
				user.Username,
				user.Phone,
				user.Mail,
				user.Wechat,
				params[0],
				params[1],
				result,
				success,
			)
			mydb.CreateActionResult(action_result)
		}(user)
	}
}

func doRun(event *types.StrategyEvent, action *types.Action) {
	params := make([]string, 0)
	params = append(params, event.IP)
	params = append(params, event.HostName)
	script := mydb.GetScript(action.ScriptID)
	var success bool
	result, err := runScript(script.FilePath, params)
	if err != nil {
		success = false
		lg.Error("run script %s | %s | %s | %s", script.Name, script.FilePath, params, result)
	} else {
		success = true
		lg.Info("run script %s | %s | %s | %s", script.Name, script.FilePath, params, result)
	}
	action_result := types.NewActionResult(event.ID,
		event.Count,
		action.ID,
		action.Type,
		action.Kind,
		action.ScriptID,
		0,
		"",
		"",
		"",
		"",
		"",
		"",
		result,
		success)
	mydb.CreateActionResult(action_result)

}

//运行脚本函数，并返回执行结果
func runScript(file_path string, params []string) (string, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
		done   chan error = make(chan error, 1)
	)

	cmd := exec.Command(file_path, params...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Start(); err != nil {
		lg.Error("run script %v error %v", file_path, err.Error())
		return err.Error(), err
	}

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Second * time.Duration(GlobalConfig.ACTION_TIMEOUT)):
		if err := cmd.Process.Kill(); err != nil {
			lg.Error("kill script process error ", err.Error())
			return err.Error(), err
		}
		lg.Warn("run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT)
		return fmt.Sprintf("run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT), errors.New(fmt.Sprintf("run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT))
	case err := <-done:
		if err != nil {
			lg.Error("run script %v error %v output %v", file_path, err.Error(), stderr.String())
			return stderr.String(), err
		}
		return stdout.String(), err
	}
}

//根据模板对象填充模板
func fillTemplate(raw_template string, template_obj notifyTemplate) string {
	var filled_template bytes.Buffer
	tmpl, err := template.New("template").Parse(raw_template)
	if err != nil {
		lg.Error(err.Error())
		return err.Error()
	}
	err = tmpl.Execute(&filled_template, template_obj)
	if err != nil {
		lg.Error(err.Error())
		return err.Error()
	}
	return filled_template.String()
}

//根据报警事件生成报警模板对象
func generateTemplateObj(strategy_event *types.StrategyEvent, trigger_events map[string]*types.TriggerEvent) notifyTemplate {
	nt := notifyTemplate{}
	nt.STRATEGY.ID = strategy_event.ID
	nt.STRATEGY.NAME = strategy_event.StrategyName
	nt.STRATEGY.CYCLE = strategy_event.Cycle
	nt.STRATEGY.STATUS = STRATEGY_STATUS_MAPPING[strategy_event.Status]
	nt.STRATEGY.PRIORITY = STRATEGY_PRIORITY_MAPPING[strategy_event.Priority]
	nt.STRATEGY.ALARM_COUNT = strategy_event.AlarmCount
	nt.STRATEGY.COUNT = strategy_event.Count
	nt.STRATEGY.UPDATE_TIME = strategy_event.UpdateTime.Format("2006-01-02 15:04:05")
	nt.STRATEGY.EXPRESSION = strategy_event.Expression
	for _, trigger_event := range trigger_events {
		nt.STRATEGY.EXPRESSION_DETAIL += trigger_event.String()
	}
	nt.HOST.NAME = strategy_event.HostName
	nt.HOST.IP = strategy_event.IP

	return nt
}

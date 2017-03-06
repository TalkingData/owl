package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"text/template"
	"time"

	"owl/common/tcp"
	. "owl/common/types"
)

var alarm_tasks *AlarmTasks

type AlarmTasks struct {
	tasks map[string]*AlarmTask
}

func (this *AlarmTasks) Add(task *AlarmTask) {
	if len(task.Triggers) == 0 {
		lg.Warn("Task %v has no triggers, skipped it.", task.ID)
		return
	}
	if task.Host.Status == "2" {
		lg.Debug("Host %v is forbidden, skipped it.", task.Host.ID)
		return
	}

	parentKey := fmt.Sprintf("%v@%v", task.Strategy.Pid, task.Host.ID)

	if _, ok := this.tasks[parentKey]; ok {
		delete(this.tasks, parentKey)
	}
	this.tasks[task.ID] = task
}

var controller *Controller

type Controller struct {
	*tcp.Server
	taskPool       *TaskPool
	highResultPool *ResultPool
	lowResultPool  *ResultPool
	nodePool       *NodePool
}

func InitController() error {
	controllerServer := tcp.NewServer(GlobalConfig.TCP_BIND, &ControllerHandle{})
	controllerServer.SetMaxPacketSize(uint32(GlobalConfig.MAX_PACKET_SIZE))
	if err := controllerServer.ListenAndServe(); err != nil {
		return err
	}
	lg.Info("Start listen: %v", GlobalConfig.TCP_BIND)

	controller = &Controller{controllerServer,
		NewTaskPool(GlobalConfig.TASK_POOL_SIZE),
		NewResultPool(GlobalConfig.RESULT_POOL_SIZE),
		NewResultPool(GlobalConfig.RESULT_POOL_SIZE),
		NewNodePool()}

	go controller.loadStrategiesiForever()
	go controller.processStrategyResultForever()
	go controller.checkNodesForever()
	return nil
}

func (this *Controller) checkNodesForever() {
	for {
		this.checkNodes()
		time.Sleep(time.Second * 10)
	}
}

func (this *Controller) checkNodes() {
	time_now := time.Now()
	for ip, node := range this.nodePool.Nodes {
		if time_now.Sub(node.Update) > time.Duration(time.Second*35) {
			delete(this.nodePool.Nodes, ip)
			lg.Warn("Inspector %v, %v lost from controller", ip, node.Hostname)
		}
	}
}

func (this *Controller) refreshNode(heartbeat *HeartBeat) {
	this.nodePool.Lock.Lock()
	defer this.nodePool.Lock.Unlock()

	if node, ok := this.nodePool.Nodes[heartbeat.IP]; ok {
		node.IP = heartbeat.IP
		node.Hostname = heartbeat.Hostname
		node.Update = time.Now()
	} else {
		node = &Node{}
		node.IP = heartbeat.IP
		node.Hostname = heartbeat.Hostname
		node.Update = time.Now()
		this.nodePool.Nodes[heartbeat.IP] = node
	}
}

func (this *Controller) loadStrategiesiForever() {
	duration := time.Duration(GlobalConfig.LOAD_STRATEGIES_INTERVAL) * time.Second
	for {
		alarm_tasks = &AlarmTasks{make(map[string]*AlarmTask)}
		for _, strategy := range mydb.GetStrategies() {
			this.loadTaskForStrategy(strategy)
		}
		this.taskPool.PutTasks(alarm_tasks.tasks)
		lg.Info("Loaded tasks %v", len(alarm_tasks.tasks))
		time.Sleep(duration)
	}
}

func (this *Controller) loadTaskForStrategy(strategy *Strategy) {
	if strategy.Enable == false {
		lg.Info("Strategy %v is not enabled, skipped it.", strategy.Name)
		return
	}

	global_hosts := make([]*Host, 0)

	for _, group := range mydb.GetGroupsByStrategyID(strategy.ID) {
		group_hosts := mydb.GetHostsByGroupID(group.ID)
		global_hosts = append(global_hosts, group_hosts...)
	}

	hosts := mydb.GetHostsByStrategyID(strategy.ID)
	global_hosts = append(global_hosts, hosts...)

	for _, host := range global_hosts {
		triggers := mydb.GetTriggersByStrategyID(strategy.ID)
		alarm_tasks.Add(NewAlarmTask(host, strategy, triggers))
	}

	lg.Info("Loaded tasks for strategy %v", strategy.Name)
}

func (this *Controller) processStrategyResultForever() {
	worker_count := GlobalConfig.WORKER_COUNT
	for worker_count > 0 {
		go this.highResultWorker()
		go this.lowResultWorker()
		worker_count -= 1
	}
}

func (this *Controller) highResultWorker() {
	for {
		select {
		case result := <-this.highResultPool.results:
			lg.Debug("Get result %v from high result pool", result.TaskID)
			this.processResult(result)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Controller) lowResultWorker() {
	for {
		select {
		case result := <-this.lowResultPool.results:
			lg.Debug("Get result %v from low result pool", result.TaskID)
			this.processResult(result)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Controller) doAlarmAction(host *Host, strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) {
	strategy_event.Status = EVENT_NEW
	actions := mydb.GetActions(strategy_event.StrategyID, ACTION_ALARM)
	for _, action := range actions {
		subject := action.AlarmSubject
		content := fillTemplate(action.AlarmTemplate, generateTemplateObj(host, strategy_event, triggerEventFilter(trigger_event_sets, ACTION_ALARM)))
		broadcastMessage(strategy_event.ID, subject, content, action)
	}
}

func (this *Controller) doRestoreAction(host *Host, strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) {
	strategy_event.Status = EVENT_CLOSED
	actions := mydb.GetActions(strategy_event.StrategyID, ACTION_RESTORE)
	for _, action := range actions {
		subject := action.RestoreSubject
		content := fillTemplate(action.RestoreTemplate, generateTemplateObj(host, strategy_event, triggerEventFilter(trigger_event_sets, ACTION_RESTORE)))
		broadcastMessage(strategy_event.ID, subject, content, action)
	}
}

func (this *Controller) doCustomAction(host *Host, strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) {
	//TODO when alarm event is triggered, the custom action would be run which is the script wrote by user.
}

func triggerEventFilter(trigger_event_sets map[string][]*TriggerEvent, action_type int) map[string][]*TriggerEvent {
	new_trigger_event_sets := make(map[string][]*TriggerEvent)
	switch action_type {
	case ACTION_ALARM:
		for index, trigger_event_set := range trigger_event_sets {
			for _, trigger_event := range trigger_event_set {
				if trigger_event.Triggered == true {
					new_trigger_event_sets[index] = append(new_trigger_event_sets[index], trigger_event)
				}
			}
		}
	case ACTION_RESTORE:
		for index, trigger_event_set := range trigger_event_sets {
			for _, trigger_event := range trigger_event_set {
				if trigger_event.Triggered == false {
					new_trigger_event_sets[index] = append(new_trigger_event_sets[index], trigger_event)
				}
			}
		}
	}
	return new_trigger_event_sets
}

func broadcastMessage(strategy_event_id int64, subject, content string, action *Action) {
	users := make(map[int]*User)
	users_obj_from_group := mydb.GetUsersByGroups(action.ID)
	users_obj_from_user := mydb.GetUsers(action.ID)
	for _, user := range users_obj_from_group {
		users[user.ID] = user
	}
	for _, user := range users_obj_from_user {
		users[user.ID] = user
	}

	for _, user := range users {
		go func(user *User) {
			params := make([]string, 0)
			params = append(params, subject)
			params = append(params, content)
			var file_path string
			switch action.SendType {
			case SEND_MAIL:
				file_path = GlobalConfig.SEND_MAIL_SCRIPT
				params = append(params, user.Mail)
			case SEND_SMS:
				file_path = GlobalConfig.SEND_SMS_SCRIPT
				params = append(params, user.Phone)
			case SEND_WECHAT:
				file_path = GlobalConfig.SEND_WECHAT_SCRIPT
				params = append(params, user.Weixin)
			case SEND_CALL:
				file_path = GlobalConfig.SEND_CALL_SCRIPT
				params = append(params, user.Phone)
			default:
				lg.Error("Unknown send type %v", action.SendType)
				return
			}
			result, err := runScript(file_path, params)
			action_result := &ActionResult{}
			action_result.StrategyEventID = strategy_event_id
			action_result.ActionID = action.ID
			action_result.ActionType = action.Type
			action_result.ActionSendType = action.SendType
			action_result.UserID = user.ID
			action_result.Username = user.Username
			action_result.Phone = user.Phone
			action_result.Mail = user.Mail
			action_result.Weixin = user.Weixin
			action_result.Subject = subject
			action_result.Content = content
			if err != nil {
				action_result.Success = false
			} else {
				action_result.Success = true
			}
			action_result.Response = result
			mydb.CreateActionResult(action_result)
		}(user)
	}
}

func fillTemplate(raw_template string, template_obj Template) string {
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

func generateTemplateObj(host *Host, strategy_event *StrategyEvent, trigger_event_sets map[string][]*TriggerEvent) Template {
	template := Template{}
	template.STRATEGY.ID = strategy_event.ID
	template.STRATEGY.NAME = strategy_event.StrategyName
	template.STRATEGY.CYCLE = strategy_event.Cycle
	template.STRATEGY.STATUS = STRATEGY_STATUS_MAPPING[strategy_event.Status]
	template.STRATEGY.PRIORITY = STRATEGY_PRIORITY_MAPPING[strategy_event.Priority]
	template.STRATEGY.ALARM_COUNT = strategy_event.AlarmCount
	template.STRATEGY.COUNT = strategy_event.Count
	template.STRATEGY.UPDATE_TIME = strategy_event.UpdateTime.Format("2006-01-02 15:04:05")
	template.STRATEGY.EXPRESSION = strategy_event.Expression
	for _, trigger_event_set := range trigger_event_sets {
		for _, trigger_event := range trigger_event_set {
			template.STRATEGY.EXPRESSION_DETAIL += trigger_event.String()
		}
	}
	template.HOST.CNAME = host.Name
	template.HOST.NAME = host.Hostname
	template.HOST.IP = host.IP
	template.HOST.STATUS = host.Status
	template.HOST.SN = host.SN

	return template
}

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
		lg.Error("Run script %v params %v error %v", file_path, params, err.Error())
		return err.Error(), err
	}

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Second * time.Duration(GlobalConfig.ACTION_TIMEOUT)):
		if err := cmd.Process.Kill(); err != nil {
			lg.Error("Kill script process error ", err.Error())
			return err.Error(), err
		}
		lg.Warn("Run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT)
		return fmt.Sprintf("Run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT), errors.New(fmt.Sprintf("Run script %v timeout in %v", file_path, GlobalConfig.ACTION_TIMEOUT))
	case err := <-done:
		if err != nil {
			lg.Error("Run script %v params %v error %v output %v", file_path, params, err.Error(), stderr.String())
			return stderr.String(), err
		}
		return stdout.String(), err
	}
}

func (this *Controller) processResult(strategy_result *StrategyResult) {
	if strategy_result.ErrorMessage != "" {
		lg.Error(fmt.Sprintf("Task %v has problem %v", strategy_result.TaskID, strategy_result.ErrorMessage))
		return
	}
	task, ok := alarm_tasks.tasks[strategy_result.TaskID]
	if !ok {
		lg.Error(fmt.Sprintf("Task %v not in cached task pool", strategy_result.TaskID))
		return
	}

	strategy_event, trigger_event_sets := generateEvent(strategy_result, task)

	new_strategy_event := mydb.GetStrategyEvent(strategy_event.StrategyID, EVENT_NEW, strategy_event.HostID)
	if new_strategy_event != nil {
		strategy_event.ID = new_strategy_event.ID
		new_strategy_event.UpdateTime = strategy_result.CreateTime
	}
	aware_strategy_event := mydb.GetStrategyEvent(strategy_event.StrategyID, EVENT_AWARE, strategy_event.HostID)
	if aware_strategy_event != nil {
		strategy_event.ID = aware_strategy_event.ID
		aware_strategy_event.UpdateTime = strategy_result.CreateTime
	}

	if aware_strategy_event != nil {
		if strategy_result.Triggered == false {
			this.doRestoreAction(task.Host, strategy_event, trigger_event_sets)
			aware_strategy_event.Status = EVENT_CLOSED
			aware_strategy_event.ProcessUser = "系统"
			aware_strategy_event.ProcessComments = "报警恢复"
			aware_strategy_event.ProcessTime = time.Now()
		}
		mydb.UpdateStrategyEvent(aware_strategy_event, trigger_event_sets)
		return
	}

	if new_strategy_event != nil {
		if strategy_result.Triggered == false {
			this.doRestoreAction(task.Host, strategy_event, trigger_event_sets)
			new_strategy_event.Status = EVENT_CLOSED
			new_strategy_event.ProcessUser = "系统"
			new_strategy_event.ProcessComments = "报警恢复"
			new_strategy_event.ProcessTime = time.Now()
			mydb.UpdateStrategyEvent(new_strategy_event, trigger_event_sets)
			return
		}
		if new_strategy_event.Count < task.Strategy.AlarmCount || task.Strategy.AlarmCount == 0 {
			new_strategy_event.Count += 1
			this.doAlarmAction(task.Host, strategy_event, trigger_event_sets)
		}
		mydb.UpdateStrategyEvent(new_strategy_event, trigger_event_sets)
		return
	}

	if new_strategy_event == nil && aware_strategy_event == nil {
		if strategy_result.Triggered == true {
			strategy_event.Status = EVENT_NEW
			last_id, _ := mydb.CreateStrategyEvent(strategy_event, trigger_event_sets)
			strategy_event.ID = last_id
			this.doAlarmAction(task.Host, strategy_event, trigger_event_sets)
		}
	}
}

func generateEvent(strategy_result *StrategyResult, task *AlarmTask) (*StrategyEvent, map[string][]*TriggerEvent) {

	var strategy_event *StrategyEvent
	trigger_event_sets := make(map[string][]*TriggerEvent)

	strategy_event = NewStrategyEvent(task.Strategy.ID,
		task.Strategy.Name,
		task.Strategy.Priority,
		task.Strategy.Cycle,
		task.Strategy.AlarmCount,
		task.Strategy.Expression,
		strategy_result.CreateTime,
		task.Host.ID,
		task.Host.Name,
		task.Host.Hostname,
		task.Host.IP,
		task.Host.SN)

	for index, trigger_result_set := range strategy_result.TriggerResultSets {
		trigger := task.Triggers[index]
		for _, trigger_result := range trigger_result_set.TriggerResults {
			trigger_event_sets[index] = append(trigger_event_sets[index], NewTriggerEvent(strategy_event.ID, index, trigger.Metric, trigger_result.Tags, trigger_result.AggregateTags, trigger.Symbol, trigger.Method, trigger.Number, trigger.Threshold, trigger_result.CurrentThreshold, trigger_result.Triggered))
		}
	}

	return strategy_event, trigger_event_sets
}

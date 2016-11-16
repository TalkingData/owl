package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"owl/common/tcp"
	"owl/common/types"
)

type Inspector struct {
	*tcp.Server
	session    *tcp.Session
	taskPool   *TaskPool
	resultPool *ResultPool
}

var inspector *Inspector

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hostname
}

func InitInspector() error {
	server := tcp.NewServer("", &InspectorHandle{})
	server.SetMaxPacketSize(uint32(GlobalConfig.MAX_PACKET_SIZE))
	inspector = &Inspector{server, nil, NewTaskPool(GlobalConfig.MAX_TASK_BUFFER), NewResultPool(GlobalConfig.MAX_RESULT_BUFFER)}
	go inspector.DialForever()
	go inspector.HeartBeatForever()
	go inspector.GetInspectorTasksForever()
	go inspector.ProcessInspectorTasksForever()
	go inspector.SendResultForever()
	return nil
}

func (this *Inspector) Dial() {
retry:
	session, err := this.Connect(GlobalConfig.CONTROLLER_ADDR, nil)
	if err != nil || session.IsClosed() {
		lg.Error("Can not connect to controller %v, error: %v", GlobalConfig.CONTROLLER_ADDR, err)
		time.Sleep(time.Second * time.Duration(3))
		goto retry
	}
	this.session = session
	lg.Info("Inspector connected to controller: %v", GlobalConfig.CONTROLLER_ADDR)
}

func (this *Inspector) DialForever() {
	for {
		if this.session == nil || this.session.IsClosed() {
			this.Dial()
		}
		time.Sleep(time.Second * time.Duration(3))
	}
}

func (this *Inspector) HeartBeatForever() {
	for {
		if this.session != nil {
			this.session.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_HEARTBEAT, types.NewHeartBeat(this.session.LocalAddr(), GetHostName())))
		}
		time.Sleep(time.Second * 30)
	}
}

func (this *Inspector) GetInspectorTasksForever() {
	for {
		if len(this.taskPool.tasks) == 0 && this.session != nil {
			this.session.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_TASK_REQUEST, types.NewHeartBeat(this.session.LocalAddr(), GetHostName())))
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (this *Inspector) SendResultForever() {
	for {
		select {
		case result := <-this.resultPool.results:
			this.session.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_RESULT, result))
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Inspector) ProcessInspectorTasksForever() {
	worker_count := GlobalConfig.WORKER_COUNT
	for worker_count > 0 {
		go this.taskWorker()
		worker_count -= 1
	}
}

func (this *Inspector) taskWorker() {
	for {
		select {
		case task := <-this.taskPool.tasks:
			this.processTask(task)
			lg.Debug("Get task %v from task pool", task.ID)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Inspector) processTask(task *types.AlarmTask) {
	info := strings.SplitN(task.ID, "@", 3)
	host_id := info[2]
	triggers_results := make(map[string]*types.TriggerResultSet)
	parameters := make(map[string]interface{})
	error_message := ""
	task.Strategy.Cycle += 1
	for index, trigger := range task.Triggers {
		var trigger_result_set *types.TriggerResultSet
		var err error
		switch trigger.Method {
		case MAX_METHOD:
			trigger_result_set, err = maxMethod(host_id, task.Strategy.Cycle, trigger)
		case MIN_METHOD:
			trigger_result_set, err = minMethod(host_id, task.Strategy.Cycle, trigger)
		case RATIO_METHOD:
			trigger_result_set, err = ratioMethod(host_id, task.Strategy.Cycle, trigger)
		case TOP_METHOD:
			trigger_result_set, err = topMethod(host_id, task.Strategy.Cycle, trigger)
		case BOTTOM_METHOD:
			trigger_result_set, err = bottomMethod(host_id, task.Strategy.Cycle, trigger)
		case LAST_METHOD:
			trigger_result_set, err = lastMethod(host_id, task.Strategy.Cycle, trigger)
		case NODATA_METHOD:
			trigger_result_set, err = nodataMethod(host_id, task.Strategy.Cycle, trigger)
		default:
			lg.Error("Trigger method %v not found", trigger.Method)
			err = errors.New(fmt.Sprintf("Trigger method %v not found", trigger.Method))
		}

		if err != nil {
			error_message = err.Error()
			break
		}

		if len(trigger_result_set.TriggerResults) == 0 {
			error_message = "no data"
			break
		}
		triggers_results[index] = trigger_result_set
		parameters[index] = trigger_result_set.Triggered
	}

	result := false
	if error_message == "" {
		compute_result, err := compute(parameters, task.Strategy.Expression)
		if err != nil {
			error_message = err.Error()
		}
		result = compute_result
	}

	strategy_result := types.NewStrategyResult(task.ID, task.Strategy.Priority, triggers_results, error_message, result, time.Now())
	this.resultPool.PutResult(strategy_result)
}

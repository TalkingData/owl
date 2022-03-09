package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"owl/common/types"

	"github.com/wuyingsong/tcp"
)

type Inspector struct {
	controller *tcp.TCPConn
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

	inspector = &Inspector{nil, NewTaskPool(GlobalConfig.MAX_TASK_BUFFER), NewResultPool(GlobalConfig.MAX_RESULT_BUFFER)}
	if err := inspector.Dial(); err != nil {
		return err
	}
	go inspector.DialForever()
	go inspector.HeartBeatForever()
	go inspector.GetInspectorTasksForever()
	go inspector.ProcessInspectorTasksForever()
	go inspector.SendResultForever()
	return nil
}

func (this *Inspector) Dial() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MAX_PACKET_SIZE))
	tcpAddr, err := net.ResolveTCPAddr("tcp", GlobalConfig.CONTROLLER_ADDR)
	if err != nil {
		return err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	this.controller = tcp.NewTCPConn(tcpConn, &callback{}, protocol)
	this.controller.AsyncWritePacket(
		tcp.NewDefaultPacket(types.ALAR_MESS_INSPECTOR_HEARTBEAT, types.NewHeartBeat(this.controller.GetLocalAddr().String(), GetHostName()).Encode()),
	)
	// this.session.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_HEARTBEAT, types.NewHeartBeat(this.session.LocalAddr(), GetHostName())))
	lg.Info("inspector connected to controller: %v", GlobalConfig.CONTROLLER_ADDR)
	return this.controller.Serve()
}

func (this *Inspector) DialForever() {
	for {
		if this.controller == nil || this.controller.IsClosed() {
			this.Dial()
		}
		time.Sleep(time.Second * time.Duration(3))
	}
}

func (this *Inspector) HeartBeatForever() {
	for {
		if this.controller != nil {
			heartBeat := types.NewHeartBeat(
				this.controller.GetLocalAddr().String(),
				GetHostName(),
			)
			this.controller.AsyncWritePacket(
				tcp.NewDefaultPacket(
					types.ALAR_MESS_INSPECTOR_HEARTBEAT,
					heartBeat.Encode(),
				),
			)
		}
		time.Sleep(time.Second * 5)
	}
}

func (this *Inspector) GetInspectorTasksForever() {
	for {
		if len(this.taskPool.tasks) == 0 && this.controller != nil {
			this.controller.AsyncWritePacket(
				tcp.NewDefaultPacket(
					types.ALAR_MESS_INSPECTOR_TASK_REQUEST,
					[]byte{},
				),
			)
		}
		time.Sleep(time.Second * 1)
	}
}

func (this *Inspector) SendResultForever() {
	var err error
	for result := range this.resultPool.results {
		for {
			if err = this.controller.AsyncWritePacket(
				tcp.NewDefaultPacket(types.ALAR_MESS_INSPECTOR_RESULT, result.Encode()),
			); err == nil {
				break
			}
			lg.Warn("send task result to controller failed error %s", err)
			time.Sleep(time.Second)
		}
		lg.Info("send task result to controller %s result:%s", result.Encode(), GlobalConfig.CONTROLLER_ADDR)
	}
}

// for {
// 	select {
// 	case result := <-this.resultPool.results:
// 		result_buffer.Results = append(result_buffer.Results, result)
// 		if len(result_buffer.Results) == GlobalConfig.RESULT_BUFFER {
// 			this.controller.AsyncWritePacket(
// 				tcp.NewDefaultPacket(types.ALAR_MESS_INSPECTOR_RESULT, result_buffer.Encode()),
// 			)
// 			lg.Info("Send %d Results to controller", len(result_buffer.Results))
// 			result_buffer.Results = result_buffer.Results[:0]
// 		}
// 	default:
// 		if len(result_buffer.Results) > 0 {
// 			this.controller.AsyncWritePacket(
// 				tcp.NewDefaultPacket(types.ALAR_MESS_INSPECTOR_RESULTS, result_buffer.Encode()),
// 			)
// 			// this.session.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_RESULTS, result_buffer))
// 			lg.Info("Send %d Results to controller", len(result_buffer.Results))
// 			result_buffer.Results = result_buffer.Results[:0]
// 		}
// 		time.Sleep(time.Millisecond * 100)
// 	}
// }
// }

func (this *Inspector) ProcessInspectorTasksForever() {

	for wc := GlobalConfig.WORKER_COUNT; wc > 0; wc-- {
		go this.taskWorker()
	}
}

func (this *Inspector) taskWorker() {
	for {
		select {
		case task := <-this.taskPool.tasks:
			this.processTask(task)
			lg.Debug("get task %v from task pool", task.ID)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Inspector) processTask(task *types.AlarmTask) {
	triggers_results := make(map[string]*types.TriggerResultSet)
	parameters := make(map[string]interface{})
	host_id := strings.SplitN(task.ID, "@", 2)[1]
	error_message := ""
	for index, trigger := range task.Triggers {
		var trigger_result_set *types.TriggerResultSet
		var err error
		if trigger.Tags == "" {
			trigger.Tags = fmt.Sprintf("uuid=%s", host_id)
		} else {
			trigger.Tags = fmt.Sprintf("%s,uuid=%s", trigger.Tags, host_id)
		}
		switch trigger.Method {
		case MAX_METHOD:
			trigger_result_set, err = maxMethod(task.Strategy.Cycle, trigger)
		case MIN_METHOD:
			trigger_result_set, err = minMethod(task.Strategy.Cycle, trigger)
		case RATIO_METHOD:
			trigger_result_set, err = ratioMethod(task.Strategy.Cycle, trigger)
		case TOP_METHOD:
			trigger_result_set, err = topMethod(task.Strategy.Cycle, trigger)
		case BOTTOM_METHOD:
			trigger_result_set, err = bottomMethod(task.Strategy.Cycle, trigger)
		case LAST_METHOD:
			trigger_result_set, err = lastMethod(task.Strategy.Cycle, trigger)
		case DIFF_METHOD:
			trigger_result_set, err = diffMethod(task.Strategy.Cycle, trigger)
		case NODATA_METHOD:
			trigger_result_set, err = nodataMethod(task.Strategy.Cycle, trigger)
		case AVG_METHOD:
			trigger_result_set, err = avgMethod(task.Strategy.Cycle, trigger)
		default:
			err = errors.New(fmt.Sprintf("trigger method %v not found", trigger.Method))
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

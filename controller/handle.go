package main

import (
	"owl/common/tcp"
	"owl/common/types"
)

type ControllerHandle struct {
}

func (this *ControllerHandle) MakeSession(sess *tcp.Session) {
	lg.Info("%s new connection ", sess.RemoteAddr())
}

func (this *ControllerHandle) LostSession(sess *tcp.Session) {
	lg.Info("%s disconnect ", sess.RemoteAddr())
}

func (this *ControllerHandle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			lg.Error("Recovered in HandleMessage", err)
		}
	}()
	mt := types.AlarmMessageType(data[0])
	switch mt {
	case types.ALAR_MESS_INSPECTOR_HEARTBEAT:
		lg.Info("Receive %v %v", types.AlarmMessageTypeText[mt], string(data[1:]))
		heartbeat := &types.HeartBeat{}
		if err := heartbeat.Decode(data[1:]); err != nil {
			lg.Error(err.Error())
			return
		}
		controller.receiveHearbeat(heartbeat)
	case types.ALAR_MESS_INSPECTOR_TASK_REQUEST:
		tasks := GetAlarmTasks()
		sess.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_TASKS, tasks))
	case types.ALAR_MESS_INSPECTOR_RESULTS:
		lg.Info("Receive %v %v", types.AlarmMessageTypeText[mt], string(data[1:]))
		results := &types.AlarmResults{}
		if err := results.Decode(data[1:]); err != nil {
			lg.Error(err.Error())
			return
		}
		controller.resultPool.PutResults(results)
	default:
		lg.Error("Unknown option: %v", mt)
	}
}

func GetAlarmTasks() *types.AlarmTasks {
	tasks := controller.taskPool.GetTasks(GlobalConfig.TASK_SIZE)
	tasks_resp := &types.AlarmTasks{tasks}
	return tasks_resp
}

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
		controller.refreshNode(heartbeat)
	case types.ALAR_MESS_INSPECTOR_TASK_REQUEST:
		tasks_resp := getInspectorTask()
		sess.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_TASKS, tasks_resp))
	case types.ALAR_MESS_INSPECTOR_RESULT:
		lg.Info("Receive %v %v", types.AlarmMessageTypeText[mt], string(data[1:]))
		result := &types.StrategyResult{}
		if err := result.Decode(data[1:]); err != nil {
			lg.Error(err.Error())
			return
		}
		switch result.Priority {
		case types.PRIORITY_HIGH_LEVEL:
			controller.highResultPool.PutResult(result)
		case types.PRIORITY_MIDDLE_LEVEL:
			controller.lowResultPool.PutResult(result)
		case types.PRIORITY_LOW_LEVEL:
			controller.lowResultPool.PutResult(result)
		default:
			lg.Error("Unknown priority: %v", result.Priority)
		}
	default:
		lg.Error("Unknown option: %v", mt)
	}
}

func getInspectorTask() *types.GetTasksResp {
	tasks := controller.taskPool.GetTasks(GlobalConfig.TASK_SIZE)
	tasks_resp := &types.GetTasksResp{tasks}
	return tasks_resp
}

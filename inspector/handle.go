package main

import (
	"owl/common/tcp"
	"owl/common/types"
)

type InspectorHandle struct {
}

func (this *InspectorHandle) MakeSession(sess *tcp.Session) {
	lg.Info("%s new connection ", sess.RemoteAddr())
}

func (this *InspectorHandle) LostSession(sess *tcp.Session) {
	lg.Info("%s disconnect ", sess.RemoteAddr())
}

func (this *InspectorHandle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			lg.Error("recovered in HandleMessage", err)
		}
	}()
	mt := types.AlarmMessageType(data[0])
	switch mt {
	case types.ALAR_MESS_INSPECTOR_TASKS:
		at := &types.AlarmTasks{}
		if err := at.Decode(data[1:]); err != nil {
			lg.Error(err.Error())
			return
		}
		if len(at.Tasks) == 0 {
			return
		}
		lg.Info("receive %v %v", types.AlarmMessageTypeText[mt], string(data[1:]))
		inspector.taskPool.PutTasks(at.Tasks)
	default:
		lg.Error("unknown option: %v", mt)
	}
}

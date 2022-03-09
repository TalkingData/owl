package main

import (
	"owl/common/types"

	"github.com/wuyingsong/tcp"
)

type callback struct {
}

func (cb *callback) OnConnected(conn *tcp.TCPConn) {
	lg.Info("callback:%s connected", conn.GetRemoteAddr().String())
}

//链接断开回调
func (cb *callback) OnDisconnected(conn *tcp.TCPConn) {
	lg.Info("callback:%s disconnect ", conn.GetRemoteAddr().String())
}

//错误回调
func (cb *callback) OnError(err error) {
	lg.Error("callback: %s", err)
}

func (cb *callback) OnMessage(conn *tcp.TCPConn, p tcp.Packet) {
	defer func() {
		if r := recover(); r != nil {
			lg.Error("Recovered in OnMessage", r)
		}
	}()
	pkt := p.(*tcp.DefaultPacket)
	switch pkt.Type {
	case types.ALAR_MESS_INSPECTOR_HEARTBEAT:
		lg.Debug("receive %v %v", types.AlarmMessageTypeText[pkt.Type], string(pkt.Body))
		heartbeat := &types.HeartBeat{}
		if err := heartbeat.Decode(pkt.Body); err != nil {
			lg.Error(err.Error())
			return
		}
		controller.receiveHearbeat(heartbeat)
	//TODO: optimized task allocate algorithm
	case types.ALAR_MESS_INSPECTOR_TASK_REQUEST:
		lg.Debug("receive get task request inspector %s", conn.GetRemoteAddr())
		tasks := &types.AlarmTasks{
			Tasks: controller.taskPool.getTasks(GlobalConfig.TASK_SIZE),
		}
		if len(tasks.Tasks) == 0 {
			lg.Debug("task pool has no tasks can be sent to inspector %s", conn.GetRemoteAddr())
			return
		}
		lg.Info("sent %d task to inspector %s", len(tasks.Tasks), conn.GetRemoteAddr())
		conn.AsyncWritePacket(
			tcp.NewDefaultPacket(
				types.ALAR_MESS_INSPECTOR_TASKS,
				tasks.Encode(),
			),
		)
		// sess.Send(types.AlarmPack(types.ALAR_MESS_INSPECTOR_TASKS, tasks))
	case types.ALAR_MESS_INSPECTOR_RESULT:
		lg.Info("receive %v %v", types.AlarmMessageTypeText[pkt.Type], string(pkt.Body))
		var (
			result types.StrategyResult
			err    error
		)
		if err = result.Decode(pkt.Body); err != nil {
			lg.Error(err.Error())
			return
		}
		if err = controller.resultPool.putResult(&result); err != nil {
			lg.Error("put task result to pool failed error:%s result:%s", err, result.Encode())
		}
	default:
		lg.Error("unknown option: %v", pkt.Type)
	}
}

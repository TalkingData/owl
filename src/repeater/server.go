package main

import (
	"fmt"
	"owl/common/types"
	"owl/repeater/backend"
	"time"

	"github.com/wuyingsong/tcp"
)

var (
	repeater *Repeater
)

type Repeater struct {
	srv     *tcp.AsyncTCPServer
	buffer  chan []byte
	backend []Backend
}

func InitRepeater() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MaxPacketSize))
	s := tcp.NewAsyncTCPServer(GlobalConfig.TCPBind, &callback{}, protocol)
	repeater = &Repeater{}
	repeater.srv = s
	repeater.buffer = make(chan []byte, GlobalConfig.BufferSize)

	var (
		err error
		bk  Backend
	)

	for _, b := range GlobalConfig.Backend {
		switch b {
		case "opentsdb":
			bk, err = backend.NewOpentsdbBackend(GlobalConfig.OpentsdbAddr)
		case "kairosdb":
			bk, err = backend.NewOpentsdbBackend(GlobalConfig.KairosdbAddr)
		case "repeater":
			bk, err = backend.NewRepeaterBackend(GlobalConfig.RepeaterAddr)
		case "kafka":
			bk, err = backend.NewKafkaBackend(GlobalConfig.KafkaBrokers, GlobalConfig.KafkaTopic)
		default:
			err = fmt.Errorf("unsupported backend %s", GlobalConfig.Backend)
		}
		if err != nil {
			return fmt.Errorf("%s error:%s", b, err)
		}
		lg.Info("new time series backend:%s", b)
		repeater.backend = append(repeater.backend, bk)
	}
	return repeater.srv.ListenAndServe()
}

func (this *Repeater) Forward() {
	var (
		err error
	)
	for data := range this.buffer {
		tsd := &types.TimeSeriesData{}
		if err = tsd.Decode(data); err != nil {
			lg.Error("decode error %s ", err)
			continue
		}
		//时间对齐
		tsd.Timestamp = tsd.Timestamp - (tsd.Timestamp % int64(tsd.Cycle))
		for index, b := range this.backend {
			for {
				if err = b.Write(tsd); err == nil {
					lg.Notice("forward to %s %v", GlobalConfig.Backend[index], *tsd)
					break
				}
				if GlobalConfig.BackendWriteFailedPolicy == SkippedPolicy {
					lg.Error("write backend %s failed, error:%s, skipped", GlobalConfig.Backend[index], err)
					break
				}
				lg.Error("write backend %s failed, error:%s, retry", GlobalConfig.Backend[index], err)
				time.Sleep(time.Second * 5)
			}
		}
	}
}

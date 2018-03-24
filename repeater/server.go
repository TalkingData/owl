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
	backend Backend
}

func InitRepeater() error {
	protocol := &tcp.DefaultProtocol{}
	protocol.SetMaxPacketSize(uint32(GlobalConfig.MaxPacketSize))
	s := tcp.NewAsyncTCPServer(GlobalConfig.TCPBind, &callback{}, protocol)
	repeater = &Repeater{}
	repeater.srv = s
	repeater.buffer = make(chan []byte, GlobalConfig.BufferSize)
	var err error
	switch GlobalConfig.Backend {
	case "opentsdb", "kairosdb":
		repeater.backend, err = backend.NewOpentsdbBackend(GlobalConfig.OpentsdbAddr)
	case "repeater":
		repeater.backend, err = backend.NewRepeaterBackend(GlobalConfig.RepeaterAddr)
	case "kafka":
		repeater.backend, err = backend.NewKafkaBackend(GlobalConfig.KafkaBrokers, GlobalConfig.KafkaTopic)
	default:
		err = fmt.Errorf("unsupported backend %s", GlobalConfig.Backend)
	}
	if err != nil {
		return fmt.Errorf("%s error:%s", GlobalConfig.Backend, err)
	}
	return repeater.srv.ListenAndServe()
}

func (this *Repeater) Forward() {
	var err error
	for {
		select {
		case data := <-this.buffer:
			tsd := &types.TimeSeriesData{}
			if err = tsd.Decode(data); err != nil {
				lg.Error("decode error %s ", err)
				continue
			}
			for {
				err = this.backend.Write(tsd)
				if err == nil {
					break
				}
				lg.Error("forward to %s error(%s)", GlobalConfig.Backend, err)
				time.Sleep(time.Second * 5)
			}
			lg.Notice("forward to %s %v", GlobalConfig.Backend, *tsd)
		}
	}
}

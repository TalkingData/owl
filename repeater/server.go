package main

import (
	"fmt"
	"owl/common/tcp"
	"owl/common/types"
	"owl/repeater/backend"
	"time"
)

var (
	repeater *Repeater
)

type Repeater struct {
	srv     *tcp.Server
	buffer  chan []byte
	backend Backend
}

func InitRepeater() error {
	s := tcp.NewServer(GlobalConfig.TCP_BIND, &handle{})
	repeater = &Repeater{}
	repeater.srv = s
	repeater.buffer = make(chan []byte, GlobalConfig.BUFFER_SIZE)
	var err error
	switch GlobalConfig.BACKEND {
	case "opentsdb":
		repeater.backend, err = backend.NewOpentsdbBackend(GlobalConfig.OPENTSDB_ADDR)
	case "repeater":
		repeater.backend, err = backend.NewRepeaterBackend(GlobalConfig.REPEATER_ADDR)
	case "kafka":
		repeater.backend, err = backend.NewKafkaBackend(GlobalConfig.KAFKA_BROKERS, GlobalConfig.KAFKA_TOPIC)
	default:
		err = fmt.Errorf("unsupported backend %s", GlobalConfig.BACKEND)
	}
	if err != nil {
		return fmt.Errorf("%s error:%s", GlobalConfig.BACKEND, err)
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
				continue
			}
			if err = this.backend.Write(tsd); err != nil {
				this.buffer <- data
				time.Sleep(time.Second * 1)
			} else {
				lg.Notice("forward to %s %v", GlobalConfig.BACKEND, *tsd)
			}
		}
	}
}

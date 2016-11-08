package main

import (
	"errors"
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
	switch GlobalConfig.BACKEND {
	case "opentsdb":
		backend, err := backend.NewOpentsdbBackend(GlobalConfig.OPENTSDB_ADDR)
		if err != nil {
			return err
		}
		repeater.backend = backend
	case "repeater":
		backend, err := backend.NewRepeaterBackend(GlobalConfig.REPEATER_ADDR)
		if err != nil {
			return err
		}
		repeater.backend = backend
	default:
		return errors.New(fmt.Sprintf("unsupported backend %s", GlobalConfig.BACKEND))
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

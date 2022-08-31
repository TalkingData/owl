package conf

import (
	"owl/common/global"
	"time"
)

type constConf struct {
	ServiceName string

	HttpServerShutdownTimeoutSecs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: global.OwlAgentServiceName,

		HttpServerShutdownTimeoutSecs: 5 * time.Second,
	}
}

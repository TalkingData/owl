package conf

import (
	"owl/common/global"
	"time"
)

type constConf struct {
	ServiceName    string
	RpcRegisterKey string

	MetricServerShutdownTimeoutSecs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.OwlRepeaterServiceName,
		RpcRegisterKey: global.OwlRepeaterRpcRegisterKey,

		MetricServerShutdownTimeoutSecs: 5 * time.Second,
	}
}

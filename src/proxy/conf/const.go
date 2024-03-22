package conf

import (
	"owl/common/global"
	"time"
)

type constConf struct {
	ServiceName string

	CfcServiceName      string
	RepeaterServiceName string

	MetricServerShutdownTimeoutSecs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: global.OwlProxyServiceName,

		CfcServiceName:      global.OwlCfcServiceName,
		RepeaterServiceName: global.OwlRepeaterServiceName,

		MetricServerShutdownTimeoutSecs: 5 * time.Second,
	}
}

package conf

import "time"

type constConf struct {
	ServiceName string

	HttpServerShutdownTimeoutSecs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: "owl-agent",

		HttpServerShutdownTimeoutSecs: 5 * time.Second,
	}
}

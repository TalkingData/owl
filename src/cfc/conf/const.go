package conf

import (
	"owl/common/global"
	"time"
)

type constConf struct {
	ServiceName    string
	RpcRegisterKey string

	MasterRpcRegisterKey string

	ExecSqlBatchLimit      int
	ExecSqlBatchIntervalMs int

	MetricServerShutdownTimeoutSecs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.OwlCfcServiceName,
		RpcRegisterKey: global.OwlCfcRpcRegisterKey,

		MasterRpcRegisterKey: global.OwlCfcMasterRegisterKey,

		ExecSqlBatchLimit:      500,
		ExecSqlBatchIntervalMs: 100,

		MetricServerShutdownTimeoutSecs: 5 * time.Second,
	}
}

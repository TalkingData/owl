package conf

import "owl/common/global"

type constConf struct {
	ServiceName string

	CleanExpiredMetricBatchLimit      int
	CleanExpiredMetricBatchIntervalMs int
	SetHostAliveBatchLimit            int
	SetHostAliveBatchIntervalMs       int
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: global.OwlCfcServiceName,

		CleanExpiredMetricBatchLimit:      500,
		CleanExpiredMetricBatchIntervalMs: 100,
		SetHostAliveBatchLimit:            500,
		SetHostAliveBatchIntervalMs:       100,
	}
}

package conf

type constConf struct {
	ServiceName string

	CleanExpiredMetricBatchLimit      int
	CleanExpiredMetricBatchIntervalMs int
	SetHostAliveBatchLimit            int
	SetHostAliveBatchIntervalMs       int
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: "owl-cfc",

		CleanExpiredMetricBatchLimit:      500,
		CleanExpiredMetricBatchIntervalMs: 100,
		SetHostAliveBatchLimit:            500,
		SetHostAliveBatchIntervalMs:       100,
	}
}

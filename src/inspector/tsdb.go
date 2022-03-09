package main

import (
	"owl/common/tsdb"
	"time"
)

var tsdbClient tsdb.TsdbClient

func InitTsdb() (err error) {
	switch GlobalConfig.BACKEND_TSDB {
	case "opentsdb":
		tsdbClient, err = tsdb.NewOpenTsdbClient(GlobalConfig.TSDB_ADDR, time.Duration(GlobalConfig.TSDB_TIMEOUT)*time.Second)
	case "kairosdb":
		tsdbClient, err = tsdb.NewKairosDbClient(GlobalConfig.TSDB_ADDR)
	}
	return err
}

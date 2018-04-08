package main

import (
	"fmt"
	"net/http"
	"owl/common/tsdb"
	"time"

	"github.com/gin-gonic/gin"
)

var tsdbClient tsdb.TsdbClient

func initTSDB() error {
	var err error
	switch config.TimeSeriesStorage {
	case "opentsdb":
		tsdbClient, err = tsdb.NewOpenTsdbClient(config.OpentsdbAddr, time.Duration(config.OpenttsdbReadTimeout)*time.Second)
	case "kairosdb":
		tsdbClient, err = tsdb.NewKairosDbClient(config.KairosdbAddr)
	default:
		err = fmt.Errorf("%s timeseries storage not support", config.TimeSeriesStorage)
	}
	return err
}

func queryTimeSeriesData(c *gin.Context) {
	metric := c.Query("metric")
	tags := c.Query("tags")
	start := c.DefaultQuery("start", time.Now().Add(-time.Hour).Format("2006/01/02-15:04:05"))
	end := c.DefaultQuery("end", time.Now().Format("2006/01/02-15:04:05"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	result, err := tsdbClient.Query(start, end, tags, "sum", metric, false)
	if err != nil {
		response["message"] = err.Error()
		response["code"] = http.StatusInternalServerError
		return
	}
	response["data"] = result
	response["code"] = http.StatusOK
}

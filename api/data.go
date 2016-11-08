package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/types"
	"time"
)

func data(c *gin.Context) {
	metric := c.Query("metric")
	tags := c.Query("tags")
	start := c.DefaultQuery("start", "1h-ago")
	end := c.Query("end")
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	client, err := types.NewTsdbClient(GlobalConfig.OPENTSDB_ADDR, time.Duration(5)*time.Second)
	if err != nil {
		response["message"] = err.Error()
		response["code"] = http.StatusInternalServerError
		return
	}
	defer client.Close()
	params := types.NewQueryParams(start, end, tags, "sum", metric)
	result, err := client.Query(params)
	if err != nil {
		response["message"] = err.Error()
		response["code"] = http.StatusInternalServerError
		return
	}
	response["data"] = result
	response["code"] = http.StatusOK
}

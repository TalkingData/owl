package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MetricSummary struct {
	ID       int    `json:"id"`
	Metric   string `json:"metric"`            //sys.cpu.idle
	DataType string `json:"data_type" db:"dt"` //COUNTER,GAUGE,DERIVE
	Cycle    int    `json:"cycle,omitempty"`
	Tags     string `json:"tags"` //{"product":"app01", "group":"dev02"}
	UpdateAt string `json:"update_at" db:"update_at"`
}

func suggestMetrics(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	productID, err := strconv.Atoi(c.DefaultQuery("product_id", "0"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	response["metrics"] = mydb.suggestMetrics(productID)
}

func suggestMetricTagSet(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	productID, err := strconv.Atoi(c.DefaultQuery("product_id", "0"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	metric := strings.TrimSpace(c.Query("metric"))
	if len(metric) == 0 {
		response["code"] = http.StatusBadRequest
		return
	}
	response["tag_set"] = mydb.suggestMetricTagSet(productID, metric)
}

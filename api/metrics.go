package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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

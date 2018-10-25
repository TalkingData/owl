package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Chart struct {
	ID       int             `json:"id"`
	Title    string          `json:"title"`
	Creator  string          `json:"creator"`
	Span     int             `json:"span"`
	Height   int             `json:"height"`
	CreateAt string          `json:"create_at"`
	Type     int             `json:"type"`
	Elements []*ChartElement `json:"elements"`
}

type ChartElement struct {
	Metric string `json:"metric"`
	Tags   string `json:"tags"`
}

func listPanelChats(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, charts := mydb.getPanelCharts(
		c.GetInt("panel_id"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetBool("paging"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["charts"] = charts
}

func createPanelChart(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		chart *Chart
		err   error
	)
	if err = c.BindJSON(&chart); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	chart.Creator = c.GetString("username")
	if chart, err = mydb.createPanelChart(c.GetInt("panel_id"), chart); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["code"] = http.StatusOK
	response["chart"] = chart
}

func updatePanelChart(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	chart := Chart{}
	if err := c.BindJSON(&chart); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.updatePanelChart(c.GetInt("panel_id"), &chart); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	response["chart"] = chart
}

func deletePanelChart(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	chartID, err := strconv.Atoi(c.Param("chart_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deletePanelChart(c.GetInt("panel_id"), chartID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

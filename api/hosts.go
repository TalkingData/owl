package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Host struct {
	ID           string  `json:"id"`
	IP           string  `json:"ip"`
	Name         string  `json:"name"`
	Hostname     string  `json:"hostname"`
	AgentVersion string  `json:"agent_version" db:"agent_version"`
	Status       string  `json:"status"`
	CreateAt     string  `json:"create_at" db:"create_at"`
	UpdateAt     string  `json:"update_at" db:"update_at"`
	Uptime       float64 `json:"uptime" db:"uptime"`
	IdlePct      float64 `json:"idle_pct" db:"idle_pct"`
	MuteTime     string  `json:"mute_time" db:"mute_time"`
}

type WarpHost struct {
	Host
	Apps []string `json:"apps"`
}

func listAllHosts(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	order := c.GetString("order")
	if order == "" {
		order = "status asc"
	}
	total, hosts := mydb.getAllHosts(
		c.GetBool("paging"),
		c.Query("query"),
		order,
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["hosts"] = warpHosts(hosts)
}

func deleteHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	if err := mydb.deleteHost(hostID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

func listHostMetrics(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		hostID string
		host   *Host
	)
	hostID = c.Param("host_id")
	if host = mydb.getHostByID(hostID); host.ID == "" {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("host [%s] is not exists", hostID)
		return
	}
	total, metrics := mydb.getHostMetrics(
		hostID,
		c.GetBool("paging"),
		c.Query("prefix"),
		c.Query("query"),
		"metric asc",
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	for _, metric := range metrics {
		if metric.Tags != "" {
			metric.Tags = fmt.Sprintf("host=%s,%s", host.Hostname, metric.Tags)
		} else {
			metric.Tags = fmt.Sprintf("host=%s", host.Hostname)
		}
	}
	response["metrics"] = metrics
	response["total"] = total
}

func listHostApps(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		hostID string
		host   *Host
	)
	hostID = c.Param("host_id")
	if host = mydb.getHostByID(hostID); host.ID == "" {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("host [%s] is not exists", hostID)
		return
	}

}

func warpHosts(hosts []Host) []WarpHost {
	warpHosts := []WarpHost{}
	for _, host := range hosts {
		warpHosts = append(warpHosts, WarpHost{host, mydb.getHostAppNames(host.ID)})
	}
	return warpHosts
}

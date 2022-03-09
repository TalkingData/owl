package main

import (
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"
	"time"

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
	Apps       []string `json:"apps,omitempty" db:"apps"`
	PluginCnt  int      `json:"plugin_cnt" db:"plugin_cnt"`
	HostGroups string   `json:"host_groups,omitempty" db:"groups"`
	Products   string   `json:"products,omitempty" db:"products"`
}

func getHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var host *Host
	hostID := c.Param("host_id")
	if host = mydb.getHostByID(hostID); host.ID == "" {
		if host = mydb.getHostByIP(hostID); host.ID == "" {
			response["code"] = http.StatusNotFound
			return
		}
	}
	response["host"] = host
	response["products"] = mydb.getHostProducts(host.ID)
	response["apps"] = mydb.getHostAppNames(host.ID)
}

//TODO: 优化查询性能
func listAllHosts(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	order := c.GetString("order")
	if order == "" {
		order = "status asc"
	}
	var (
		noProduct bool
		err       error
	)
	noProduct, err = strconv.ParseBool(c.DefaultQuery("no_product", "false"))
	if err != nil {
		noProduct = false
	}
	total, hosts := mydb.getAllHosts(
		c.GetBool("paging"),
		noProduct,
		c.GetString("query"),
		order,
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["hosts"] = hosts
}

func deleteHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	if err := mydb.deleteHost(hostID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	mydb.CleanupHostEvents(hostID)
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
		c.GetString("query"),
		"metric asc",
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["metrics"] = metrics
	response["total"] = total
}

func deleteHostMetrics(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil || len(ids.IDS) == 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}

	if err := mydb.deleteHostMetrics(hostID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
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
	response["apps"] = mydb.getHostAppNames(hostID)
}

func muteHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	muteTime := c.Query("mute_time")
	if len(hostID) == 0 {
		return
	}
	if err := mydb.muteHost(hostID, muteTime); err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["host"] = mydb.getHostByID(hostID)
}

func unmuteHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	if len(hostID) == 0 {
		return
	}
	unmuteTime := time.Now().Add(-time.Hour * 24).Format("2006-01-02 15:04:05")
	if err := mydb.muteHost(hostID, unmuteTime); err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["host"] = mydb.getHostByID(hostID)
}

func listHostPlugins(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, plugins := mydb.getHostPlugins(
		c.Param("host_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		"metric asc",
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["plugins"] = plugins
	response["total"] = total
}

func createHostPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var err error
	plugin := &types.Plugin{}
	if err = c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err = plugin.Validate(); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	hostID := c.Param("host_id")
	if plugin, err = mydb.createHostPlugin(hostID, plugin); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["plugin"] = plugin
}

func updateHostPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	plugin := &types.Plugin{}
	var err error
	if err = c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err = plugin.Validate(); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err = mydb.updateHostPlugin(c.Param("host_id"), plugin); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["plugin"] = plugin
}

func deleteHostPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostID := c.Param("host_id")
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = "plugin id avalied"
	}

	if err := mydb.deleteHostPlugin(hostID, pluginID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

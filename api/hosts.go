package main

import (
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"

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
	Apps       []string    `json:"apps"`
	PluginCnt  int         `json:"plugin_cnt"`
	HostGroups []HostGroup `json:"host_groups,omitempty"`
	Products   []Product   `json:"products,omitempty"`
}

func getHost(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	response["host"] = warpHost(*mydb.getHostByID(c.Param("host_id")))
}

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
		warpHosts = append(warpHosts, warpHost(host))
	}
	return warpHosts
}

func warpProductHosts(productID int, hosts []Host) []WarpHost {
	warpHosts := []WarpHost{}
	for _, host := range hosts {
		warpHosts = append(warpHosts, warpProductHost(productID, host))
	}
	return warpHosts
}

func warpHost(host Host) WarpHost {
	pluginCnt, _ := mydb.getHostPlugins(host.ID, false, "", "", 0, 0)
	return WarpHost{
		host,
		mydb.getHostAppNames(host.ID),
		pluginCnt,
		mydb.getHostHostGroups(0, host.ID),
		mydb.getHostProducts(host.ID),
	}
}

func warpProductHost(productID int, host Host) WarpHost {
	pluginCnt, _ := mydb.getHostPlugins(host.ID, false, "", "", 0, 0)
	return WarpHost{
		host,
		mydb.getHostAppNames(host.ID),
		pluginCnt,
		mydb.getHostHostGroups(productID, host.ID),
		nil,
		// mydb.getHostProducts(host.ID),
	}
}

func listHostPlugins(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, plugins := mydb.getHostPlugins(
		c.Param("host_id"),
		c.GetBool("paging"),
		c.Query("query"),
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

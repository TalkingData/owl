package main

import (
	"log"
	"net/http"
	"owl/common/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PluginHostGroup struct {
	HostGroup
	Product string `json:"product"`
}

func listPlugins(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, plugins := mydb.getPlugins(
		c.GetString("query"),
		c.GetBool("paging"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["plugins"] = plugins
}

func createPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		plugin *types.Plugin
		err    error
	)
	if err = c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	username := c.GetString("username")
	if plugin, err = mydb.createPlugin(plugin, username); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["code"] = http.StatusOK
	response["plugin"] = plugin
}

func updatePlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	plugin := types.Plugin{}
	if err := c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.updatePlugin(plugin); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	response["plugin"] = plugin
}

func deletePlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deletePlugin(pluginID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

func addHostGroups2Plugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.addHostGroups2Plugin(pluginID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}
func removeHostGroupsFromPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.removeHostGroupsFromPlugin(pluginID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
}

func listPluginHostGroups(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	total, hostGroups := mydb.getPluginHostGroups(
		pluginID,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["host_groups"] = hostGroups
	response["total"] = total
}

func listNotInPluginHostGroups(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	total, hostGroups := mydb.getNotInPluginHostGroups(
		pluginID,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["host_groups"] = hostGroups
	response["total"] = total
}

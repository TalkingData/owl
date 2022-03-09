package main

import (
	"fmt"
	"log"
	"net/http"
	"owl/common/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HostGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	CreateAt    string `json:"create_at" db:"create_at"`
	UpdateAt    string `json:"update_at" db:"update_at"`
}

type WarpHostGroup struct {
	HostGroup
	PluginCnt   int `json:"plugin_cnt,omitempty" db:"plugin_cnt"`
	HostCnt     int `json:"host_cnt,omitempty" db:"host_cnt"`
	StrategyCnt int `json:"strategy_cnt,omitempty" db:"strategy_cnt"`
}

func listProductHostGroupHosts(c *gin.Context) {
	response := gin.H{"status": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, hosts := mydb.getProductHostGroupHosts(
		c.GetInt("product_id"),
		c.GetInt("host_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func listNotInProductHostGroupHosts(c *gin.Context) {
	response := gin.H{"status": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, hosts := mydb.getNotInProductHostGroupHosts(
		c.GetInt("product_id"),
		c.GetInt("host_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func listProductHostGroups(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var username string
	if c.DefaultQuery("my", "false") == "true" {
		username = c.GetString("username")
	}
	total, hostGroups := mydb.getProductHostGroups(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		username,
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["host_groups"] = hostGroups
}

func updateProductHostGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostGroup := HostGroup{}
	if err := c.BindJSON(&hostGroup); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.updateHostGroup(hostGroup); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["host_group"] = mydb.getProductHostGroupByID(
		c.GetInt("product_id"),
		hostGroup.ID)
}

func createProductHostGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostGroup := HostGroup{}
	if err := c.BindJSON(&hostGroup); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	productID := c.GetInt("product_id")
	hostGroup.Creator = c.GetString("username")
	if err := mydb.createProductHostGroup(productID, hostGroup); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["host_group"] = mydb.getProductHostGroupByName(
		productID,
		hostGroup.Name)
}

func deleteProductHostGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	hostGroupID, err := strconv.Atoi(c.Param("host_group_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deleteProductHostGroup(c.GetInt("product_id"), hostGroupID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

func addHosts2ProductHostGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []string `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	groupID := c.GetInt("host_group_id")
	if err := mydb.addHost2HostGroup(groupID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	total, hosts := mydb.getProductHostGroupHosts(
		c.GetInt("product_id"),
		groupID,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func removeHostsFromProductHostGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []string `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	groupID := c.GetInt("host_group_id")
	if err := mydb.removeHostFromHostGroup(groupID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	total, hosts := mydb.getProductHostGroupHosts(
		c.GetInt("product_id"),
		groupID,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func listHostGroupPlugins(c *gin.Context) {
	response := gin.H{"status": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	fmt.Println(c.Param("host_group_id"))
	total, plugins := mydb.getHostGroupPlugins(
		c.GetInt("host_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["plugins"] = plugins
}

func updateHostGroupPlugin(c *gin.Context) {
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
	if err = mydb.updateHostGroupPlugin(c.GetInt("host_group_id"), plugin); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["plugin"] = plugin
}

func createHostGroupPlugin(c *gin.Context) {
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
	if plugin, err = mydb.createHostGroupPlugin(c.GetInt("host_group_id"), plugin); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["plugin"] = plugin
}

func deleteHostGroupPlugin(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	pluginID, err := strconv.Atoi(c.Param("plugin_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deleteHostGroupPlugin(c.GetInt("host_group_id"), pluginID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

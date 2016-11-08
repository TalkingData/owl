package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/types"
	"strconv"
)

type Plugin struct {
	types.Plugin
	Hosts []*types.Host `json:"hosts" gorm:"many2many:host_plugin;"`
}

func pluginList(c *gin.Context) {
	q := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))

	plugins := []*Plugin{}
	db := mydb.Table("plugin").Order("create_at desc")
	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		db = db.Where("phone like ? or weixin like ? or mail like ?", q, q, q)
	}
	var total int
	db.Count(&total)
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Find(&plugins)

	for _, plugin := range plugins {
		mydb.Model(&plugin).Association("Hosts").Find(&plugin.Hosts)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "plugins": plugins, "total": total})
}

func pluginCreate(c *gin.Context) {
	plugin := Plugin{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	cnt := 0
	mydb.Table("plugin").Where("name=?", plugin.Name).Count(&cnt)
	if cnt > 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = "plugin already exists"
		return
	}
	if err := mydb.Set("gorm:save_associations", false).Save(&plugin).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&plugin).Association("Hosts").Replace(plugin.Hosts).Find(&plugin.Hosts)

	response["code"] = http.StatusOK
	response["message"] = "user create successful"
	response["plugin"] = plugin
}

func pluginDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	plugin := Plugin{}
	if mydb.Table("plugin").Where("id = ?", id).Find(&plugin).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The plugin not found"
		return
	}
	if err := mydb.Delete(&plugin).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["code"] = http.StatusOK
	response["message"] = fmt.Sprintf("%s delete successful", plugin.Name)
}

func pluginUpdate(c *gin.Context) {
	plugin := Plugin{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)

	if err := c.BindJSON(&plugin); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.Table("plugin").Save(&plugin).Association("Hosts").Replace(plugin.Hosts).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["plugin"] = plugin
	response["code"] = http.StatusOK
	response["message"] = "plugin update successful"
}

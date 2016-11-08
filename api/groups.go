package main

import (
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Group struct {
	types.Group
	Hosts []*types.Host `json:"hosts" gorm:"many2many:host_group"`
}

func groupList(c *gin.Context) {
	groups := []*types.Group{}
	q := c.DefaultQuery("q", "")
	db := mydb.Table("group")
	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ?", q)
	}
	db.Find(&groups)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "groups": groups})
}

func groupUpdate(c *gin.Context) {
	group := Group{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&group); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.Set("gorm:save_associations", false).Table("group").Save(&group).
		Association("Hosts").Replace(group.Hosts).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["code"] = http.StatusOK
	response["message"] = "group update successful"
	response["group"] = group
}

func groupCreate(c *gin.Context) {
	group := Group{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&group); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	cnt := 0
	mydb.Table("group").Where("name=?", group.Name).Count(&cnt)
	if cnt > 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = "group already exists"
		return
	}
	if err := mydb.Set("gorm:save_associations", false).Save(&group).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&group).Association("Hosts").Replace(group.Hosts).Find(&group.Hosts)

	response["code"] = http.StatusOK
	response["message"] = "group create successful"
	response["group"] = group
}

func groupDelete(c *gin.Context) {
	id := c.Param("id")
	group := types.Group{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	mydb.Table("group").Where("id = ?", id).First(&group)
	if group.ID == 0 {
		response["code"] = http.StatusNotFound
		response["message"] = "group not found"
	} else {
		var cnt int
		mydb.Table("host_group").Where("group_id = ?", group.ID).Count(&cnt)
		mydb.Delete(&group)
		response["message"] = "delete success"
	}
	c.JSON(http.StatusOK, response)
}

func groupHostList(c *gin.Context) {
	id := c.Param("id")
	group := types.Group{}
	response := gin.H{"code": http.StatusOK}
	mydb.Table("group").Where("id = ?", id).First(&group)
	if group.ID == 0 {
		response["code"] = http.StatusNotFound
		response["message"] = "group not found"
	} else {
		hosts := []*Host{}
		response["hosts"] = hosts
		db := mydb.Table("host").Joins("JOIN host_group ON host.id = host_group.host_id").
			Where("host_group.group_id = ?", group.ID)
		//根据状态过滤主机
		if status := c.Query("status"); len(status) > 0 {
			db = db.Where("status = ?", status)
		}
		q := fmt.Sprintf("%%%s%%", c.Query("q"))
		db = db.Where("name like ? or hostname like ? or ip like ?", q, q, q)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))

		offset := (page - 1) * pageSize
		db.Offset(offset).Limit(pageSize).Find(&hosts)

		metrics := 0
		if len(hosts) > 0 {
			for _, h := range hosts {
				//获取主机关联的组
				mydb.Joins("JOIN host_group ON host_group.group_id = group.id").
					Where("host_group.host_id = ?", h.ID).Find(&h.Groups)

				//获取metric数量
				mydb.Table("metric").Where("host_id = ?", h.ID).Count(&h.Metrics)

				//获取主机应用的组策略数量
				mydb.Table("strategy").Joins("JOIN strategy_group ON strategy_group.strategy_id = strategy.id").
					Joins("JOIN host_group ON host_group.group_id = strategy_group.group_id").
					Where("host_group.host_id = ?", h.ID).Count(&h.GroupStrategies)

				//获取主机应用的主机策略数量
				mydb.Table("strategy").Joins("JOIN strategy_host ON strategy_host.strategy_id = strategy.id").
					Where("strategy_host.host_id = ?", h.ID).Count(&h.HostStrategies)

				metrics += h.Metrics

			}
			response["hosts"] = hosts
			response["metric_count"] = metrics
			response["host_count"] = len(hosts)
		}
	}
	c.JSON(http.StatusOK, response)
}

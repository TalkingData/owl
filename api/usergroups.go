package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/types"
	"strconv"
)

type userGroup struct {
	types.UserGroup
	Users []*types.User `json:"users" gorm:"many2many:user_user_group;"`
}

func userGroupList(c *gin.Context) {
	q := c.DefaultQuery("q", "")

	groups := []*userGroup{}
	db := mydb.Table("user_group")
	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ?", q)
	}
	db.Find(&groups)
	for _, group := range groups {
		mydb.Model(&group).Association("Users").Find(&group.Users)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "groups": groups})
}

func userGroupCreate(c *gin.Context) {
	group := userGroup{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&group); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	var cnt int
	if mydb.Table("user_group").Where("name=?", group.Name).Count(&cnt); cnt != 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = "user group already exists"
		return
	}
	if err := mydb.Set("gorm:save_associations", false).Save(&group).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&group).Association("Users").Replace(group.Users).Find(&group.Users)

	response["code"] = http.StatusOK
	response["message"] = "group create successful"
	response["group"] = group
}

func userGroupUpdate(c *gin.Context) {
	groups := []*userGroup{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&groups); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	var message string
	for _, group := range groups {

		if mydb.Table("user_group").Where("id=?", group.ID).First(&group).RecordNotFound() {
			response["code"] = http.StatusBadRequest
			message += "user group not found\n"
			continue
		}
		if err := mydb.Set("gorm:save_associations", false).Save(&group).Association("Users").Replace(group.Users).Error; err != nil {
			response["code"] = http.StatusInternalServerError
			message += err.Error()
			continue
		}
	}
	response["message"] = message
	response["code"] = http.StatusOK
	response["groups"] = groups
	response["message"] = "user group update successful"
}

func userGroupDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	group := userGroup{}
	if mydb.Table("user_group").Where("id = ?", id).Find(&group).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The usergroup not found"
		return
	}
	if err := mydb.Delete(&group).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return

	}
	response["code"] = http.StatusOK
	response["message"] = fmt.Sprintf("%s delete successful", group.Name)
}

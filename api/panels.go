package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/types"
	"strconv"
)

func panelList(c *gin.Context) {
	panels := []*types.Panel{}
	response := gin.H{"code": http.StatusOK}
	db := mydb.Table("panel")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))
	q := c.DefaultQuery("q", "")
	user_id := GetUser(c).ID
	favor, _ := strconv.Atoi(c.DefaultQuery("favor", "0"))

	if favor == 1 {
		db = db.Where("favor = 1")
	}
	db = db.Where("user_id = ? ", user_id)

	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		db = db.Where("name like ?", q)
	}
	total := 0
	db.Count(&total)
	response["total"] = total
	if page != 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	db.Find(&panels)
	for _, panel := range panels {
		mydb.Model(&panel).Association("Charts").Find(&panel.Charts)
		for _, chart := range panel.Charts {
			mydb.Model(&chart).Related(&chart.Elements)
		}
	}
	response["panels"] = panels
	c.JSON(http.StatusOK, response)
}

func panelCreate(c *gin.Context) {
	panel := types.Panel{}
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&panel); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if !mydb.Table("panel").Where("user_id = ? and name = ?", GetUser(c).ID, panel.Name).First(&panel).RecordNotFound() {
		response["code"] = http.StatusBadRequest
		response["message"] = "panel already exists"
		return
	}
	panel.UserID = GetUser(c).ID
	if err := mydb.Select("name", "thumbnail", "user_id").Save(&panel).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&panel).Association("Charts").Replace(panel.Charts)
	mydb.Model(&panel).Find(&panel)
	for _, chart := range panel.Charts {
		mydb.Model(&chart).Related(&chart.Elements)
	}
	response["panel"] = panel
}

func panelUpdate(c *gin.Context) {
	panel := types.Panel{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&panel); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return

	}
	if mydb.Table("panel").Where("id = ? and user_id = ?", panel.ID, GetUser(c).ID).Find(&types.Panel{}).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The panel not found"
		return
	}
	if err := mydb.Model(&panel).Select("name", "thumbnail").Updates(&panel).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&panel).Association("Charts").Replace(panel.Charts)
	for _, chart := range panel.Charts {
		mydb.Model(&chart).Related(&chart.Elements)
	}
	response["code"] = http.StatusOK
	response["message"] = "panel update successful"
	response["panel"] = panel
}

func panelDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	panel := types.Panel{}
	if mydb.Table("panel").Where("id = ? and user_id = ?", id, GetUser(c).ID).Find(&panel).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The panel not found"
		return
	}
	if err := mydb.Delete(&panel).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["code"] = http.StatusOK
	response["message"] = fmt.Sprintf("%s delete successful", panel.Name)
}

func getPanelChart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	panel := types.Panel{}
	charts := []*types.Chart{}
	if mydb.Table("panel").Where("id = ?", id).Find(&panel).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The panel not found"
		return
	}
	mydb.Model(&panel).Association("Charts").Find(&charts)
	for _, chart := range charts {
		mydb.Model(chart).Related(&chart.Elements)
	}
	response["charts"] = charts
}

func panelFavor(c *gin.Context) {
	cnt := 0
	mydb.Table("panel").Where("user_id=? and favor = 1", GetUser(c).ID).Count(&cnt)
	if cnt > 3 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotModified,
			"message": "favor the panel can not be more than four",
		})
		return
	}
	Favor(c, 1)
}

func panelUnFavor(c *gin.Context) {
	Favor(c, 0)
}

func Favor(c *gin.Context, favor int) {
	id, _ := strconv.Atoi(c.Param("id"))
	user_id := GetUser(c).ID
	panel := types.Panel{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if mydb.Table("panel").Where("id = ? and user_id = ?", id, user_id).First(&panel).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "panel not found"
		return
	}
	panel.Favor = favor
	mydb.Table("panel").Select("favor").Save(&panel)
	response["code"] = http.StatusOK
	response["panel"] = panel
}

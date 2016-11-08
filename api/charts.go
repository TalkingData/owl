package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/types"
	"strconv"
)

func chartList(c *gin.Context) {
	charts := []*types.Chart{}
	response := gin.H{"code": http.StatusOK}
	db := mydb.Table("chart").Where("user_id = ?", GetUser(c).ID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))
	q := c.DefaultQuery("q", "")
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
	db.Find(&charts)
	for _, chart := range charts {
		mydb.Model(chart).Related(&chart.Elements)
	}
	response["charts"] = charts
	c.JSON(http.StatusOK, response)
}

func chartCreate(c *gin.Context) {
	chart := types.Chart{}
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&chart); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if !mydb.Table("chart").Where("user_id = ? and name = ?", GetUser(c).ID, chart.Name).First(&chart).RecordNotFound() {
		response["code"] = http.StatusBadRequest
		response["message"] = "chart already exists"
		return
	}
	chart.UserID = GetUser(c).ID
	if err := mydb.Omit("create_at", "update_at").Create(&chart).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&chart).Find(&chart)
	response["chart"] = chart
}

func chartUpdate(c *gin.Context) {
	chart := types.Chart{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&chart); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if mydb.Table("chart").Where("id = ? and user_id = ?", chart.ID, GetUser(c).ID).Find(&types.Chart{}).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The chart not found"
		return
	}
	if err := mydb.Model(&chart).Set("gorm:save_associations", false).
		Omit("create_at", "update_at", "user_id").Save(&chart).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	for _, ele := range chart.Elements {
		ele.ChartID = chart.ID
	}
	mydb.Model(chart).Association("Elements").Replace(&chart.Elements)
	mydb.Table("chart_element").Where("chart_id is NULL").Delete(types.ChartElement{})
	response["code"] = http.StatusOK
	response["message"] = "chart update successful"
	response["chart"] = chart
}

func chartDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	chart := types.Chart{}
	if mydb.Table("chart").Where("id = ? and user_id = ?", id, GetUser(c).ID).Find(&chart).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The chart not found"
		return
	}
	if err := mydb.Delete(&chart).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["code"] = http.StatusOK
	response["message"] = fmt.Sprintf("%s delete successful", chart.Name)
}

func chartDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	chart := types.Chart{}
	if mydb.Table("chart").Where("id = ? and user_id = ?", id, GetUser(c).ID).Find(&chart).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The chart not found"
		return
	}
	mydb.Model(&chart).Related(&chart.Elements)
	response["code"] = http.StatusOK
	response["chart"] = chart
}

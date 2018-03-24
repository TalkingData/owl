package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StrategyTemplateDetail 策略模板详细
type StrategyTemplateDetail struct {
	StrategyTemplate
	TriggerTemplates []*TriggerTemplate `json:"triggers"`
}

// StrategyTemplate 策略模板
type StrategyTemplate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	AlarmCount  int    `json:"alarm_count" db:"alarm_count"`
	Cycle       int    `json:"cycle"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
}

// TriggerTemplate 表达式模板
type TriggerTemplate struct {
	ID                 int64   `json:"id" `
	StrategyTemplateID int     `json:"strategy_template_id" db:"strategy_template_id"`
	Metric             string  `json:"metric"`
	Tags               string  `json:"tags"`
	Number             int     `json:"number"`
	Index              string  `json:"index" `
	Method             string  `json:"method" `
	Symbol             string  `json:"symbol" `
	Threshold          float64 `json:"threshold" `
	Description        string  `json:"description"`
}

func templateGet(c *gin.Context) {
	templateID, err := strconv.Atoi(c.Param("template_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest})
	}
	template := mydb.GetStrategyTemplateDetail(templateID)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "template": &template})
}

func templateList(c *gin.Context) {
	templates := mydb.GetStrategyTemplates()
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "templates": &templates})
}

func templateCreate(c *gin.Context) {
	var std StrategyTemplateDetail
	if err := c.BindJSON(&std); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := mydb.CreateStrategyTemplate(&std); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func templateUpdate(c *gin.Context) {
	var std StrategyTemplateDetail
	if err := c.BindJSON(&std); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := mydb.UpdateStrategyTemplate(&std); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func templateDelete(c *gin.Context) {
	templateIDs := c.QueryArray("template_id")
	if len(templateIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params not be applied"})
		return
	}
	if err := mydb.DeleteStrategyTemplates(templateIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Script 执行脚本结构体
type Script struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	FilePath string `db:"file_path" json:"file_path"`
}

func scriptList(c *gin.Context) {
	scripts := mydb.GetScripts()
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "scripts": &scripts})
}

func scriptCreate(c *gin.Context) {
	var script Script
	if err := c.BindJSON(&script); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := mydb.CreateScript(&script); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func scriptUpdate(c *gin.Context) {
	var script Script
	if err := c.BindJSON(&script); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if err := mydb.UpdateScript(&script); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func scriptDelete(c *gin.Context) {
	scriptIDs := c.QueryArray("script_id")
	if len(scriptIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params not be applied"})
		return
	}
	if err := mydb.DeleteScript(scriptIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

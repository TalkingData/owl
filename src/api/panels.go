package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Panel struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
}

func listProductPanel(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, panels := mydb.getProductPanels(
		c.GetInt("product_id"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetBool("paging"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["panels"] = panels
}

func createProductPanel(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		panel *Panel
		err   error
	)
	if err = c.BindJSON(&panel); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	panel.Creator = c.GetString("username")
	if panel, err = mydb.createProductPanel(c.GetInt("product_id"), panel); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["code"] = http.StatusOK
	response["panel"] = panel
}

func updateProductPanel(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	panel := Panel{}
	if err := c.BindJSON(&panel); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if err := mydb.updatePanel(panel); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	response["panel"] = panel
}

func deleteProductPanel(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	panelID, err := strconv.Atoi(c.Param("panel_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deletePanel(panelID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

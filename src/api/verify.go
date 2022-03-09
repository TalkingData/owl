package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func verifyAdminPermission(c *gin.Context) {
	user := mydb.getUserProfile(c.GetString("username"))
	if user == nil || !user.isAdmin() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
	}
}

func verifyAndInjectProductID(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest})
	}
	var (
		exists bool
		user   *User
	)
	user = mydb.getUserProfile(c.GetString("username"))
	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
	}

	if !user.isAdmin() {
		for _, p := range mydb.getUserProducts(user) {
			if productID == p.ID {
				exists = true
				break
			}
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		}
	}
	c.Set("product_id", productID)
}

func verifyAndInjectHostGroupID(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("host_group_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest})
	}
	c.Set("host_group_id", productID)
}

func verifyAndInjectUserGroupID(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("user_group_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest})
	}
	c.Set("user_group_id", productID)
}

func verifyAndInjectPanelID(c *gin.Context) {
	panelID, err := strconv.Atoi(c.Param("panel_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest})
	}
	c.Set("panel_id", panelID)
}

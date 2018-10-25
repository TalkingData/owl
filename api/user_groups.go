package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description" db:"description"`
}

func listProductUserGroups(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, userGroups := mydb.getProductUserGroups(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["user_groups"] = userGroups
	response["total"] = total
}

func listProductUserGroupUsers(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, users := mydb.getProductUserGroupUsers(
		c.GetInt("user_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func listNotInProductUserGroupUsers(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, users := mydb.getNotInProductUserGroupUsers(
		c.GetInt("product_id"),
		c.GetInt("user_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func createProductUserGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	productID := c.GetInt("product_id")
	var (
		err       error
		userGroup UserGroup
	)
	if err := c.BindJSON(&userGroup); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if mydb.findProductUserGroup(productID, userGroup.Name) != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("%s is exists", userGroup.Name)
		return
	}
	if userGroup, err = mydb.createProductUserGroup(productID, userGroup); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["user_group"] = userGroup
}

func updateProductUserGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var err error
	productID := c.GetInt("product_id")
	userGroup := UserGroup{}
	if err = c.BindJSON(&userGroup); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err = mydb.updateProductUserGroup(productID, userGroup); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["user_group"] = userGroup
}

func deleteProductUserGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	userGroupID, err := strconv.Atoi(c.Param("user_group_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deleteProductUserGroup(c.GetInt("product_id"), userGroupID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

func addUsers2ProductUserGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	groupID := c.GetInt("user_group_id")
	if err := mydb.addUsers2UserGroup(groupID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	total, users := mydb.getProductUserGroupUsers(
		c.GetInt("user_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func removeUsersFromProductUserGroup(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	groupID := c.GetInt("user_group_id")
	if err := mydb.removeUsersFromUserGroup(groupID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	total, users := mydb.getProductUserGroupUsers(
		c.GetInt("user_group_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

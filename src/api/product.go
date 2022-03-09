package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Creator     string `json:"creator"`
}

/*
	/api/v1/users/userID/products
	{
		"id":string,
		"name":string
	}
*/
func listUserProducts(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	username := c.GetString("username")
	var user *User
	if user = mydb.getUserProfile(username); user == nil {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("getUserProfile failed, username:%s", username)
		return
	}
	response["products"] = mydb.getUserProducts(user)
}

func listAllProducts(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, products := mydb.getProducts(
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.Query("is_delete"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	type comProduct struct {
		*Product
		HostCnt int `json:"host_cnt"`
		UserCnt int `json:"user_cnt"`
	}
	comProducts := []comProduct{}
	for _, p := range products {
		product := comProduct{
			p,
			0, 0,
		}
		product.UserCnt = mydb.getProductUsersCnt(p.ID, "")
		product.HostCnt = mydb.getProductHostsCnt(p.ID, false, "")
		comProducts = append(comProducts, product)
	}
	response["total"] = total
	response["products"] = comProducts
}

func createProduct(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	username := c.GetString("username")
	var (
		user    *User
		product *Product
		err     error
	)
	if user = mydb.getUserProfile(username); user == nil {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("get user id failed, username:%s", username)
		return
	}
	if err = c.BindJSON(&product); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if p := mydb.getProductByName(product.Name); p != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = fmt.Sprintf("%s is exists.", product.Name)
		return
	}
	product.Creator = user.Username
	product, err = mydb.createProduct(product)
	if err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["product"] = product
}

func updateProduct(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	product := &Product{}
	if err := c.BindJSON(&product); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if p := mydb.getProductByName(product.Name); p != nil {
		if p.ID != product.ID {
			response["code"] = http.StatusBadRequest
			response["message"] = fmt.Sprintf("product name %s already exists", product.Name)
			return
		}
	}
	if err := mydb.updateProduct(product); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["product"] = product
}

func deleteProduct(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	if err := mydb.deleteProduct(c.GetInt("product_id")); err != nil {
		response["code"] = http.StatusInternalServerError
	}
}

/*
TODO. 判断是否在该产品线
*/
func listProductUsers(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, users := mydb.getProductUsers(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func listNotInProductUsers(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, users := mydb.getNotInProductUsers(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func addUsers2Product(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	productID := c.GetInt("product_id")
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.addUsers2Product(productID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	total, users := mydb.getProductUsers(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func removeUsersFromProduct(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	productID := c.GetInt("product_id")
	ids := struct {
		IDS []int `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.removeUsersFromProduct(productID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	total, users := mydb.getProductUsers(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["users"] = users
}

func listProductHosts(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	order := c.GetString("order")
	if len(order) == 0 {
		order = "status asc"
	}
	var (
		noGroup   bool
		err       error
		productID = c.GetInt("product_id")
	)
	noGroup, err = strconv.ParseBool(c.DefaultQuery("no_group", "false"))
	if err != nil {
		noGroup = false
	}

	total, hosts := mydb.getProductHosts(
		productID,
		noGroup,
		c.GetBool("paging"),
		c.GetString("query"),
		order,
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func listNotInProductHosts(c *gin.Context) {
	total, hosts := mydb.getNotInProductHosts(
		c.GetInt("product_id"),
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"hosts": hosts,
		"total": total})
}

func addHosts2Product(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []string `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	productID := c.GetInt("product_id")
	if err := mydb.addHosts2Product(productID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	total, hosts := mydb.getProductHosts(
		c.GetInt("product_id"),
		false,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["total"] = total
	response["hosts"] = hosts
}

func removeHostsFromProduct(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	ids := struct {
		IDS []string `json:"ids"`
	}{}
	if err := c.BindJSON(&ids); err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	productID := c.GetInt("product_id")
	if err := mydb.removeHostsFromProduct(productID, ids.IDS); err != nil {
		response["code"] = http.StatusInternalServerError
		log.Println(err)
		return
	}
	total, hosts := mydb.getProductHosts(
		c.GetInt("product_id"),
		false,
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["hosts"] = hosts
	response["total"] = total

}

func muteProductHosts(c *gin.Context) {
	ids := c.PostFormArray("ids")
	fmt.Println(ids)
	var hids []int
	c.BindJSON(&hids)
	fmt.Println(hids)
	c.JSON(http.StatusOK, gin.H{})
}

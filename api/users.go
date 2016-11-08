package main

import (
	"fmt"
	"net/http"
	"owl/common/types"
	"owl/common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	types.User
	Groups []types.UserGroup `gorm:"many2many:user_user_group;" json:"groups"`
}

func userList(c *gin.Context) {
	group_id, _ := strconv.Atoi(c.DefaultQuery("group_id", "0"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	q := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", DefaultPageSize))
	hasNull, _ := strconv.Atoi(c.DefaultQuery("hasNull", "0"))

	users := []*User{}
	db := mydb.Table("user").Order("create_at desc")
	if group_id != 0 {
		db = db.Joins("JOIN user_user_group ON user.id = user_user_group.user_id").
			Where("user_user_group.user_group_id = ?", group_id)
	}
	if status != -1 {
		db = db.Where("status = ?", status)
	}
	if hasNull != 0 {
		db = db.Where("phone = '' or mail = '' or weixin=''")
	}
	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		db = db.Where("username like ? or phone like ? or weixin like ? or mail like ?", q, q, q, q)
	}
	var cnt int
	db.Count(&cnt)
	if page != 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	db.Find(&users)
	for _, user := range users {
		mydb.Model(&user).Association("Groups").Find(&user.Groups)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "users": users, "total": cnt})
}

func userCreate(c *gin.Context) {
	user := User{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if user.Username == "" || user.Phone == "" {
		response["code"] = http.StatusBadRequest
		response["message"] = "The username and phone number are not allowed to be empty"
		return
	}
	cnt := 0
	mydb.Table("user").Where("username=?", user.Username).Count(&cnt)
	if cnt > 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = "username already exists"
		return
	}
	user.Password = utils.Md5(user.Username)
	if err := mydb.Set("gorm:save_associations", false).Save(&user).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	mydb.Model(&user).Association("Groups").Replace(user.Groups).Find(&user.Groups)

	response["code"] = http.StatusOK
	response["message"] = "user create successful"
	response["user"] = user
}

func userDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := types.User{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if mydb.Table("user").Where("id = ?", id).Find(&user).RecordNotFound() {
		response["code"] = http.StatusNotFound
		response["message"] = "The user not found"
		return
	}
	if err := mydb.Delete(&user).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return

	}
	response["code"] = http.StatusOK
	response["message"] = fmt.Sprintf("%s delete successful", user.Username)
}

func userUpdate(c *gin.Context) {
	user := User{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)

	if err := c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	currentUser := GetUser(c)
	if !currentUser.IsAdmin() && currentUser.ID != user.ID {
		response["code"] = 403
		response["message"] = "You don't have permission to access."
		return
	}
	if err := mydb.Set("gorm:save_associations", false).Table("user").
		Select("username", "role", "phone", "mail", "weixin", "status").Save(&user).
		Association("Groups").Replace(user.Groups).Error; err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	response["code"] = http.StatusOK
	response["message"] = "user update successful"
}

func changePassword(c *gin.Context) {
	request := struct {
		ID              int    `json:'id'`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}{}
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	if err := c.BindJSON(&request); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	currentUser := GetUser(c)
	if !currentUser.IsAdmin() && currentUser.ID != request.ID {
		response["code"] = 403
		response["message"] = "You don't have permission to access."
		return
	}
	if request.CurrentPassword == request.NewPassword {
		response["code"] = http.StatusNotModified
		response["message"] = "new password equal current password"
		return
	}
	user := types.User{}
	mydb.Table("user").Where("id = ?", request.ID).Find(&user)
	if user.Password != utils.Md5(request.CurrentPassword) {
		response["code"] = http.StatusNotModified
		response["message"] = "current password is incorrect"
		return
	}
	mydb.Table("user").Where("id = ?", request.ID).Update("password", utils.Md5(request.NewPassword))
	response["code"] = http.StatusOK
	response["message"] = "change password success"
}

func userInfo(c *gin.Context) {
	user := GetUser(c)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "user": user})
}

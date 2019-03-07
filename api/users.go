package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wuyingsong/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	LocalUser = "local"
	IAMUser   = "iam"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"-" db:"password"`
	DisplayName  string `json:"display_name" db:"display_name"`
	EmailAddress string `json:"email_address" db:"mail"`
	PhoneNum     string `json:"phone_number" db:"phone"`
	Wechat       string `json:"wechat" db:"wechat"`
	Role         int    `json:"role" db:"role"` // 1:admin 0:user
	Status       int    `json:"status" db:"status"`
	Type         string `json:"type" db:"type"`
	CreateAt     string `json:"created_date" db:"create_at"`
	UpdateAt     string `json:"updated_date" db:"update_at"`
}

func (u *User) isAdmin() bool {
	return u.Role == 1
}

func (u *User) isDisabled() bool {
	return u.Status == 0
}

func SyncUsers(c *gin.Context) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sync_users", config.IamURL), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("App-ID", config.AppID)
	req.Header.Add("Api-Key", config.AppKey)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	s := struct {
		Code    int    `json:"code"`
		Message []User `json:"message"`
	}{}
	if err := json.Unmarshal(buf, &s); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, users := mydb.getAllUsers(false, "", "", 0, 0)
	userMap := map[string]struct{}{}
	for _, user := range s.Message {
		userMap[user.Username] = struct{}{}
		if user.Username == "" {
			continue
		}
		user.Type = IAMUser
		user, _ := mydb.createUser(&user)
		log.Println("create user:", user)
	}
	for _, user := range users {
		if _, ok := userMap[user.Username]; ok {
			continue
		}
		if user.Type == IAMUser {
			mydb.deleteUser(user.ID)
			log.Printf("delete iam user %v", user)
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": s})
}

func listAllUsers(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	total, users := mydb.getAllUsers(
		c.GetBool("paging"),
		c.GetString("query"),
		c.GetString("order"),
		c.GetInt("offset"),
		c.GetInt("limit"),
	)
	response["code"] = http.StatusOK
	response["total"] = total
	response["users"] = users
}

/*
1. 获取用户token
2.
*/

func getUserProfile(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var user *User
	if user = mydb.getUserProfile(c.GetString("username")); user == nil {
		response["code"] = http.StatusBadRequest
		response["message"] = "user not found"
		return
	}
	response["result"] = user
}

func createUser(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		user *User
		err  error
	)
	if err = c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = err.Error()
		return
	}
	if user.Username == "" {
		response["code"] = http.StatusBadRequest
		response["message"] = "username is empty"
		return
	}
	if user := mydb.getUserProfile(user.Username); user != nil {
		response["code"] = http.StatusBadRequest
		response["message"] = "user is exsits"
		return
	}
	user.Password = utils.Md5(user.Username)
	user.Type = LocalUser
	if user, err = mydb.createUser(user); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["code"] = http.StatusOK
	response["user"] = user
}

func deleteUser(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if err := mydb.deleteUser(userID); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
}

func resetUserPassword(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var user *User
	if err := c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		log.Println(err)
		return
	}
	username := user.Username
	if user = mydb.getUserProfile(username); user == nil {
		response["code"] = http.StatusBadRequest
		return
	}
	user.Password = utils.Md5(user.Username)
	if err := mydb.updateUserProfile(user); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["user"] = user
}

func updateUserProfile(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var user User
	if err := c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		log.Println(err)
		return
	}
	if user.ID == 0 {
		response["code"] = http.StatusBadRequest
		return
	}
	currUser := mydb.getUserProfile(c.GetString("username"))
	if user.ID != currUser.ID {
		if !currUser.isAdmin() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}
	}
	user.Password = currUser.Password
	if err := mydb.updateUserProfile(&user); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["user"] = mydb.getUserProfile(user.Username)
}

func changeUserRole(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var user User
	if err := c.BindJSON(&user); err != nil {
		response["code"] = http.StatusBadRequest
		log.Println(err)
		return
	}
	if user.ID == 0 {
		response["code"] = http.StatusBadRequest
		response["message"] = "user id is Illegal"
		return
	}
	if err := mydb.setUserRole(user.ID, user.Role); err != nil {
		response["code"] = http.StatusInternalServerError
		return
	}
	response["user"] = user
}

func changeUserPassword(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
	newPassword, _ := c.GetPostForm("new_password")
	user := mydb.getUserProfile(username)
	if user == nil {
		response["code"] = http.StatusBadRequest
		return
	}
	if user.Password != password {
		response["code"] = http.StatusBadRequest
		response["message"] = "password wrong"
		return
	}
	user.Password = newPassword
	if err := mydb.updateUserProfile(user); err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = err.Error()
		return
	}
	c.SetCookie("token", "", -1, "/", "", false, true)
}

func Login(c *gin.Context) {
	response := gin.H{"code": http.StatusOK}
	defer c.JSON(http.StatusOK, response)
	var (
		tokenString string
		err         error
	)
	switch c.Request.Method {
	case "GET", "get":
		tokenString = c.DefaultQuery("token", "")
		if _, err = ValidateToken(tokenString); err != nil {
			fmt.Println("validateToken error:", err.Error())
			response["code"] = http.StatusBadRequest
			response["message"] = err.Error()
			return
		}

	case "POST", "post":
		username, _ := c.GetPostForm("username")
		password, _ := c.GetPostForm("password")
		username, password = strings.TrimSpace(username), strings.TrimSpace(password)
		user := mydb.getUserProfile(username)
		if user == nil {
			response["code"] = http.StatusUnauthorized
			response["message"] = "user not found"
			return
		}
		if user.Password != password || user.Username != username {
			response["code"] = http.StatusUnauthorized
			response["message"] = "username or password is wrong"
			return
		}
		tokenString, err = GenerateToken(username, user)
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["message"] = "generate token failed"
			return
		}
		response["token"] = tokenString
	}
	c.SetCookie("token", tokenString, 86400, "/", "", false, true)
	response["code"] = http.StatusTemporaryRedirect
	response["link"] = "/"
}

func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil || len(tokenString) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized})
		return
	}
	token, err := ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": "token is valied"})
		return
	}
	if config.AuthType == "iam" {
		if err = iamLogout(token); err != nil {
			log.Println("iam logout failed", err)
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "message": "iam logout failed"})
			return
		}
	}
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "link": "/", "message": "logout sucessfully"})
}

func iamLogout(token *jwt.Token) error {
	claims := token.Claims.(jwt.MapClaims)
	username := claims["email"].(string)
	url := fmt.Sprintf("%s/logout", strings.TrimRight(config.IamURL, "/"))
	jsonStr := fmt.Sprintf(`{"email":"%s"}`, username)
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", config.AppKey)
	req.Header.Set("App-Id", config.AppID)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", res)
	}
	return nil
}

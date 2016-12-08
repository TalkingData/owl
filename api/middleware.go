package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/appleboy/gin-jwt.v2"

	"owl/common/types"
	"owl/common/utils"
)

var authMiddleware *jwt.GinJWTMiddleware
var globalMiddleware gin.HandlerFunc

func InitMiddleware() {
	globalMiddleware = GlobalMiddleware()
	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:      "TalkingData",
		Key:        []byte(GlobalConfig.SECRET_KEY),
		Timeout:    time.Hour * 24 * 7,
		MaxRefresh: time.Hour * 24 * 7,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			var user types.User

			if err := mydb.Table("user").Where("username = ?", userId).First(&user).Error; err != nil {
				c.Set("code", http.StatusNotFound)
				c.Set("message", "user not found")
				return userId, false
			}

			if user.Status != 1 {
				c.Set("code", http.StatusUnauthorized)
				c.Set("message", "user is disable")
				return userId, false
			}

			if userId == user.Username && utils.Md5(strings.TrimSpace(password)) == user.Password {
				return userId, true
			}
			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			var user types.User

			if err := mydb.Table("user").Where("username = ?", userId).First(&user).Error; err != nil {
				c.Set("code", http.StatusNotFound)
				c.Set("message", "user not found")
				return false
			}

			if user.Role == types.USER {
				for _, handler := range AdminHandlers {
					if c.HandlerName() == handler {
						return false
					}
				}
			}
			return true
		},
		PayloadFunc: func(userId string) map[string]interface{} {
			user_info := make(map[string]interface{})
			user_info["userId"] = userId
			return user_info
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if new_code, ok := c.Get("code"); ok {
				code = new_code.(int)
				new_message, _ := c.Get("message")
				message = new_message.(string)
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
				return
			}

			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}
}

func GlobalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUserFromDataBase(c)

		//for user not found
		if user == nil {
			c.Abort()
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "user not found",
			})
			return
		}

		c.Set("user", user)

		//for user status
		if user.Status != 1 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "user is disable",
			})
			return
		}

		//for opertion journal
		operations := Operations{}
		operations.OperationResult = true
		operations.OperationTime = time.Now()
		operations.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
		rip := c.Request.Header.Get("X-Real-Ip")
		if len(rip) > 0 {
			operations.IP = rip
		}
		operations.Operator = user.Username
		operations.Content = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String())

		c.Next()

		if strings.Contains(c.Request.URL.String(), "/operations") != true {
			if c.Writer.Status() <= 599 && c.Writer.Status() >= 400 {
				operations.OperationResult = false
			}
			mydb.Create(&operations)
		}
	}
}

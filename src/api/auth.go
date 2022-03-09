package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthIAM() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tokenString string
			token       *jwt.Token
			err         error
			link        string
		)
		tokenString, _ = c.Cookie("token")
		illegal := false
		//本地校验
		if token, err = ValidateToken(tokenString); err != nil {
			illegal = true

			//IAM远程校验
		} else if err = validateTokenFromIAM(tokenString); err != nil {
			illegal = true
		}
		//校验失败
		if illegal {
			c.SetCookie("token", "", -1, "/", "", false, true)
			link = fmt.Sprintf("%s/login?app_id=%s", config.IamURL, config.AppID)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "link": link})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		if username, ok := claims["email"]; ok {
			c.Set("username", username)
			action := strings.Split(c.HandlerName(), ".")[1]
			if !JudgePermission(action, username.(string), config.AppID, config.AppKey) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
				return
			}
		}
		c.Next()
	}
}

func AuthMySQL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tokenString string
			token       *jwt.Token
			err         error
			link        string
		)
		tokenString, _ = c.Cookie("token")
		if token, err = ValidateToken(tokenString); err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "link": link, "message": err.Error()})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if username, ok := claims["email"]; ok {
			if user := mydb.getUserProfile(username.(string)); user == nil {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": "user not found"})
				c.SetCookie("token", "", -1, "/", "", false, true)
				return
			}
			c.Set("username", username)
		}
		c.Next()
	}
}

func JudgePermission(actionName string, username string, appID string, apiKey string) bool {
	url := fmt.Sprintf("%s/permission?action=%s&username=%s", strings.TrimRight(config.IamURL, "/"), actionName, username)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("App-Id", appID)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusOK:
		return true
	default:
		return false
	}
}

func validateTokenFromIAM(tokenString string) error {
	var jsonStr = []byte(fmt.Sprintf(`{"token":"%s"}`, tokenString))
	resp, err := http.Post(fmt.Sprintf("%s/validate", config.IamURL), "application/json", bytes.NewReader(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(res))
		return fmt.Errorf("%s", res)
	}
	return nil
}

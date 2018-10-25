package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Operation 操作日志的结构体
type Operation struct {
	IP       string    `json:"ip"`
	Operator string    `json:"operator"`
	Method   string    `json:"method"`
	API      string    `json:"api"`
	Body     string    `json:"body"`
	Result   int       `json:"result"`
	Time     time.Time `json:"time"`
}

// MarshalJSON 时间格式转换
func (o Operation) MarshalJSON() ([]byte, error) {
	type Alias Operation
	return json.Marshal(&struct {
		Time string `json:"time"`
		Alias
	}{
		Time:  o.Time.Format("2006-01-02 15:04:05"),
		Alias: (Alias)(o),
	})
}

// OperationRecord 操作记录中间件
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.ToLower(c.Request.Method) == "get" {
			c.Next()
			return
		}
		operation := &Operation{}
		operation.Result = 1
		operation.Time = time.Now()
		operation.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
		rip := c.Request.Header.Get("X-Real-Ip")
		if len(rip) > 0 {
			operation.IP = rip
		}
		operation.Operator = c.GetString("username")
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		operation.Method = c.Request.Method
		operation.API = c.Request.URL.String()
		operation.Body = string(body)

		c.Next()

		if strings.Contains(c.Request.URL.String(), "/operations") != true {
			if c.Writer.Status() <= 599 && c.Writer.Status() >= 400 {
				operation.Result = 0
			}
			mydb.CreateOperation(operation)
		}
	}
}

func operationList(c *gin.Context) {
	startTime := c.DefaultQuery("start_time", time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05"))
	endTime := c.DefaultQuery("end_time", time.Now().Format("2006-01-02 15:04:05"))
	query := c.GetString("query")
	where := fmt.Sprintf("(time BETWEEN '%s' AND '%s')", startTime, endTime)
	if query != "" {
		where += fmt.Sprintf("AND (ip LIKE '%%%s%%' OR method LIKE '%%%s%%' OR operator LIKE '%%%s%%' OR api LIKE '%%%s%%')", query, query, query, query)
	}
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	operations := mydb.GetOperations(where, limit)
	total := mydb.GetOperationsCount(where)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "operations": operations, "total": total})
}

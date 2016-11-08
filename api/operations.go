package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Operations struct {
	OperationTime   time.Time `json:"operation_time"`
	IP              string    `json:"ip"`
	Operator        string    `json:"operator"`
	Content         string    `json:"content"`
	OperationResult bool      `json:"operation_result"`
}

func (o Operations) MarshalJSON() ([]byte, error) {
	type Alias Operations
	return json.Marshal(&struct {
		OperationTime string `json:"operation_time"`
		Alias
	}{
		OperationTime: o.OperationTime.Format("2006-01-02 15:04:05"),
		Alias:         (Alias)(o),
	})
}

func operationsList(c *gin.Context) {
	var operations []Operations
	var fail int
	var success int
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	operation_result, err := strconv.Atoi(c.DefaultQuery("operation_result", "-1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
	}
	key := c.DefaultQuery("key", "")
	cycle := c.DefaultQuery("cycle", "now")
	year := 0
	month := 0
	day := 0
	switch cycle {
	case "all":
		year = -10
	case "now":
		day = 0
	case "1d":
		day = -1
	case "3d":
		day = -3
	case "7d":
		day = -7
	case "1m":
		month = -1
	case "3m":
		month = -3
	}
	start := c.DefaultQuery("start", time.Now().AddDate(year, month, day).Format("2006-01-02 00:00:00"))
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02 15:04:05"))

	offset := (page - 1) * page_size
	sort := "operation_time desc"

	where := fmt.Sprintf("`operation_time` BETWEEN '%s' AND '%s'", start, end)
	if key != "" {
		key = strings.TrimSpace(key)
		where += fmt.Sprintf(" AND (`operator` LIKE '%%%s%%' OR `content` LIKE '%%%s%%')", key, key)
	}

	mydb.Table("operations").Where(where).Where("operation_result = ?", 0).Count(&fail)
	mydb.Table("operations").Where(where).Where("operation_result = ?", 1).Count(&success)

	if operation_result != -1 {
		where += fmt.Sprintf(" AND `operation_result` = %d", operation_result)
	}

	mydb.Where(where).Order(sort).Offset(offset).Limit(page_size).Find(&operations)

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "operations": &operations, "success": success, "fail": fail})
}

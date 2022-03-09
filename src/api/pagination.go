package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	//DefaultPageSize  设置默认每页返回条数
	DefaultPageSize = 10
	//MaxPageSize 设置每页最大条数
	MaxPageSize = 500
)

func pagination(c *gin.Context) {
	var (
		page     int
		pageSize int
	)
	page, _ = strconv.Atoi(c.Query("page"))
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	pagingString := c.DefaultQuery("paging", "true")
	orderString := strings.TrimSpace(c.Query("order"))
	queryString := strings.TrimSpace(c.Query("query"))
	if len(orderString) > 0 {
		fields := strings.Split(orderString, "|")
		if len(fields) == 2 {
			orderKey, orderMethod := fields[0], fields[1]
			switch orderMethod {
			case "asc", "desc":
			default:
				orderMethod = "asc"
			}
			orderString = strings.Join([]string{orderKey, orderMethod}, " ")
		} else {
			orderString = ""
		}

	}
	c.Set("query", queryString)
	c.Set("order", orderString)
	c.Set("paging", true)

	if pagingString == "false" {
		c.Set("paging", false)
		return
	}
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	c.Set("offset", (page-1)*pageSize)
	c.Set("limit", pageSize)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.String(http.StatusOK, "TalkingData owl api, use /apidoc to show the detail information")
}

func apiDoc(c *gin.Context) {
	var routes string
	for _, route := range Engine.Routes() {
		routes += fmt.Sprintf("%s %s \n", route.Method, route.Path)
	}
	c.String(http.StatusOK, routes)
}

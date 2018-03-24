package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func nodesStatus(c *gin.Context) {
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(config.AlarmHealthCheckURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	u.Path = path.Join(u.Path, "nodes", "status")
	resp, err := client.Get(u.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "application/json", body)
}

func queuesStatus(c *gin.Context) {
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(config.AlarmHealthCheckURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	u.Path = path.Join(u.Path, "queues", "status")
	resp, err := client.Get(u.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "application/json", body)
}

func queuesClean(c *gin.Context) {
	id := c.Param("id")
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(config.AlarmHealthCheckURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	u.Path = path.Join(u.Path, "queues", id, "clean")
	resp, err := client.Post(u.String(), "application/json", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "application/json", body)
}

func queuesMute(c *gin.Context) {
	id := c.Param("id")
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(config.AlarmHealthCheckURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	u.Path = path.Join(u.Path, "queues", id, "mute")
	resp, err := client.Post(u.String(), "application/json", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "application/json", body)
}

func queuesUnMute(c *gin.Context) {
	id := c.Param("id")
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	u, err := url.Parse(config.AlarmHealthCheckURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	u.Path = path.Join(u.Path, "queues", id, "unmute")
	resp, err := client.Post(u.String(), "application/json", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "application/json", body)
}

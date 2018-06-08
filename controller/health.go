package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) startHttpServer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("%v", r)
			debug.PrintStack()
		}
	}()

	Engine := gin.Default()
	Engine.GET("/debug/pprof/", gin.WrapF(http.HandlerFunc(pprof.Index)))
	Engine.GET("/debug/pprof/cmdline", gin.WrapF(http.HandlerFunc(pprof.Cmdline)))
	Engine.GET("/debug/pprof/profile", gin.WrapF(http.HandlerFunc(pprof.Profile)))
	Engine.GET("/debug/pprof/symbol", gin.WrapF(http.HandlerFunc(pprof.Symbol)))
	Engine.GET("/debug/pprof/trace", gin.WrapF(http.HandlerFunc(pprof.Trace)))
	nodes := Engine.Group("/nodes")
	{
		nodes.GET("/status", nodesStatus)
	}
	queues := Engine.Group("/queues")
	{
		queues.GET("/status", queueStatus)
		queues.POST("/:id/clean", queueClean)
		queues.POST("/:id/mute", queueMute)
		queues.POST("/:id/unmute", queueUnMute)
	}

	Engine.Run(GlobalConfig.HTTP_SERVER)
}

func nodesStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"nodes": controller.nodePool.Nodes})
}

func queueStatus(c *gin.Context) {
	qs := make([]gin.H, 0)
	controller.eventQueuesMutex.RLock()
	defer controller.eventQueuesMutex.RUnlock()
	for _, product := range mydb.GetProducts() {
		if queue, ok := controller.eventQueues[product.ID]; ok {
			qs = append(qs, gin.H{"id": product.ID, "name": product.Name, "count": queue.len(), "mute": queue.mute})
		}
	}
	c.JSON(http.StatusOK, gin.H{"queues": qs})
}

func queueClean(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	if queue, ok := controller.eventQueues[i]; ok {
		queue.clean()
		c.JSON(http.StatusOK, gin.H{"message": "cleaned!"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "queue not found!"})
}

func queueMute(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	if queue, ok := controller.eventQueues[i]; ok {
		queue.mute = true
		c.JSON(http.StatusOK, gin.H{"message": "muted!"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "queue not found!"})
}

func queueUnMute(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	if queue, ok := controller.eventQueues[i]; ok {
		queue.mute = false
		c.JSON(http.StatusOK, gin.H{"message": "unmuted!"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "queue not found!"})
}

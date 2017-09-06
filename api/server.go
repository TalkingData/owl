package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"

	"owl/common/types"
)

const (
	DefaultPageSize = "10"
)

var Engine *gin.Engine
var AdminHandlers = []string{"main.userCreate", "main.userDelete"}

func InitServer() error {
	InitMiddleware()
	Engine = gin.Default()
	Engine.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	Engine.GET("/", index)
	Engine.GET("/apidoc", apiDoc)

	v1 := Engine.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)

		authorized := v1.Group("/")
		authorized.Use(authMiddleware.MiddlewareFunc())
		authorized.Use(globalMiddleware)
		{
			authorized.GET("/refresh_token", authMiddleware.RefreshHandler)

			suggest := authorized.Group("/suggest")
			{
				suggest.GET("/metric", suggestMetric)
				suggest.GET("/tagk", suggestTagk)
				suggest.GET("/tagv", suggestTagv)
				suggest.GET("/buildindex", BuildMetricAndTagIndex)
			}

			statistics := authorized.Group("/statistics")
			{
				statistics.GET("/host/info", hostInfo)
				statistics.GET("/host/status", hostStatus)
				statistics.GET("/event/count", eventsCount)
				statistics.GET("/event/status", eventsStatus)
				statistics.GET("/event/panel", eventsPanel)

			}

			host := authorized.Group("/hosts")
			{
				host.GET("", hostList)
				host.GET("/:id/strategies", strategiesByHostId)
				host.GET("/:id/plugins", pluginByHostID)
				host.GET("/:id/metrics", metricsByHostId)
				host.POST("/:id/rename", hostRename)
				host.POST("/:id/enable", hostEnable)
				host.POST("/:id/disable", hostDisable)
				host.DELETE("/:id", hostDelete)
			}

			group := authorized.Group("/groups")
			{
				group.GET("", groupList)
				group.PUT("", groupCreate)
				group.POST("", groupUpdate)
				group.DELETE("/:id", groupDelete)
				group.GET("/:id/hosts", groupHostList)
			}

			event := authorized.Group("/events")
			{
				event.GET("", eventsList)
				event.GET("/:id/detail", eventDetail)
				event.GET("/:id/process", processInfo)
				event.POST("/:id/inform", eventInform)
				event.POST("/:id/close", eventClose)
				event.DELETE("", eventDelete)
			}

			operation := authorized.Group("/operations")
			{
				operation.GET("", operationsList)
			}

			panel := authorized.Group("/panels")
			{
				panel.GET("", panelList)
				panel.PUT("", panelCreate)
				panel.POST("", panelUpdate)
				panel.DELETE("/:id", panelDelete)
				panel.GET("/:id/charts", getPanelChart)
				panel.POST("/:id/favor", panelFavor)
				panel.POST("/:id/unfavor", panelUnFavor)
			}

			chart := authorized.Group("/charts")
			{
				chart.GET("", chartList)
				chart.GET("/:id", chartDetail)
				chart.PUT("", chartCreate)
				chart.POST("", chartUpdate)
				chart.DELETE("/:id", chartDelete)
			}

			strategy := authorized.Group("/strategies")
			{
				strategy.GET("", strategiesList)
				strategy.GET("/:id/info", strategyInfo)
				strategy.PUT("", strategyCreate)
				strategy.DELETE("/:id", strategyDelete)
				strategy.POST("", strategyUpdate)
				strategy.POST("/:id/switch", strategySwitch)
			}

			user := authorized.Group("/users")
			{
				user.GET("", userList)
				user.POST("/changepassword", changePassword)
				user.GET("/info", userInfo)
				user.POST("", userUpdate)
				user.PUT("", userCreate)
				user.DELETE("/:id", userDelete)
			}

			ugroup := authorized.Group("/usergroups")
			{
				ugroup.GET("", userGroupList)
				ugroup.PUT("", userGroupCreate)
				ugroup.POST("", userGroupUpdate)
				ugroup.DELETE("/:id", userGroupDelete)
			}

			plugin := authorized.Group("/plugins")
			{
				plugin.GET("", pluginList)
				plugin.PUT("", pluginCreate)
				plugin.POST("", pluginUpdate)
				plugin.DELETE("/:id", pluginDelete)
			}
			query := authorized.Group("/query")
			{
				query.GET("", data)
			}
		}

	}
	return Engine.Run(GlobalConfig.HTTP_BIND)
}

func GetUserFromDataBase(c *gin.Context) *types.User {
	if userID, ok := c.Get("userID"); ok {
		var user types.User
		if err := mydb.Table("user").Where("username = ?", userID).First(&user).Error; err != nil {
			return nil
		}
		return &user
	}
	return nil
}

func GetUser(c *gin.Context) *types.User {
	if user, ok := c.Get("user"); ok {
		return user.(*types.User)
	}
	return nil
}

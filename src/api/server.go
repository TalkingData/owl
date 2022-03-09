package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

var Engine *gin.Engine

func InitServer() error {
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
	v1 := Engine.Group("/api/v1", pagination)
	{
		v1.GET("/login", Login)
		v1.POST("/login", Login)
		v1.GET("/sync_user", SyncUsers)

		switch config.AuthType {
		case "mysql":
			v1.Use(AuthMySQL())
		case "iam":
			v1.Use(AuthIAM())
		}

		//操作日志
		v1.Use(OperationRecord())
		v1.GET("/query", queryTimeSeriesData)
		v1.GET("/suggest/metrics", suggestMetrics)
		v1.GET("/suggest/tags", suggestMetricTagSet)
		v1.POST("/logout", Logout)

		plugin := v1.Group("/plugins")
		{
			plugin.GET("", listPlugins)
			plugin.POST("", verifyAdminPermission, createPlugin)
			plugin.PUT("", verifyAdminPermission, updatePlugin)
			plugin.DELETE("/:plugin_id", verifyAdminPermission, deletePlugin)

			// plugin.PUT("/:plugin_id/host_groups/add", addHostGroups2Plugin)
			// plugin.PUT("/:plugin_id/host_groups/remove", removeHostGroupsFromPlugin)
			// plugin.GET("/:plugin_id/host_groups", listPluginHostGroups)
			// plugin.GET("/:plugin_id/host_groups/not_in", listNotInPluginHostGroups)
		}

		// user interface
		user := v1.Group("/users")
		{
			user.GET("", verifyAdminPermission, listAllUsers)
			if config.AuthType == "mysql" {
				//创建用户
				user.POST("", verifyAdminPermission, createUser)
				//修改密码
				user.POST("changepassword", changeUserPassword)
				user.POST("resetpassword", verifyAdminPermission, resetUserPassword)
			}
			// 获取用户所在的产品列表
			user.GET("/products", listUserProducts)
			//获取用户信息
			user.GET("/profile", getUserProfile)
			//更新用户信息
			user.PUT("/profile", updateUserProfile)
			//修改角色
			user.PUT("/change_role", verifyAdminPermission, changeUserRole)
			//删除用户
			user.DELETE("/:user_id", verifyAdminPermission, deleteUser)
		}

		hosts := v1.Group("/hosts")
		{
			//获取所有主机列表
			hosts.GET("", verifyAdminPermission, listAllHosts)
			// 获取主机的metric列表
			hosts.GET("/:host_id/metrics", listHostMetrics)
			hosts.DELETE("/:host_id/metrics", deleteHostMetrics)

			//获取主机 metric 前缀
			hosts.GET("/:host_id/apps", listHostApps)

			hosts.GET("/:host_id", getHost)

			// 添加 plugin 到主机
			hosts.POST("/:host_id/plugins", createHostPlugin)

			hosts.PUT("/:host_id/plugins", updateHostPlugin)

			hosts.PUT("/:host_id/mute", muteHost)
			hosts.PUT("/:host_id/unmute", unmuteHost)

			hosts.GET("/:host_id/plugins", listHostPlugins)

			hosts.DELETE("/:host_id/plugins/:plugin_id", deleteHostPlugin)

			hosts.DELETE("/:host_id", verifyAdminPermission, deleteHost)
		}

		product := v1.Group("/products")
		{
			product.POST("", verifyAdminPermission, createProduct)
			product.PUT("", verifyAdminPermission, updateProduct)
			product.GET("", verifyAdminPermission, listAllProducts)
		}

		// 报警策略模板
		template := v1.Group("/templates")
		{
			template.GET("", templateList)
			template.GET("/:template_id", templateGet)
			template.POST("", verifyAdminPermission, templateCreate)
			template.PUT("", verifyAdminPermission, templateUpdate)
			template.DELETE("", verifyAdminPermission, templateDelete)
		}

		// 报警执行脚本
		script := v1.Group("/scripts")
		{
			script.GET("", scriptList)
			script.POST("", verifyAdminPermission, scriptCreate)
			script.PUT("", verifyAdminPermission, scriptUpdate)
			script.DELETE("", verifyAdminPermission, scriptDelete)
		}

		// 操作日志
		operation := v1.Group("/operations", verifyAdminPermission)
		{
			operation.GET("", operationList)
		}

		// 健康检查
		health := v1.Group("/health", verifyAdminPermission)
		{
			health.GET("/nodes/status", nodesStatus)
			health.GET("/queues/status", queuesStatus)
			health.POST("/queues/:id/clean", queuesClean)
			health.POST("/queues/:id/mute", queuesMute)
			health.POST("/queues/:id/unmute", queuesUnMute)
		}

		productDetail := product.Group("/:product_id", verifyAndInjectProductID)
		{
			productDetail.DELETE("", verifyAdminPermission, deleteProduct)

			productPanel := productDetail.Group("/panels")
			{
				productPanel.GET("", listProductPanel)
				productPanel.POST("", createProductPanel)
				productPanel.PUT("", updateProductPanel)
				productPanel.DELETE("/:panel_id", deleteProductPanel)
				panelChart := productPanel.Group("/:panel_id/charts", verifyAndInjectPanelID)
				{
					panelChart.GET("", listPanelChats)
					panelChart.POST("", createPanelChart)
					panelChart.PUT("", updatePanelChart)
					panelChart.DELETE("/:chart_id", deletePanelChart)
				}
			}

			productHost := productDetail.Group("/hosts")
			{
				// 获取产品线下的主机列表
				productHost.GET("", listProductHosts)
				// 向产品线中分配主机
				productHost.PUT("/add", verifyAdminPermission, addHosts2Product)
				// 从产品线中移除主机
				productHost.PUT("/remove", verifyAdminPermission, removeHostsFromProduct)
				// 获取不在产品线中的主机
				productHost.GET("/not_in", verifyAdminPermission, listNotInProductHosts)
			}
			productUser := productDetail.Group("/users")
			{
				// 获取产品线下的用户
				productUser.GET("", listProductUsers)
				// 向产品线中添加用户
				productUser.PUT("/add", addUsers2Product)
				// 从产品线中移除用户
				productUser.PUT("/remove", removeUsersFromProduct)
				// 获取不在产品线中的用户
				productUser.GET("/not_in", listNotInProductUsers)
			}
			productHostGroup := productDetail.Group("/host_groups")
			{
				// 获取产品线下的主机组
				productHostGroup.GET("", listProductHostGroups)
				// 更新产品线下的主机组
				productHostGroup.PUT("", updateProductHostGroup)
				// 创建产品线主机组
				productHostGroup.POST("", createProductHostGroup)
				//删除产品线主机组
				productHostGroup.DELETE("/:host_group_id", verifyAndInjectHostGroupID, deleteProductHostGroup)

				productHostGroupPlugin := productHostGroup.Group("/:host_group_id/plugins", verifyAndInjectHostGroupID)
				{
					productHostGroupPlugin.GET("", listHostGroupPlugins)
					productHostGroupPlugin.PUT("", updateHostGroupPlugin)
					productHostGroupPlugin.POST("", createHostGroupPlugin)
					productHostGroupPlugin.DELETE("/:plugin_id", deleteHostGroupPlugin)
				}

				productHostGroupHost := productHostGroup.Group("/:host_group_id/hosts", verifyAndInjectHostGroupID)
				{
					// 获取主机组下的主机列表
					productHostGroupHost.GET("", listProductHostGroupHosts)
					// 向主机组中添加主机
					productHostGroupHost.PUT("/add", addHosts2ProductHostGroup)
					// 从主机组中移除主机
					productHostGroupHost.PUT("/remove", removeHostsFromProductHostGroup)
					// 获取不在主机组内的主机
					productHostGroupHost.GET("/not_in", listNotInProductHostGroupHosts)
				}
			}
			productUserGroup := productDetail.Group("/user_groups")
			{
				// 获取产品线下的用户组
				productUserGroup.GET("", listProductUserGroups)
				// 更新产品线用户组信息
				productUserGroup.PUT("", updateProductUserGroup)
				// 创建产品线用户组
				productUserGroup.POST("", createProductUserGroup)
				//删除产品线用户组
				productUserGroup.DELETE("/:user_group_id", deleteProductUserGroup)

				productUserGroupUser := productUserGroup.Group("/:user_group_id/users", verifyAndInjectUserGroupID)
				{
					// 获取用户组下的用户列表
					productUserGroupUser.GET("", listProductUserGroupUsers)
					// 向用户组中添加用户
					productUserGroupUser.PUT("/add", addUsers2ProductUserGroup)
					// 从用户组中移除用户
					productUserGroupUser.PUT("/remove", removeUsersFromProductUserGroup)
					// 获取不在用户组内的用户
					productUserGroupUser.GET("/not_in", listNotInProductUserGroupUsers)
				}
			}
			strategy := productDetail.Group("/strategies")
			{
				// 获取报警策略列表
				strategy.GET("", strategyList)
				// 创建报警策略
				strategy.POST("", strategyCreate)
				// 修改报警策略
				strategy.PUT("", strategyUpdate)
				// 删除报警策略
				strategy.DELETE("", strategyDelete)
				// 获取单个报警策略详细信息
				strategy.GET("/info/:strategy_id", strategyInfo)
				// 禁用报警策略
				strategy.PUT("/switch", strategySwitch)
				// 获取主机组下的策略列表
				strategy.GET("/list/:host_group_id", strategyHostGroup)

			}
			event := productDetail.Group("/events")
			{
				// 获取报警事件列表
				event.GET("", eventsList)
				// 获取失败的报警事件列表
				event.GET("/failed", eventsFailed)
				// 知悉报警事
				event.PUT("/aware", eventAware)
				// 获取某个报警事件的历史处理记录
				event.GET("/process/:event_id", eventProcessRecord)
				// 获取某个报警事件的详细信息
				event.GET("/detail/:event_id", eventDetail)
			}
		}
	}

	return Engine.Run(config.HTTPBind)
}

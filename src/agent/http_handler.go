package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	httpHandler "owl/common/http_handler"
	"owl/common/logger"
	"owl/dto"
	"time"
)

// newHttpHandler 创建http handler
func (a *agent) newHttpHandler() http.Handler {
	h := gin.New()
	h.Use(gin.Recovery(), httpHandler.GinAccessLogMiddleware(a.logger))

	h.POST("/ts_data", a.tsDataHandler(true))
	h.POST("/ts_data/raw", a.tsDataHandler(false))

	h.NoRoute(httpHandler.PageNotFoundHandler)

	return h
}

func (a *agent) tsDataHandler(fillAgentInfo bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tsdArr := dto.TsDataArray{}

		// 序列化request body
		if err := c.ShouldBindJSON(&tsdArr); err != nil {
			a.logger.WarnWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while calling c.ShouldBindJSON.")
			httpHandler.WriteResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		currTs := time.Now().Unix()
		// 验证所有发来的数据
		for idx, tsd := range tsdArr {
			if err := tsd.Validate(); err != nil {
				msg := fmt.Sprintf(
					"Ts data validate failed at Index=%v, Metric=%v, Error=%v",
					idx,
					tsd.Metric,
					err.Error(),
				)
				httpHandler.WriteResponse(c, http.StatusBadRequest, msg)
				a.logger.WarnWithFields(logger.Fields{
					"index":             idx,
					"ts_data_metric":    tsd.Metric,
					"ts_data_data_type": tsd.DataType,
					"ts_data_value":     tsd.Value,
					"ts_data_timestamp": tsd.Timestamp,
					"ts_data_cycle":     tsd.Cycle,
					"ts_data_tags":      tsd.Tags,
					"error":             err,
				}, "An error occurred while calling c.ShouldBindJSON")
				return
			}

			// 对于没有设置时间戳的metric，自动设置为当前时间戳
			if tsd.Timestamp < 1 {
				tsd.Timestamp = currTs
			}
		}

		// 预处理数据
		a.preprocessTsData(tsdArr, fillAgentInfo)
		httpHandler.WriteResponse(c, 0, "success")
	}
}

// newMetricHttpHandler 创建metric http handler
func (a *agent) newMetricHttpHandler() http.Handler {
	h := gin.New()
	h.Use(gin.Recovery())

	h.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h.NoRoute(httpHandler.PageNotFoundHandler)

	return h
}

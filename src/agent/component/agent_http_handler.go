package component

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/logger"
	"owl/dto"
	"time"
)

func (agent *agent) newHttpHandler() http.Handler {
	h := gin.New()
	h.Use(gin.Recovery(), agent.ginAccessLogMiddleware())

	h.POST("/ts_data", agent.tsDataHandler(true))
	h.POST("/ts_data/raw", agent.tsDataHandler(false))

	h.NoRoute(pageNotFoundHandler)

	return h
}

func (agent *agent) tsDataHandler(fillAgentInfo bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tsdArr := dto.TsDataArray{}

		// 序列化request body
		if err := c.ShouldBindJSON(&tsdArr); err != nil {
			agent.logger.WarnWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while c.ShouldBindJSON in agent.tsDataHandler")
			writeResponse(c, http.StatusBadRequest, err.Error())
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
				writeResponse(c, http.StatusBadRequest, msg)
				agent.logger.WarnWithFields(logger.Fields{
					"index":             idx,
					"ts_data_metric":    tsd.Metric,
					"ts_data_data_type": tsd.DataType,
					"ts_data_value":     tsd.Value,
					"ts_data_timestamp": tsd.Timestamp,
					"ts_data_cycle":     tsd.Cycle,
					"ts_data_tags":      tsd.Tags,
					"error":             err,
				}, "An error occurred while c.ShouldBindJSON in agent.tsDataHandler")
				return
			}

			// 对于没有设置时间戳的metric，自动设置为当前时间戳
			if tsd.Timestamp < 1 {
				tsd.Timestamp = currTs
			}
		}

		// 发送数据
		agent.sendTsDataArray(tsdArr, fillAgentInfo)
		writeResponse(c, 0, "success")
	}
}

func (agent *agent) ginAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent.logger.InfoWithFields(logger.Fields{
			"request_uri":    c.Request.RequestURI,
			"request_method": c.Request.Method,
			"user_agent":     c.Request.UserAgent(),
		}, "Got agent http server request.")
		c.Next()
	}
}

func pageNotFoundHandler(c *gin.Context) {
	writeResponse(c, http.StatusNotFound, "Page not found.")
}

func writeResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
	})
}

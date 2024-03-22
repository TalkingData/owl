package http_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"owl/common/logger"
)

func WriteResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
	})
}

func PageNotFoundHandler(c *gin.Context) {
	WriteResponse(c, http.StatusNotFound, "Page not found.")
}

func GinAccessLogMiddleware(lg *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.InfoWithFields(logger.Fields{
			"request_uri":    c.Request.RequestURI,
			"request_method": c.Request.Method,
			"user_agent":     c.Request.UserAgent(),
		}, "Got http server request.")
		c.Next()
	}
}

package prom

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	httpHandler "owl/common/http_handler"
)

type Prom interface {
	// GrpcInterceptor 返回一个grpc拦截器
	GrpcInterceptor(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler,
	) (interface{}, error)

	// GoMicroHandlerWrapper 返回一个go-micro的handler wrapper
	GoMicroHandlerWrapper() server.HandlerWrapper

	// MetricServerStart 启动MetricServer
	MetricServerStart() error
	// MetricServerStop 关闭MetricServer
	MetricServerStop(context.Context)
}

type prom struct {
	RequestCounter *prometheus.CounterVec

	metricServer *http.Server
}

func NewProm(metricListen string) Prom {
	reqCt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "The total number of requests.",
		},
		[]string{"method", "status"},
	)
	prometheus.MustRegister(reqCt)

	reqCt.WithLabelValues("_all", "total").Add(0)
	reqCt.WithLabelValues("_all", "error").Add(0)
	reqCt.WithLabelValues("_all", "success").Add(0)

	h := gin.New()
	h.Use(gin.Recovery())

	h.GET("/metrics", gin.WrapH(promhttp.Handler()))
	h.NoRoute(httpHandler.PageNotFoundHandler)

	return &prom{
		RequestCounter: reqCt,

		metricServer: &http.Server{
			Addr:    metricListen,
			Handler: h,
		},
	}
}

package prom

import (
	"context"
	"google.golang.org/grpc"
)

// GrpcInterceptor 返回一个grpc拦截器
func (p *prom) GrpcInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler,
) (interface{}, error) {
	p.RequestCounter.WithLabelValues(info.FullMethod, "total").Inc()
	p.RequestCounter.WithLabelValues("_all", "total").Inc()

	iface, err := h(ctx, req)
	if err != nil {
		p.RequestCounter.WithLabelValues(info.FullMethod, "error").Inc()
		p.RequestCounter.WithLabelValues("_all", "error").Inc()
		return iface, err
	}

	p.RequestCounter.WithLabelValues(info.FullMethod, "success").Inc()
	p.RequestCounter.WithLabelValues("_all", "success").Inc()
	return iface, err
}

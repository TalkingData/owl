package prom

import (
	"context"
	"github.com/micro/go-micro/v2/server"
)

// GoMicroHandlerWrapper 返回一个go-micro的handler wrapper
func (p *prom) GoMicroHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			p.RequestCounter.WithLabelValues(req.Method(), "total").Inc()
			p.RequestCounter.WithLabelValues("_all", "total").Inc()

			if err := h(ctx, req, rsp); err != nil {
				p.RequestCounter.WithLabelValues(req.Method(), "error").Inc()
				p.RequestCounter.WithLabelValues("_all", "error").Inc()
				return err
			}

			p.RequestCounter.WithLabelValues(req.Method(), "success").Inc()
			p.RequestCounter.WithLabelValues("_all", "success").Inc()
			return nil
		}
	}
}

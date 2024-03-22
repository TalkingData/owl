package prom

import (
	"context"
	"net/http"
)

func (p *prom) MetricServerStart() error {
	if err := p.metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (p *prom) MetricServerStop(ctx context.Context) {
	_ = p.metricServer.Shutdown(ctx)
}

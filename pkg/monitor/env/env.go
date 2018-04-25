package env

import (
	"github.com/spacelavr/monitor/pkg/monitor/metrics"
)

type env struct {
	metrics *metrics.Metrics
}

var (
	e = &env{}
)

// SetMetrics set metrics to env
func SetMetrics(m *metrics.Metrics) {
	e.metrics = m
}

// GetMetrics get metrics from env
func GetMetrics() *metrics.Metrics {
	return e.metrics
}

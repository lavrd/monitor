package env

import (
	"github.com/spacelavr/monitor/pkg/monitor/cri"
	"github.com/spacelavr/monitor/pkg/monitor/metrics"
)

type env struct {
	cri     *cri.Cri
	metrics *metrics.Metrics
}

var (
	e = &env{}
)

func SetMetrics(m *metrics.Metrics) {
	e.metrics = m
}

func GetMetrics() *metrics.Metrics {
	return e.metrics
}

func SetCri(cri *cri.Cri) {
	e.cri = cri
}

func GetCri() *cri.Cri {
	return e.cri
}

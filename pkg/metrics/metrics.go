package metrics

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"monitor/pkg/docker"
)

// Metrics
type Metrics struct {
	*metricsMap
	docker    *docker.Docker
	cInterval time.Duration
	mInterval time.Duration
}

type metricsMap struct {
	sync.RWMutex
	metrics map[string]*docker.ContainerStats
}

// Info
type Info struct {
	Metrics []*docker.ContainerStats `json:"metrics,omitempty"`
	Alert   string                   `json:"alert,omitempty"`
}

// New returns new metrics
func New(d *docker.Docker, containersInterval, metricsInterval time.Duration) *Metrics {
	return &Metrics{
		docker: d,
		metricsMap: &metricsMap{
			metrics: make(map[string]*docker.ContainerStats),
		},
		cInterval: containersInterval,
		mInterval: metricsInterval,
	}
}

// Collect collect metrics and check for new containers
func (m *Metrics) Collect() error {
	for range time.Tick(m.cInterval) {
		containers, err := m.docker.ContainerList()
		if err != nil {
			return err
		}

		for _, container := range containers {
			if _, ok := m.metrics[container.Names[0][1:]]; !ok {
				go func() {
					if err := m.collect(container.Names[0][1:]); err != nil {
						log.Fatal().Err(err).Msg("collection metrics error")
					}
				}()
				log.Info().Msgf("new container %s", container.Names[0][1:])
			}
		}
	}

	return nil
}

// Info returns info about containers metrics
func (m *Metrics) Info(id string) *Info {
	if len(m.metrics) == 0 {
		return &Info{
			Alert: "no running containers",
		}
	}
	if metrics := m.accumulate(m.parse(id)); metrics == nil {
		return &Info{
			Alert: "these containers are not running",
		}
	} else {
		return &Info{
			Metrics: metrics,
		}
	}
}

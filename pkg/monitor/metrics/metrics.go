package metrics

import (
	"sync"
	"time"

	"monitor/pkg/cri"
	"monitor/pkg/types"
	"monitor/pkg/utils/log"

	"github.com/spf13/viper"
)

// Metrics
type Metrics struct {
	*metricsMap
	Cri        *docker.Cri
	cInterval  time.Duration
	cmInterval time.Duration
}

type metricsMap struct {
	sync.RWMutex
	metrics map[string]*docker.ContainerStats
}

// Public
type Public struct {
	Metrics []*docker.ContainerStats `json:"metrics,omitempty"`
	Alert   string                   `json:"alert,omitempty"`
}

// New returns new metrics
func New() (*Metrics, error) {
	r, err := docker.New()
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Cri: r,
		metricsMap: &metricsMap{
			metrics: make(map[string]*docker.ContainerStats),
		},
		cInterval:  time.Second * time.Duration(viper.GetInt(types.FCInterval)),
		cmInterval: time.Second * time.Duration(viper.GetInt(types.FCMInterval)),
	}, nil
}

// Collect collect metrics and check for new containers
func (m *Metrics) Collect() error {
	for range time.Tick(m.cInterval) {
		containers, err := m.Cri.ContainerList()
		if err != nil {
			return err
		}

		for _, container := range containers {
			if _, ok := m.metrics[container.Names[0][1:]]; !ok {
				go func() {
					if err := m.collect(container.Names[0][1:]); err != nil {
						log.Fatal(err)
					}
				}()
				log.Debug("new container `", container.Names[0][1:], "`")
			}
		}
	}

	return nil
}

// Public returns public info about containers metrics
func (m *Metrics) Public(id string) *Public {
	if len(m.metrics) == 0 {
		return &Public{
			Alert: "no running containers",
		}
	}

	if metrics := m.accumulate(m.parse(id)); metrics == nil {
		return &Public{
			Alert: "these containers are not running",
		}
	} else {
		return &Public{
			Metrics: metrics,
		}
	}
}

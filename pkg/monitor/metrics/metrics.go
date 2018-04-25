package metrics

import (
	"sync"
	"time"

	"github.com/spacelavr/monitor/pkg/cri"
	"github.com/spacelavr/monitor/pkg/types"
	"github.com/spacelavr/monitor/pkg/utils/log"
	"github.com/spf13/viper"
)

// Metrics
type Metrics struct {
	Cri        *cri.Cri
	*metricsMap
	cInterval  time.Duration
	cmInterval time.Duration
}

type metricsMap struct {
	sync.RWMutex
	metrics map[string]*cri.ContainerStats
}

// Public
type Public struct {
	Metrics []*cri.ContainerStats `json:"metrics,omitempty"`
	Alert   string                `json:"alert,omitempty"`
}

// New returns new metrics
func New() *Metrics {
	return &Metrics{
		Cri: cri.New(),
		metricsMap: &metricsMap{
			metrics: make(map[string]*cri.ContainerStats),
		},
		cInterval:  time.Second * time.Duration(viper.GetInt(types.FCInterval)),
		cmInterval: time.Second * time.Duration(viper.GetInt(types.FCMInterval)),
	}
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

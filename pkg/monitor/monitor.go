package monitor

import (
	"time"

	"github.com/rs/zerolog/log"

	"monitor/pkg/docker"
	"monitor/pkg/metrics"
	"monitor/pkg/monitor/api"
)

// Daemon start monitor daemon
func Daemon(metricsInterval, containersInterval time.Duration, port int) error {
	log.Debug().Msg("starting monitor daemon")

	d, err := docker.New()
	if err != nil {
		return err
	}
	defer d.Close()

	m := metrics.New(d, containersInterval, metricsInterval)

	a := api.New(metricsInterval, m, port)

	go func() {
		if err := m.Collect(); err != nil {
			log.Fatal().Err(err).Msg("collect metrics error")
		}
	}()

	go func() {
		defer a.Stop()
		if err := a.Start(); err != nil {
			log.Fatal().Err(err).Msg("start http server error")
		}
	}()

	return nil
}

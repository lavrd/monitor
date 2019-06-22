package monitor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"monitor/pkg/monitor/metrics"

	"monitor/pkg/monitor/api"
)

// Daemon start monitor daemon
func Daemon(metricsInterval, containersInterval time.Duration, port int) error {
	log.Debug().Msg("starting monitor daemon")

	m, err := metrics.New()
	if err != nil {
		return err
	}
	defer m.Cri

	env.SetMetrics(m)

	go func() {
		if err := m.Collect(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		srv := &http.Server{
			Handler: api.Router(),
			Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

}

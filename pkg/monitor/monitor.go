package monitor

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"monitor/pkg/monitor/env"
	"monitor/pkg/monitor/metrics"
	"monitor/pkg/monitor/router"

	"github.com/spf13/viper"
)

// Daemon start monitor daemon
func Daemon() error {
	log.Debug().Msg("starting monitor daemon")

	m, err := metrics.New()
	if err != nil {
		return err
	}
	defer m.Cri.Close()

	env.SetMetrics(m)

	go func() {
		if err := m.Collect(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		srv := &http.Server{
			Handler: router.Router(),
			Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

}

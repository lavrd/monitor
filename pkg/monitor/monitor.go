package monitor

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monitor/pkg/monitor/env"
	"monitor/pkg/monitor/metrics"
	"monitor/pkg/monitor/router"
	"monitor/pkg/utils/log"

	"github.com/spf13/viper"
)

// Daemon start monitor daemon
func Daemon() {
	log.Debug("start monitor daemon")

	var (
		sig = make(chan os.Signal)
	)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	m, err := metrics.New()
	if err != nil {
		log.Fatal(err)
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

	<-sig
	log.Debug("handle SIGINT and SIGTERM")
}

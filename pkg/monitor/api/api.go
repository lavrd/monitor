package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"monitor/pkg/metrics"
)

type API struct {
	metricsInterval time.Duration
	metrics         *metrics.Metrics
	srv             *http.Server
}

func New(metricsInterval time.Duration, metrics *metrics.Metrics, port int) *API {
	api := &API{
		metrics:         metrics,
		metricsInterval: metricsInterval,
	}
	api.srv = &http.Server{
		Handler: api.router(),
		Addr:    fmt.Sprintf(":%d", port),
	}
	return api
}

func (a *API) Start() error {
	log.Info().Msgf("start http server on %s", a.srv.Addr)
	return a.srv.ListenAndServe()
}

func (a *API) Stop() {
	if err := a.srv.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("shutdown server error")
	}
	log.Info().Msg("server shutdown")
}

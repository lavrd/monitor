package api

import (
	"html/template"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

// execute execute template
func (a *API) execute(path string, w http.ResponseWriter) {
	html, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal().Err(err).Msg("parse files error")
	}

	tpl := template.Must(html, err)

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("execute template error")
	}
}

// DashboardH
func (a *API) DashboardH(w http.ResponseWriter, _ *http.Request) {
	a.execute("./dashboard/index.html", w)
}

// P404H
func (a *API) P404H(w http.ResponseWriter, _ *http.Request) {
	a.execute("./dashboard/404.html", w)
}

// MetricsH
func (a *API) MetricsH(conn *websocket.Conn) {
	ids := make([]byte, 512)

	n, err := conn.Read(ids)
	if err != nil {
		log.Error().Err(err).Msg("read ids from websocket error")
		return
	}

	for range time.Tick(a.metricsInterval) {
		info := a.metrics.Info(string(ids[:n]))

		if err := websocket.JSON.Send(conn, info); err != nil {
			log.Error().Err(err).Msg("send json error")
			return
		}
	}
}

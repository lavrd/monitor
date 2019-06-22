package api

import (
	"net/http"

	"monitor/pkg/monitor/router/handlers"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

// Router returns router configuration
func Router() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dashboard/static/"))))

	r.HandleFunc("/dashboard", handlers.DashboardH)
	r.Handle("/metrics", websocket.Handler(handlers.MetricsH))
	r.NotFoundHandler = http.HandlerFunc(handlers.P404H)

	return r
}

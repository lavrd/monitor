package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spacelavr/monitor/pkg/monitor/router/handlers"
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

// todo продумтаь логи

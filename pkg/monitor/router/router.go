package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spacelavr/monitor/pkg/monitor/metrics"
)

// Router returns router configuration
func Router() http.Handler {
	r := mux.NewRouter()

	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dashboard/static/"))))
	r.PathPrefix("/dashboard/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dashboard/static/"))))

	r.HandleFunc("/dashboard", dashboard)
	r.NotFoundHandler = http.HandlerFunc(p404)

	r.HandleFunc("/api/metrics/{id}", getMetrics)
	r.HandleFunc("/api/status", status)
	r.HandleFunc("/api/stopped", getStopped)
	r.HandleFunc("/api/launched", getLaunched)
	r.HandleFunc("/api/logs/{id}", getLogs)

	return handlers.LoggingHandler(os.Stdout, r)
}

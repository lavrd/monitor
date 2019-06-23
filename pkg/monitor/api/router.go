package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

// router returns router configuration
func (a *API) router() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dashboard/static/"))))

	r.HandleFunc("/dashboard", a.DashboardH)
	r.Handle("/metrics", websocket.Handler(a.MetricsH))
	r.NotFoundHandler = http.HandlerFunc(a.P404H)

	return r
}

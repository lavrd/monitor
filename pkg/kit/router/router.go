package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spacelavr/dlm/pkg/kit/metrics"
)

// Router returns router configuration
func Router() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dashboard/static/"))))
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

// get API status
func status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// dashboard page
func dashboard(w http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("./dashboard/index.html")
	tpl := template.Must(html, err)
	tpl.Execute(w, nil)
}

// 404 page
func p404(w http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("./dashboard/404.html")
	tpl := template.Must(html, err)
	tpl.Execute(w, nil)
}

// get container metrics
func getMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics.Get().Get(mux.Vars(r)["id"])); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// get stopped containers
func getStopped(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics.Get().GetStoppedContainers()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// get launched containers
func getLaunched(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics.Get().GetLaunchedContainers()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// get container logs
func getLogs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics.GetContainerLogs(mux.Vars(r)["id"])); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

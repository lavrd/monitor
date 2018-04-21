package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// todo write response 200 func for all handlers

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

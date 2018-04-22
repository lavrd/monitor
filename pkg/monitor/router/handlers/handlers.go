package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spacelavr/monitor/pkg/log"
	"github.com/spacelavr/monitor/pkg/monitor/env"
)

func DashboardH(w http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("./dashboard/index.html")
	if err != nil {
		log.Fatal(err)
	}

	tpl := template.Must(html, err)

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func P404H(w http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("./dashboard/404.html")
	if err != nil {
		log.Fatal(err)
	}

	tpl := template.Must(html, err)

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func MetricsH(w http.ResponseWriter, r *http.Request) {
	var (
		id      = mux.Vars(r)["id"]
		m       = env.GetMetrics()
		metrics = m.Get(id)
	)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Fatal(err)
	}
}

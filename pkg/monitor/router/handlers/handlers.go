package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/spacelavr/monitor/pkg/monitor/env"
	"github.com/spacelavr/monitor/pkg/types"
	"github.com/spacelavr/monitor/pkg/utils/log"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
)

// execute execute template and render to client
func execute(path string, w http.ResponseWriter) {
	html, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	tpl := template.Must(html, err)

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// DashboardH
func DashboardH(w http.ResponseWriter, _ *http.Request) {
	execute("./dashboard/index.html", w)
}

// P404H
func P404H(w http.ResponseWriter, _ *http.Request) {
	execute("./dashboard/404.html", w)
}

// MetricsH
func MetricsH(ws *websocket.Conn) {
	var (
		m   = env.GetMetrics()
		ids = make([]byte, 512)
	)

	n, err := ws.Read(ids)
	if err != nil {
		return
	}

	duration := time.Second * time.Duration(viper.GetInt(types.FCMInterval))

	for range time.Tick(duration) {
		metrics := m.Public(string(ids[:n]))
		if err := websocket.JSON.Send(ws, metrics); err != nil {
			return
		}
	}
}

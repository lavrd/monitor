package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/spacelavr/monitor/pkg/monitor/env"
	"github.com/spacelavr/monitor/pkg/utils/log"
	"golang.org/x/net/websocket"
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

func MetricsH(ws *websocket.Conn) {
	var (
		m   = env.GetMetrics()
		ids = make([]byte, 512)
	)

	n, err := ws.Read(ids)
	if err != nil {
		return
	}

	fmt.Println(string(ids[:n]))

	for range time.Tick(time.Second * 1) {
		metrics := m.Get(string(ids[:n]))

		fmt.Println(metrics)

		if err := websocket.JSON.Send(ws, metrics); err != nil {
			return
		}
	}
}

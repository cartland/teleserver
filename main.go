package main

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
)

var startTime = time.Now()

type Metric struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func getSpeed() Metric {
	t := time.Since(startTime)
	return Metric{
		Type:  "speed",
		Value: 20 + 2*math.Cos(t.Seconds()),
	}
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	e := json.NewEncoder(ws)
	for {
		e.Encode(getSpeed())
		time.Sleep(250 * time.Millisecond)
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/ws", websocket.Handler(EchoServer))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", r)
}

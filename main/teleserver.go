package main

import (
	"net/http"

	"bitbucket.org/stvnrhodes/teleserver"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/ws", websocket.Handler(teleserver.MetricsServer))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", r)
}

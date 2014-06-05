package main

import (
	"flag"
	"fmt"
	"net/http"

	"bitbucket.org/stvnrhodes/teleserver/lib"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.Int("port", 8080, "Port for the webserver")
	flag.Parse()

	r := mux.NewRouter()
	r.Handle("/ws", websocket.Handler(teleserver.MetricsServer))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

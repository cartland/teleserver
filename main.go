package main

import (
	"flag"
	"fmt"
	"net/http"

	"bitbucket.org/stvnrhodes/teleserver/lib"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stvnrhodes/broadcaster"
)

func main() {
	port := flag.Int("port", 8080, "Port for the webserver")
	flag.Parse()

	b := broadcaster.New()
	go lib.GenFake(b)

	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWs(b, &websocket.Upgrader{}))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

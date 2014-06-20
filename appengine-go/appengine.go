// +build appengine

package main

import (
	"net/http"

	"github.com/calsol/teleserver/embedded"
	"github.com/calsol/teleserver/lib"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stvnrhodes/broadcaster"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// db, err := lib.NewDB("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	b := broadcaster.New()

	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWS(b, &websocket.Upgrader{}))
	// r.HandleFunc("/api/graphs", lib.ServeFlotGraphs(db))
	// r.HandleFunc("/api/latest", lib.ServeLatest(db))

	r.PathPrefix("/").HandlerFunc(embedded.ServeFiles)

	http.Handle("/", r)
}

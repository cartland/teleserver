// +build appengine

package main

import (
	"log"
	"net/http"
	"path"

	"go/build"

	"github.com/calsol/teleserver/embedded"
	"github.com/calsol/teleserver/lib"
	"github.com/gorilla/mux"
	"github.com/stvnrhodes/broadcaster"
)

func init() {

	// All data is distributed to all web connections,
	b := broadcaster.New()

	// We handle websocket connections and allow fetching limited historical data.
	r := mux.NewRouter()
	r.HandleFunc("/data/{name}.json", lib.ServeJSON(b))

	// We either serve from embedded so that the binary is standalone, or from
	// /public for rapid development and live changes.
	p, err := build.Default.Import(embedded.BasePkg, "", build.FindOnly)
	if err != nil {
		log.Fatalf("Couldn't find resource files: %v", err)
	}
	root := path.Join(p.Dir, "public")
	log.Println("Serving content from", root)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(root)))

	http.Handle("/", r)
}

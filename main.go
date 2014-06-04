package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":8080", r)
}

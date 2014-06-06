package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/stvnrhodes/teleserver/lib"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stvnrhodes/broadcaster"
	"github.com/tarm/goserial"
)

func main() {
	port := flag.Int("port", 8080, "Port for the webserver")
	uart := flag.String("serial", "", "Serial port for talking to the car")
	baud := flag.Int("baud", 115200, "Baud rate for the serial port")
	flag.Parse()

	if *uart == "" {
		fmt.Println("    Must supply valid port to -serial.")
		fmt.Println("    If on windows this will look like COM45.")
		fmt.Println("    If on a *nix, this will look like /dev/tty.usbmodem1412")
		return
	}

	b := broadcaster.New()
	go lib.GenFake(b)

	p, err := serial.OpenPort(&serial.Config{Name: *uart, Baud: *baud})
	if err != nil {
		log.Fatal(err)
	}
	go lib.Read(p, b)

	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWs(b, &websocket.Upgrader{}))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

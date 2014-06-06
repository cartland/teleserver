package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	fake := flag.Bool("fake", false, "Generate fake data and ignore serial")
	file := flag.String("log_file", "", "If provided, create a log file and write logs to it")
	flag.Parse()

	b := broadcaster.New()
	if *file != "" {
		f, err := os.OpenFile(*file, os.O_CREATE|os.O_APPEND, os.ModeAppend|0755)
		if err != nil {
			log.Fatal(err)
		}
		_ = f
		log.Fatal("Logging is not implemented yet :-(")
	}

	if *fake {
		go lib.GenFake(b)
	} else {
		if *uart == "" {
			fmt.Println("    Must supply valid port to -serial.")
			fmt.Println("    If on windows this will look like COM45.")
			fmt.Println("    If on a *nix, this will look like /dev/tty.usbmodem1412")
			return
		}
		p, err := serial.OpenPort(&serial.Config{Name: *uart, Baud: *baud})
		if err != nil {
			log.Fatal(err)
		}
		go lib.Read(p, b)
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWS(b, &websocket.Upgrader{}))
	r.HandleFunc("/data/{name}.json", lib.ServeJSON(b))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

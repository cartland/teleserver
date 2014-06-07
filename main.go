package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/calsol/teleserver/lib"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stvnrhodes/broadcaster"
	"github.com/tarm/goserial"
)

func main() {
	port := flag.Int("port", 8080, "Port for the webserver")
	uart := flag.String("serial", "", "Serial port for talking to the car")
	baud := flag.Int("baud", 115200, "Baud rate for the serial port")
	can := flag.Bool("can", false, "Treat the serial as streaming CAN, not JSON.")
	fake := flag.Bool("fake", false, "Generate fake data and ignore serial")
	file := flag.String("log_file", "_tmp", "Prefix for the log file. The log file name is based on the time")
	flag.Parse()

	b := broadcaster.New()
	go lib.LogToFile(*file, b)

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
		if *can {
			go lib.ReadCAN(p, b)
		} else {
			go lib.ReadJSON(p, b)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWS(b, &websocket.Upgrader{}))
	r.HandleFunc("/data/{name}.json", lib.ServeJSON(b))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Printf("Starting server on port %v", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

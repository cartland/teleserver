// +build !appengine

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/calsol/teleserver/can"
	"github.com/calsol/teleserver/embedded"
	"github.com/calsol/teleserver/lib"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stvnrhodes/broadcaster"
	"github.com/tarm/goserial"
)

func main() {
	port := flag.Int("port", 8080, "Port for the webserver")
	uart := flag.String("serial", "", "Serial port for talking to the car")
	baud := flag.Int("baud", 115200, "Baud rate for the serial port")
	canAddr := flag.String("can_addr", "", "Port for SocketCAN.")
	fake := flag.Bool("fake", false, "Generate fake data and ignore serial")
	logFile := flag.String("alsologto", "", "Log to stdout and this file")
	sqlite := flag.String("sqlite", "", "Sqlite file that stores all messages sent over CAN. By default it creates a temporary db")
	embed := flag.Bool("use_embedded", false, "Use the embedded files generated by go-bindata to make the binary portable")
	flag.Parse()

	s, err := sql.Open("sqlite3", *sqlite)
	if err != nil {
		log.Fatal(err)
	}
	db, err := lib.NewDB(s)
	if err != nil {
		log.Fatal(err)
	}

	// Split logging to another file if necessary.
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		log.SetOutput(io.MultiWriter(os.Stdout, f))
		log.Printf("Logging to stdout and %s", *logFile)
	}

	// All data is distributed to all web connections,
	b := broadcaster.New()

	go db.WriteMessages(b)

	// sendToCan is a dummy function that can be replaced by a real implementation
	sendToCAN := lib.HandleSendToCAN(nil)

	// We select one of the possible sources to cast data nto the broadcaster
	if *fake {
		go lib.GenFake(b)
	} else if *uart != "" {
		p, err := serial.OpenPort(&serial.Config{Name: *uart, Baud: *baud})
		if err != nil {
			log.Fatal(err)
		}
		defer p.Close()
		go lib.ReadCAN(lib.NewXSPCANReader(p), b)
	} else if *canAddr != "" {
		c, err := can.Dial(*canAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		sendToCAN = lib.HandleSendToCAN(c)
		go lib.ReadCAN(lib.NewSocketCANReader(c), b)
	} else {
		fmt.Println("    Must use -fake or supply valid port to -serial or -can.")
		fmt.Println("    - fake will generate fake sinusoidal data.")
		fmt.Println("    - If on windows serial port will look like COM45.")
		fmt.Println("    - If on a *nix, serial port will look like /dev/tty.usbmodem1412")
		fmt.Println("    - can only works on a system with SocketCAN, and will probably be can0 or similar.")
		return
	}

	// We handle websocket connections and allow fetching limited historical data.
	r := mux.NewRouter()
	r.HandleFunc("/ws", lib.ServeWS(b, &websocket.Upgrader{}))
	r.HandleFunc("/graphs/{name}.json", lib.ServeJSONGraphs(b))
	r.HandleFunc("/send{type}", sendToCAN).Methods("POST")

	// We either serve from embedded so that the binary is standalone, or from
	// /public for rapid development and live changes.
	if *embed {
		log.Println("Serving embedded content")
		r.PathPrefix("/").HandlerFunc(embedded.ServeFiles)
	} else {
		p, err := build.Default.Import(embedded.BasePkg, "", build.FindOnly)
		if err != nil {
			log.Fatalf("Couldn't find resource files: %v", err)
		}
		root := path.Join(p.Dir, "public")
		log.Println("Serving content from", root)
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(root)))
	}

	log.Printf("Starting server on port %v", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

package lib

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between pinging for presence.
	pingPeriod = 30 * time.Second
	// Maximum time to wait when trying to send a message to the client.
	writeWait = time.Second
)

func writer(ws *websocket.Conn, ch <-chan interface{}) {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	for {
		select {
		case m := <-ch:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteJSON(m); err != nil {
				log.Println("Disconnected from client: ", err)
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Disconnected from client: ", err)
				return
			}
		}
	}
}

// ServeWS creates a http.HandlerFunc that upgrades a connection to websockets
// and sends it broadcasted data as JSON.
func ServeWS(b broadcaster.Caster, u *websocket.Upgrader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := u.Upgrade(w, r, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
			return
		}
		defer ws.Close()

		done := make(chan struct{})
		ch := b.Subscribe(done)
		writer(ws, ch)
		close(done)
	}
}

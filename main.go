package main

import (
	"fmt"
	"io"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) { io.Copy(ws, ws) }

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, webpage)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloWorld)
	r.Handle("/ws", websocket.Handler(EchoServer))
	http.ListenAndServe(":8080", r)
}

const webpage = `
<!DOCTYPE html>
<html>
<body>
There's a script here, trust me!
<script>
var ws = new WebSocket("ws://" + window.location.host + "/ws");
ws.onopen = function (e) {
  ws.send("Here's some text that the server is urgently awaiting!");
};
ws.onmessage = function (e) {
  console.log(e.data);
}
</script>
</body>
</html>
`

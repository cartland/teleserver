package teleserver

import (
	"encoding/json"
	"math"
	"time"

	"code.google.com/p/go.net/websocket"
)

var startTime = time.Now()

type Metric struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func getSpeed() Metric {
	t := time.Since(startTime)
	return Metric{Type: "speed", Value: 20 + 2*math.Cos(t.Seconds())}
}

func getVolt() Metric {
	t := time.Since(startTime)
	return Metric{Type: "voltage", Value: 120 + 20*math.Sin(t.Seconds())}
}

func getSolar() Metric {
	t := time.Since(startTime)
	return Metric{Type: "solar", Value: 1000 + 200*math.Sin(t.Seconds())}
}

func MetricsServer(ws *websocket.Conn) {
	e := json.NewEncoder(ws)
	for {
		e.Encode(getSpeed())
		e.Encode(getVolt())
		e.Encode(getSolar())
		time.Sleep(100 * time.Millisecond)
	}
}

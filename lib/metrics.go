package lib

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"time"

	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between fake data readings
	metricPeriod = 500 * time.Millisecond
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

func Read(r io.Reader, b broadcaster.Caster) {
	d := json.NewDecoder(r)
	for {
		m := &Metric{}
		if err := d.Decode(m); err != nil {
			log.Print("issues reading serial: ", err)
		}
		b.Cast(m)
	}
}

func GenFake(b broadcaster.Caster) {
	for {
		// b.Cast(getSpeed())
		b.Cast(getVolt())
		b.Cast(getSolar())
		time.Sleep(metricPeriod)
	}
}

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

func getVolt() Metric {
	t := time.Since(startTime)
	return Metric{Type: "voltage", Value: 120 + 20*math.Sin(t.Seconds())}
}

func getSolar() Metric {
	t := time.Since(startTime)
	return Metric{Type: "solar", Value: 1000 + 200*math.Cos(t.Seconds())}
}

func Read(r io.Reader, b broadcaster.Caster) {
	// Read until new line in case we start in the middle
	p := make([]byte, 1)
	for _, err := r.Read(p); p[0] != '\n'; _, err = r.Read(p) {
		if err != nil {
			log.Fatal(err)
		}
	}

	d := json.NewDecoder(r)
	for {
		m := &Metric{}
		if err := d.Decode(m); err != nil {
			log.Fatal("issues reading serial: ", err)
		}
		b.Cast(m)
	}
}

func GenFake(b broadcaster.Caster) {
	for {
		b.Cast(getVolt())
		b.Cast(getSolar())
		time.Sleep(metricPeriod)
	}
}

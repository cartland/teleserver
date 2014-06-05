package lib

import (
	"math"
	"time"
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

func GenFake(b *Broadcaster) {
	for {
		b.Cast(getSpeed())
		b.Cast(getVolt())
		b.Cast(getSolar())
		time.Sleep(metricPeriod)
	}
}

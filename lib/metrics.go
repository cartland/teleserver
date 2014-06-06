package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between fake data readings
	metricPeriod = 200 * time.Millisecond

	// How long to look back for json data
	bufferedTime = 20 * time.Second
)

var startTime = time.Now()

type Metric struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

func ms() int64 {
	return toMS(time.Now())
}
func toMS(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func getSpeed() Metric {
	t := time.Since(startTime)
	return Metric{Type: "speed", Value: 50 + 10*math.Sin(t.Seconds()), Time: ms()}
}

func getVolt() Metric {
	t := time.Since(startTime)
	return Metric{Type: "voltage", Value: 120 + 20*math.Cos(t.Seconds()), Time: ms()}
}

func getSolar() Metric {
	t := time.Since(startTime)
	return Metric{Type: "solar", Value: 1000 + 200*math.Cos(t.Seconds()), Time: ms()}
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
		m.Time = ms()
		b.Cast(m)
	}
}

func GenFake(b broadcaster.Caster) {
	for {
		b.Cast(getSpeed())
		b.Cast(getVolt())
		b.Cast(getSolar())
		time.Sleep(metricPeriod)
	}
}

type point struct {
	x int64
	y float64
}

func (p point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%f]", p.x, p.y)), nil
}

type GraphData struct {
	Label string  `json:"label"`
	Data  []point `json:"data"`
}

func fromMS(milliseconds int64) time.Time {
	return time.Unix(0, 1000000*milliseconds)
}

// updateSeries adds p to ps and removes any points from longer ago than oldest.
func updateSeries(ps []point, p point, oldest time.Duration) []point {
	now, then := fromMS(p.x), fromMS(ps[0].x)
	for len(ps) > 10 && now.Sub(then) > oldest {
		ps = ps[1:]
		then = fromMS(ps[0].x)
	}
	return append(ps, p)
}

// ServeJSON remembers broadcast metrics for the 5 minutes and serves them up
// when requested based on the type field.
func ServeJSON(b broadcaster.Caster) func(http.ResponseWriter, *http.Request) {
	var mu sync.Mutex
	data := make(map[string]GraphData)
	dataCh := b.Subscribe(nil)
	go func() {
		for d := range dataCh {
			switch d := d.(type) {
			case Metric:
				p := point{x: d.Time, y: d.Value}
				mu.Lock()
				if series, ok := data[d.Type]; ok {
					series.Data = updateSeries(series.Data, p, bufferedTime)
					data[d.Type] = series
				} else {
					data[d.Type] = GraphData{
						Label: d.Type,
						Data:  []point{p},
					}
				}
				mu.Unlock()
			}
		}
	}()
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		e := json.NewEncoder(w)
		mu.Lock()
		d := data[name]
		mu.Unlock()
		if err := e.Encode(d); err != nil {
			log.Print("Failed to send json: ", err)
		}
	}
}

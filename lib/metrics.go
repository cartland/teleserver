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
	// This constant should be kept in sync with public/main.js
	bufferedTime = 20 * time.Second
)

var startTime = time.Now()

// Metric represents a json object with the type of metric, the value of the
// metric, and the timestamp in milliseconds since epoch.
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

// Read will continually read json from the io.Reader, interpret it as a metric,
// and send the metric through the broadcaster.
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

// GenFake broadcasts fake data for speed, voltage, and power.
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

// GraphData represents the data for a flot graph.
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

// ServeJSON remembers broadcast metrics for bufferedTime and serves them up
// when requested based on the type field. It uses the {name} variable from
// mux.Vars to determine which data to serve, and will serve all graphs in an
// array if {name} == "all".
func ServeJSON(b broadcaster.Caster) func(http.ResponseWriter, *http.Request) {
	var mu sync.Mutex
	data := make(map[string]GraphData)
	dataCh := b.Subscribe(nil)

	go func() {
		// Get broadcast data and process it into a map for requests.
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		name := mux.Vars(r)["name"]
		e := json.NewEncoder(w)
		var d interface{}
		mu.Lock()
		// Special case "all" to return all the data.
		if name == "all" {
			var s []GraphData
			for _, g := range data {
				s = append(s, g)
			}
			d = s
		} else {
			// A nonexistant name should 404.
			if s, ok := data[name]; ok {
				d = s
			}
		}
		mu.Unlock()
		if d == nil {
			w.WriteHeader(http.StatusNotFound)
			d = GraphData{Label: name}
		}
		if err := e.Encode(d); err != nil {
			log.Print("Failed to send json: ", err)
		}
	}
}

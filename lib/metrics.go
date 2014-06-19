package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/calsol/teleserver/msgs"
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

func getSpeed() msgs.CANPlus {
	t := time.Since(startTime)
	v := float32(50 + 20*math.Cos(t.Seconds()))
	return msgs.NewCANPlus(&msgs.VelocityMeasurement{VehicleVelocity: v})
}

func getPower() msgs.CANPlus {
	t := time.Since(startTime)
	v := float32(100 + 10*math.Sin(t.Seconds()))
	a := v/4 + float32(10*math.Sin(t.Seconds()/1.8))
	return msgs.NewCANPlus(&msgs.BusMeasurement{BusVoltage: v, BusCurrent: a})
}

func getMPPT(id uint16) msgs.CANPlus {
	t := time.Since(startTime)
	msg := msgs.IDToMessage(id).(*msgs.MPPTStatus)
	fid := float64(id)
	msg.ArrayVoltage = uint16(5000 + 500*math.Sin(t.Seconds()+fid))
	msg.ArrayCurrent = uint16(5000 + 500*math.Cos(t.Seconds()+fid))
	msg.BatteryVoltage = uint16(10000 + 1000*math.Cos(t.Seconds()/2+fid))
	msg.Temperature = uint16(2500 + 5000*math.Cos(t.Seconds()/3+fid))
	return msgs.NewCANPlus(msg)
}

// readTill takes bytes from the reader until it sees b.
func readTill(r io.Reader, b byte) {
	p := make([]byte, 1)
	for _, err := r.Read(p); p[0] != b; _, err = r.Read(p) {
		if err != nil {
			log.Fatal(err)
		}
	}
}

// GenFake broadcasts fake data for speed, voltage, and power.
func GenFake(b broadcaster.Caster) {
	for {
		b.Cast(getSpeed())
		b.Cast(getPower())
		b.Cast(getMPPT(0x600))
		b.Cast(getMPPT(0x601))
		b.Cast(getMPPT(0x602))
		b.Cast(getMPPT(0x603))
		time.Sleep(metricPeriod)
	}
}

type point struct {
	x time.Time
	y float64
}

func (p point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%f]", p.x.UnixNano()/1000000, p.y)), nil
}

// GraphData represents the data for a flot graph.
type GraphData struct {
	Label string  `json:"label"`
	Data  []point `json:"data"`
}

// updateSeries adds p to ps and removes any points from longer ago than oldest.
func updateSeries(ps []point, p point, oldest time.Duration) []point {
	now, then := p.x, ps[0].x
	for len(ps) > 10 && now.Sub(then) > oldest {
		ps, then = ps[1:], ps[0].x
		then = ps[0].x
	}
	return append(ps, p)
}

func updateData(data map[string]GraphData, mu *sync.Mutex, name string, p point) {
	mu.Lock()
	defer mu.Unlock()
	if series, ok := data[name]; ok {
		series.Data = updateSeries(series.Data, p, bufferedTime)
		data[name] = series
	} else {
		data[name] = GraphData{
			Label: name,
			Data:  []point{p},
		}
	}
}

// ServeJSON remembers broadcast metrics for bufferedTime and serves them up
// when requested based on the type field. It uses the {name} variable from
// mux.Vars to determine which data to serve, and will serve all graphs in an
// array if {name} == "all".
func ServeJSON(b broadcaster.Caster) func(http.ResponseWriter, *http.Request) {
	var mu sync.Mutex
	data := make(map[string]GraphData)
	dataCh := b.Subscribe(nil)

	updateData := func(name string, p point) {
		mu.Lock()
		defer mu.Unlock()
		if series, ok := data[name]; ok {
			series.Data = updateSeries(series.Data, p, bufferedTime)
			data[name] = series
		} else {
			data[name] = GraphData{
				Label: name,
				Data:  []point{p},
			}
		}
	}

	go func() {
		// Get broadcast data and process it into a map for requests.
		for d := range dataCh {
			switch d := d.(type) {

			case msgs.CANPlus, *msgs.CANPlus:
				m, ok := d.(msgs.CANPlus)
				if !ok {
					m = *d.(*msgs.CANPlus)
				}
				v := reflect.ValueOf(m.CAN).Elem()
				t := reflect.TypeOf(m.CAN).Elem()
				for i := 0; i < t.NumField(); i++ {
					f := t.Field(i)
					val := v.FieldByName(f.Name)
					if k := val.Kind(); k == reflect.Float32 || k == reflect.Float64 {
						updateData(f.Name, point{x: m.Time, y: val.Float()})
					}
				}

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
			s := []GraphData{}
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

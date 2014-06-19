package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/calsol/teleserver/msgs"
)

const (
	// Default time to look back for json data
	bufferedTime = 2 * time.Minute
)

type point struct {
	x time.Time
	y interface{} // This should always be a number
}

func (p point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%v]", p.x.UnixNano()/1000000, p.y)), nil
}

// ServeLatest will find the most recent result for a CAN ID in the database
// and return it. ServeLatest can return multiple CAN IDs.
func ServeLatest(db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		metrics := []*msgs.CANPlus{}
		for _, id := range r.Form["canid"] {
			canid, err := strconv.Atoi(id)
			if err != nil {
				continue
			}
			metric, err := db.GetLatest(uint16(canid))
			if err != nil {
				continue
			}
			metrics = append(metrics, metric)
		}

		e := json.NewEncoder(w)
		if err := e.Encode(metrics); err != nil {
			log.Printf("Failed to encode json: %v", err)
		}
	}
}

// GraphData represents the data for a flot graph.
type GraphData struct {
	Label string  `json:"label"`
	Data  []point `json:"data"`
}

func hasField(c *msgs.CANPlus, name string) bool {
	_, ok := reflect.TypeOf(c.CAN).Elem().FieldByName(name)
	return ok
}

// ServeFlotGraphs serves fields from CAN messages as input data to flot graphs.
func ServeFlotGraphs(db *DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Use buffered time by default
		d := bufferedTime
		if len(r.Form["time"]) > 0 {
			t, err := time.ParseDuration(r.Form["time"][0])
			if err == nil {
				d = t
			}
		}

		metrics := [][]*msgs.CANPlus{}
		for _, id := range r.Form["canid"] {
			canid, err := strconv.Atoi(id)
			if err != nil {
				continue
			}
			metric, err := db.GetSince(d, uint16(canid))
			if err != nil {
				continue
			}
			metrics = append(metrics, metric)
		}

		s := []GraphData{}
		for _, field := range r.Form["field"] {
			for _, metric := range metrics {
				if len(metric) > 0 && hasField(metric[0], field) {
					graph := GraphData{Label: fmt.Sprintf("0x%x - %s", metric[0].CANID, field)}
					for _, m := range metric {
						graph.Data = append(graph.Data, point{m.Time, reflect.ValueOf(m.CAN).Elem().FieldByName(field).Interface()})
					}
					s = append(s, graph)
				}
			}
		}

		e := json.NewEncoder(w)
		if err := e.Encode(s); err != nil {
			log.Printf("Failed to encode json: %v", err)
		}
	}
}

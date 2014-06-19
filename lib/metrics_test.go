package lib_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/msgs"
	"github.com/gorilla/mux"
	"github.com/stvnrhodes/broadcaster"
)

func getHTTP(t *testing.T, url string) string {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestServeJSON(t *testing.T) {
	b := broadcaster.New()
	defer b.Close()
	r := mux.NewRouter()
	r.HandleFunc("/{name}.json", lib.ServeJSONGraphs(b))
	ts := httptest.NewServer(r)
	defer ts.Close()

	bus := &msgs.BusMeasurement{BusVoltage: 0, BusCurrent: 0.5}
	b.Cast(msgs.CANPlus{bus, msgs.GetID(bus), time.Unix(30, 0)})
	bus = &msgs.BusMeasurement{BusVoltage: 1.5, BusCurrent: 3}
	b.Cast(&msgs.CANPlus{bus, msgs.GetID(bus), time.Unix(40, 0)})
	b.Cast(123)
	v := &msgs.VelocityMeasurement{MotorVelocity: 1, VehicleVelocity: 2}
	b.Cast(msgs.CANPlus{v, msgs.GetID(v), time.Unix(40, 0)})

	tests := []struct{ path, want string }{
		{"/badurl", `404 page not found` + "\n"},
		{"/nil.json", `{"label":"nil","data":null}` + "\n"},
		{"/BusVoltage.json", `{"label":"BusVoltage","data":[[30000,0.000000],[40000,1.500000]]}` + "\n"},
		{"/BusCurrent.json", `{"label":"BusCurrent","data":[[30000,0.500000],[40000,3.000000]]}` + "\n"},
		{"/MotorVelocity.json", `{"label":"MotorVelocity","data":[[40000,1.000000]]}` + "\n"},
		{"/VehicleVelocity.json", `{"label":"VehicleVelocity","data":[[40000,2.000000]]}` + "\n"},
		{"/all.json", `[{"label":"BusVoltage","data":[[30000,0.000000],[40000,1.500000]]},{"label":"BusCurrent","data":[[30000,0.500000],[40000,3.000000]]},{"label":"MotorVelocity","data":[[40000,1.000000]]},{"label":"VehicleVelocity","data":[[40000,2.000000]]}]` + "\n"},
	}

	for _, c := range tests {
		got := getHTTP(t, ts.URL+c.path)
		if got != c.want {
			t.Errorf("%v: got %v, want %v", c.path, got, c.want)
		}
	}
}

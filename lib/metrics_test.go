package lib_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/msgs"
	"github.com/gorilla/mux"
)

func makeDB(t *testing.T) *lib.DB {
	s, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db, err := lib.NewDB(s)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

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

func TestServeLatest(t *testing.T) {
	db := makeDB(t)
	bus := &msgs.BusMeasurement{BusVoltage: 0, BusCurrent: 0.5}
	db.WriteCAN(&msgs.CANPlus{bus, msgs.GetID(bus), time.Unix(30, 0)})
	bus = &msgs.BusMeasurement{BusVoltage: 1.5, BusCurrent: 3}
	db.WriteCAN(&msgs.CANPlus{bus, msgs.GetID(bus), time.Unix(40, 0)})
	v := &msgs.VelocityMeasurement{MotorVelocity: 1, VehicleVelocity: 2}
	db.WriteCAN(&msgs.CANPlus{v, msgs.GetID(v), time.Unix(40, 0)})

	r := mux.NewRouter()
	r.HandleFunc("/data", lib.ServeLatest(db))
	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct{ path, want string }{
		{"/data", `[]` + "\n"},
		{"/data?canid=123", `[]` + "\n"},
		{fmt.Sprintf("/data?canid=%d", msgs.GetID(bus)), `[{"CAN":{"BusVoltage":1.5,"BusCurrent":3},"canID":1026,"time":"1970-01-01T00:00:40Z"}]` + "\n"},
		{fmt.Sprintf("/data?canid=%d", msgs.GetID(v)), `[{"CAN":{"MotorVelocity":1,"VehicleVelocity":2},"canID":1027,"time":"1970-01-01T00:00:40Z"}]` + "\n"},
		{fmt.Sprintf("/data?canid=%d&canid=%d", msgs.GetID(bus), msgs.GetID(v)), `[{"CAN":{"BusVoltage":1.5,"BusCurrent":3},"canID":1026,"time":"1970-01-01T00:00:40Z"},{"CAN":{"MotorVelocity":1,"VehicleVelocity":2},"canID":1027,"time":"1970-01-01T00:00:40Z"}]` + "\n"},
	}

	for _, c := range tests {
		got := getHTTP(t, ts.URL+c.path)
		if got != c.want {
			t.Errorf("%v: got %v, want %v", c.path, got, c.want)
		}
	}
}

func TestServeFlot(t *testing.T) {

	db := makeDB(t)

	// TODO(stvn): Test more intersections
	tm := time.Now()
	bus := &msgs.BusMeasurement{BusVoltage: 0, BusCurrent: 0.5}
	db.WriteCAN(&msgs.CANPlus{bus, msgs.GetID(bus), tm.Add(-time.Hour)})
	bus = &msgs.BusMeasurement{BusVoltage: 1.5, BusCurrent: 3}
	db.WriteCAN(&msgs.CANPlus{bus, msgs.GetID(bus), tm.Add(-30 * time.Minute)})
	v := &msgs.VelocityMeasurement{MotorVelocity: 1, VehicleVelocity: 2}
	db.WriteCAN(&msgs.CANPlus{v, msgs.GetID(v), tm.Add(-30 * time.Minute)})

	r := mux.NewRouter()
	r.HandleFunc("/data", lib.ServeFlotGraphs(db))
	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct{ path, want string }{
		{"/data", `[]` + "\n"},
		{"/data?canid=123", `[]` + "\n"},
		{fmt.Sprintf("/data?canid=%d&field=BusVoltage", msgs.GetID(bus)), `[]` + "\n"},
		{
			fmt.Sprintf("/data?canid=%d&field=BusVoltage&time=45m", msgs.GetID(bus)),
			`[{"label":"0x402BusVoltage","data":[[1403157354184,1.5]]}]` + "\n",
		},
		{
			fmt.Sprintf("/data?canid=%d&field=VehicleVelocity&time=45m", msgs.GetID(v)),
			`[{"label":"0x403VehicleVelocity","data":[[1403157734191,2]]}]` + "\n",
		},
		{
			fmt.Sprintf("/data?canid=%d&canid=%d&time=45m&field=BusVoltage", msgs.GetID(v), msgs.GetID(bus)),
			`[{"label":"0x402BusVoltage","data":[[1403157462723,1.5]]}]` + "\n",
		},
		{
			fmt.Sprintf("/data?canid=%d&time=1h20m&field=BusVoltage", msgs.GetID(bus)),
			`[{"label":"0x402BusVoltage","data":[[1403155825986,0],[1403157625986,1.5]]}]` + "\n",
		},
		{
			fmt.Sprintf("/data?canid=%[1]d&canid=%[1]d&time=1h20m&field=BusVoltage", msgs.GetID(bus)),
			`[{"label":"0x402BusVoltage","data":[[1403155825986,0],[1403157625986,1.5]]}]` + "\n",
		},
		{
			fmt.Sprintf("/data?canid=%d&canid=%d&time=45m&field=BusVoltage&field=BusCurrent&field=VehicleVelocity", msgs.GetID(v), msgs.GetID(bus)),
			`[{"label":"0x402BusVoltage","data":[[1403157734191,1.5]]},{"label":"0x402BusCurrent","data":[[1403157734191,3]]},{"label":"0x403VehicleVelocity","data":[[1403157734191,2]]}]` + "\n",
		},
	}

	for _, c := range tests {
		got := getHTTP(t, ts.URL+c.path)
		// We compare lengths to avoid parsing the time.
		if len(got) != len(c.want) {
			t.Errorf("%v: got %v, want %v", c.path, got, c.want)
		}
	}
}

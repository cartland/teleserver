package lib_test

import (
	"database/sql"
	"reflect"
	"sync"
	"testing"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/msgs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stvnrhodes/broadcaster"
)

func TestSQLMessages(t *testing.T) {
	s, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db, err := lib.NewDB(s)
	if err != nil {
		t.Fatal(err)
	}
	b := broadcaster.New()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { wg.Done(); db.WriteMessages(b); wg.Done() }()
	wg.Wait()

	raw := []msgs.CAN{
		&msgs.BusMeasurement{BusVoltage: 0, BusCurrent: 0.5},
		&msgs.BusMeasurement{BusVoltage: 1.5, BusCurrent: 3},
		&msgs.VelocityMeasurement{MotorVelocity: 1, VehicleVelocity: 2},
	}

	msgPlus := []interface{}{
		msgs.NewCANPlus(raw[0]),
		msgs.NewCANPlus(raw[1]),
		msgs.NewCANPlus(raw[2]),
		123,
	}

	for _, m := range msgPlus {
		b.Cast(m)
	}

	wg.Add(1)
	b.Close()
	wg.Wait()

	// Only check for the latest
	for _, i := range []int{1, 2} {
		msg, err := db.GetLatest(msgs.GetID(raw[i]))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(msgPlus[i], *msg) {
			t.Fatalf("%d: got %#v, want %#v", i, msg, msgPlus[i])
		}
	}
}

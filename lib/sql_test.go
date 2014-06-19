package lib_test

import (
	"database/sql"
	"reflect"
	"sync"
	"testing"
	"time"

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

	var msgPlus []*msgs.CANPlus
	for _, msg := range raw {
		m := msgs.NewCANPlus(msg)
		msgPlus = append(msgPlus, &m)
	}

	for _, m := range msgPlus {
		b.Cast(m)
	}

	wg.Add(1)
	b.Close()
	wg.Wait()

	// Only check for the latest messages.
	for _, i := range []int{1, 2} {
		msg, err := db.GetLatest(msgs.GetID(raw[i]))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(msgPlus[i], msg) {
			t.Fatalf("%d: got %#v, want %#v", i, msg, msgPlus[i])
		}
	}

	multiple, err := db.GetSince(time.Minute, msgs.GetID(raw[0]))
	if err != nil {
		t.Fatal(err)
	}
	original := []*msgs.CANPlus{msgPlus[0], msgPlus[1]}
	if !reflect.DeepEqual(multiple, original) {
		t.Fatalf("got %v, want %v", multiple, original)
	}
}

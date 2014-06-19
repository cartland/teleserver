package lib_test

import (
	"database/sql"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/messages"
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

	msgs := []messages.CAN{
		&messages.BusMeasurement{BusVoltage: 0, BusCurrent: 0.5},
		&messages.BusMeasurement{BusVoltage: 1.5, BusCurrent: 3},
		&messages.VelocityMeasurement{MotorVelocity: 1, VehicleVelocity: 2},
	}

	msgPlus := []interface{}{
		&messages.CANPlus{msgs[0], messages.GetID(msgs[0]), time.Unix(50, 0)},
		messages.CANPlus{msgs[1], messages.GetID(msgs[1]), time.Unix(30, 0)},
		messages.CANPlus{msgs[2], messages.GetID(msgs[2]), time.Unix(40, 0)},
		123,
	}

	for _, m := range msgPlus {
		b.Cast(m)
	}

	wg.Add(1)
	b.Close()
	wg.Wait()

	// Only check for the latest
	for _, i := range []int{0, 2} {
		msg, err := db.GetLatest(messages.GetID(msgs[i]))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(msgPlus[i])).Interface(), *msg) {
			t.Fatalf("%d: got %#v, want %#v", i, msg, msgPlus[i])
		}
	}
}

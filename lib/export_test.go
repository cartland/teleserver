package lib_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/msgs"
	"github.com/gorilla/mux"
	"github.com/stvnrhodes/broadcaster"
)

var totallyFakeSecretKey = []byte("123")

func TestSendMessage(t *testing.T) {
	db := makeDB(t)
	i := lib.NewImporter(db, totallyFakeSecretKey)
	r := mux.NewRouter()
	r.Handle("/", i)
	ts := httptest.NewServer(r)
	defer ts.Close()

	b := broadcaster.New()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { wg.Done(); lib.PostOnBroadcast(b, ts.URL, totallyFakeSecretKey); wg.Done() }()
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

	resp, err := http.PostForm(ts.URL, url.Values{"data": {"{1,2}"}, "key": {"AQ=="}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Errorf("Should not accept request: %v", resp)
		s, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(string(s))
	}

	resp, err = lib.PostToURL(ts.URL, "{1,2}", totallyFakeSecretKey)

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

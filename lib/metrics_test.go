package lib_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bitbucket.org/stvnrhodes/teleserver/lib"
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
	r.HandleFunc("/{name}.json", lib.ServeJSON(b))
	ts := httptest.NewServer(r)
	defer ts.Close()

	b.Cast(lib.Metric{Type: "foo", Value: 1.0, Time: 10000})
	b.Cast(lib.Metric{Type: "bar", Value: 2.0, Time: 30000})
	b.Cast(lib.Metric{Type: "foo", Value: 3.0, Time: 40000})
	b.Cast(123)

	tests := []struct{ path, want string }{
		{"/badurl", `404 page not found` + "\n"},
		{"/nil.json", `{"label":"nil","data":null}` + "\n"},
		{"/foo.json", `{"label":"foo","data":[[10000,1.000000],[40000,3.000000]]}` + "\n"},
		{"/bar.json", `{"label":"bar","data":[[30000,2.000000]]}` + "\n"},
		{"/all.json", `[{"label":"foo","data":[[10000,1.000000],[40000,3.000000]]},{"label":"bar","data":[[30000,2.000000]]}]` + "\n"},
	}

	for _, c := range tests {
		got := getHTTP(t, ts.URL+c.path)
		if got != c.want {
			t.Errorf("%v: got %v, want %v", c.path, got, c.want)
		}
	}
}

func TestReadData(t *testing.T) {
	testString := `:3.0}
{"type":"foo","value":1.0}
{"type":"bar","value":2.0}
`
	then := time.Now()
	time.Sleep(time.Millisecond)
	b := broadcaster.New()
	defer b.Close()
	ch := b.Subscribe(nil)

	lib.Read(strings.NewReader(testString), b)

	expected := []*lib.Metric{{Type: "foo", Value: 1.0}, {Type: "bar", Value: 2.0}}
	for _, e := range expected {
		got := <-ch
		m, ok := got.(*lib.Metric)
		if !ok {
			t.Errorf("Wrong type %[1]T for %+[1]v", got)
			continue
		}
		now := time.Now()
		if m.Type != e.Type || m.Value != e.Value {
			t.Errorf("Wrong metric: got %+v, want %+v", m, e)
		}
		if time := lib.FromMS(m.Time); time.Before(then) || time.After(now) {
			t.Errorf("Time %v should be between 0 and %v", time.Sub(then), now.Sub(then))
		}
	}
}

package lib_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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

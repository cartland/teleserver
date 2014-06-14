package lib_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/calsol/teleserver/lib"
	"github.com/gorilla/mux"
)

func postHTTP(t *testing.T, url string, vals url.Values) string {
	resp, err := http.PostForm(url, vals)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestSendCAN(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/send{type}", lib.HandleSendToCAN(nil))
	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		path, want string
		vals       url.Values
	}{
		{"/sendnotype", "notype is an invalid send type\n", url.Values{"id": {"123"}}},
		{
			"/sendbytes",
			`cannot parse length: strconv.ParseInt: parsing "": invalid syntax` + "\n",
			url.Values{"id": {"123"}},
		},
		{
			"/sendbytes",
			`cannot parse byte 0: strconv.ParseInt: parsing "ac": invalid syntax` + "\n",
			url.Values{"id": {"123"}, "length": {"1"}, "byte0": {"ac"}},
		},
		{
			"/sendbytes",
			"CANSocket not running, " +
				"cannot send &can.Frame{ID:0x7b, DataLen:0x3, Padding:[3]uint8{0x0, 0x0, 0x0}, " +
				"Data:[8]uint8{0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0}} over CAN\n",
			url.Values{"id": {"123"}, "length": {"3"},
				"byte0": {"1"}, "byte1": {"2"}, "byte2": {"3"}, "byte3": {"4"}},
		},
		{
			"/sendfloats",
			"CANSocket not running, " +
				"cannot send &can.Frame{ID:0x7b, DataLen:0x8, Padding:[3]uint8{0x0, 0x0, 0x0}, " +
				"Data:[8]uint8{0x0, 0x0, 0x0, 0x3f, 0x9a, 0x99, 0x99, 0x3f}} over CAN\n",
			url.Values{"id": {"123"}, "float0": {"0.5"}, "float1": {"1.2"}},
		},
	}

	for i, c := range tests {
		got := postHTTP(t, ts.URL+c.path, c.vals)
		if got != c.want {
			t.Errorf("%d: %v: got %q, want %q", i, c.path, got, c.want)
		}
	}
}

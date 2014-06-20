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

func TestSendMessage(t *testing.T) {

	r := mux.NewRouter()
	r.HandleFunc("/", lib.ServeHTTPWithHMAC)
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.PostForm(ts.URL, url.Values{"data": {"{1,2}"}, "key": {"AQ=="}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Error("Should not accept request: %v", resp)
		s, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(string(s))
	}

	resp, err = lib.PostToURL(ts.URL, "{1,2}")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Should accept request: %v", resp)
		s, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(string(s))
	}

}

package embedded_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"go/build"

	"github.com/calsol/teleserver/embedded"
	"github.com/gorilla/mux"
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

func TestUpdatedEmbedded(t *testing.T) {
	embed := mux.NewRouter()
	embed.PathPrefix("/").HandlerFunc(embedded.ServeFiles)
	embedServ := httptest.NewServer(embed)
	defer embedServ.Close()

	fresh := mux.NewRouter()
	p, err := build.Default.Import(embedded.BasePkg, "", build.FindOnly)
	if err != nil {
		t.Fatalf("Couldn't find resource files: %v", err)
	}
	fresh.PathPrefix("/").Handler(http.FileServer(http.Dir(path.Join(p.Dir, "public"))))
	freshServ := httptest.NewServer(fresh)
	defer freshServ.Close()

	for _, f := range embedded.AssetNames() {
		f = strings.TrimPrefix(f, "public")
		wantURL, gotURL := freshServ.URL+f, embedServ.URL+f
		if want, got := getHTTP(t, wantURL), getHTTP(t, gotURL); want != got {
			t.Errorf("Embedded file %v does not match expected: got %v, want %v", f, got, want)
			t.Errorf(`Try running 'go-bindata -o embedded/assets.go -ignore \\.bower\.json -ignore bower_components/marked -ignore \demos -ignore \core-tests -ignore bower_components/highlightjs -nomemcopy -pkg embedded public/...'`)
		}
	}
}

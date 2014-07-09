package embedded

import (
	"net/http"
	"path"
	"time"
)

// BasePkg is the package where files to serve are normally found
const BasePkg = "github.com/calsol/teleserver"

var modtime = time.Now()

// checkLastModified is taken from net/http/fs.go
// return value is whether this request is now complete.
func checkLastModified(w http.ResponseWriter, r *http.Request) bool {
	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := w.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return false
}

// ServeFiles will act as a file server for all Assets
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	if checkLastModified(w, r) {
		return
	}
	file := path.Join("public", r.URL.Path)
	b, err := Asset(file)
	if err != nil {
		b, err = Asset(path.Join(file, "index.html"))
		if err != nil {
			http.Error(w, "Could not find "+file, 404)
			return
		}
	}
	switch path.Ext(file) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	}
	w.Write(b)
}

package embedded

import (
	"net/http"
	"path"
)

// BasePkg is the package where files to serve are normally found
const BasePkg = "github.com/calsol/teleserver"

// ServeFiles will act as a file server for all Assets
func ServeFiles(w http.ResponseWriter, r *http.Request) {
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

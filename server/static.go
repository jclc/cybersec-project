package server

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	RegisterHandler("/static", ServeStatic)
}

// ServeStatic serves a static file from the static/ directory. Boy, I wish there
// was a built-in method to do this!!
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Instead of having the file path be the part of the url itself like how everyone
	// else does it (super lame) we'll make it part of a query (super cool)
	q := r.URL.Query()
	file, ok := q["file"]
	if !ok || len(file) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Assume the working directory is the same as the repo root directory.
	// static/ should be a subdir. We just join the paths, ez
	fp := filepath.Join("./static", file[0])
	f, err := os.Open(fp)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	contentType := mime.TypeByExtension(filepath.Ext(file[0]))
	if contentType == "" {
		// Honestly it's probably just text. Who cares?
		contentType = "text/plain"
	}
	w.Header().Set("Content-Type", contentType)
	io.Copy(w, f)
}

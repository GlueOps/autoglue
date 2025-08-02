package api

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed ui/*
var embeddedUI embed.FS

func StaticHandler() http.Handler {
	contentFS, _ := fs.Sub(embeddedUI, "ui")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") ||
			strings.HasPrefix(r.URL.Path, "/swagger") ||
			strings.HasPrefix(r.URL.Path, "/debug/pprof/") {
			http.NotFound(w, r) // or let it fall through to the mux router
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		file, err := contentFS.Open(path)
		if err != nil {
			// fallback to index.html
			serveIndexHTML(w, r, contentFS)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			serveIndexHTML(w, r, contentFS)
			return
		}

		// Read full file into memory to serve as ReadSeeker
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), bytesReader(data))
	})
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request, fsys fs.FS) {
	file, err := fsys.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "unable to stat index.html", http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read index.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	http.ServeContent(w, r, "index.html", stat.ModTime(), bytesReader(data))
}

func bytesReader(b []byte) io.ReadSeeker {
	return io.NewSectionReader(strings.NewReader(string(b)), 0, int64(len(b)))
}

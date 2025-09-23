package ui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// NOTE: Vite outputs to ui/dist with assets in dist/assets.
// If you add more nested folders in the future, include them here too.

//go:embed dist
var distFS embed.FS

// spaFileSystem serves embedded dist/ files with SPA fallback to index.html
type spaFileSystem struct {
	fs fs.FS
}

func (s spaFileSystem) Open(name string) (fs.File, error) {
	// Normalize, strip leading slash
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	// Try exact file
	f, err := s.fs.Open(name)
	if err == nil {
		return f, nil
	}

	// If the requested file doesn't exist, fall back to index.html for SPA routes
	// BUT only if it's not obviously a static asset extension
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".js", ".css", ".map", ".json", ".txt", ".ico", ".png", ".jpg", ".jpeg",
		".svg", ".webp", ".gif", ".woff", ".woff2", ".ttf", ".otf", ".eot", ".wasm":
		return nil, fs.ErrNotExist
	}

	return s.fs.Open("index.html")
}

func newDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}

// SPAHandler returns an http.Handler that serves the embedded UI (with caching)
func SPAHandler() (http.Handler, error) {
	sub, err := newDistFS()
	if err != nil {
		return nil, err
	}

	// Wrap with our SPA filesystem and our own file server to control headers.
	spa := spaFileSystem{fs: sub}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent /api, /swagger, /debug/pprof from being eaten by SPA fallback.
		if strings.HasPrefix(r.URL.Path, "/api/") ||
			r.URL.Path == "/api" ||
			strings.HasPrefix(r.URL.Path, "/swagger") ||
			strings.HasPrefix(r.URL.Path, "/debug/pprof") {
			http.NotFound(w, r)
			return
		}

		// Open file (or fallback to index.html)
		filePath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if filePath == "" {
			filePath = "index.html"
		}
		f, err := spa.Open(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		// Guess content-type by suffix (let Go detect if possible)
		// Serve with gentle caching: long for assets, short for HTML
		if strings.HasSuffix(filePath, ".html") {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			// Vite assets are hashed; safe to cache
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		info, _ := f.Stat()
		modTime := time.Now()
		if info != nil {
			modTime = info.ModTime()
		}

		// Serve content
		http.ServeContent(w, r, filePath, modTime, file{f})
	}), nil
}

// file wraps fs.File to implement io.ReadSeeker if possible (for ServeContent)
type file struct{ fs.File }

func (f file) Seek(offset int64, whence int) (int64, error) {
	if s, ok := f.File.(io.Seeker); ok {
		return s.Seek(offset, whence)
	}
	// Fallback: not seekable
	return 0, fs.ErrInvalid
}

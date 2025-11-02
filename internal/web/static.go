package web

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

// NOTE: Vite outputs to web/dist with assets in dist/assets.
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
		".svg", ".webp", ".gif", ".woff", ".woff2", ".ttf", ".otf", ".eot", ".wasm", ".br", ".gz":
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
	spa := spaFileSystem{fs: sub}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") ||
			r.URL.Path == "/api" ||
			strings.HasPrefix(r.URL.Path, "/swagger") ||
			strings.HasPrefix(r.URL.Path, "/debug/pprof") {
			http.NotFound(w, r)
			return
		}

		filePath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// Try compressed variants for assets and HTML
		// NOTE: we only change *Content-Encoding*; Content-Type derives from original ext
		// Always vary on Accept-Encoding
		w.Header().Add("Vary", "Accept-Encoding")

		enc := r.Header.Get("Accept-Encoding")
		if tryServeCompressed(w, r, spa, filePath, enc) {
			return
		}

		// Fallback: normal open (or SPA fallback)
		f, err := spa.Open(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		if strings.HasSuffix(filePath, ".html") {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		info, _ := f.Stat()
		modTime := time.Now()
		if info != nil {
			modTime = info.ModTime()
		}
		http.ServeContent(w, r, filePath, modTime, file{f})
	}), nil
}

func tryServeCompressed(w http.ResponseWriter, r *http.Request, spa spaFileSystem, filePath, enc string) bool {
	wantsBR := strings.Contains(enc, "br")
	wantsGZ := strings.Contains(enc, "gzip")

	type cand struct {
		logical  string // MIME/type decision uses this (uncompressed name)
		physical string // actual file we open (with .br/.gz)
		enc      string
	}

	var cands []cand

	// 1) direct compressed variant of requested path (rare for SPA routes, but cheap to try)
	if wantsBR {
		cands = append(cands, cand{logical: filePath, physical: filePath + ".br", enc: "br"})
	}
	if wantsGZ {
		cands = append(cands, cand{logical: filePath, physical: filePath + ".gz", enc: "gzip"})
	}

	// 2) SPA route: fall back to compressed index.html
	if filepath.Ext(filePath) == "" {
		if wantsBR {
			cands = append(cands, cand{logical: "index.html", physical: "index.html.br", enc: "br"})
		}
		if wantsGZ {
			cands = append(cands, cand{logical: "index.html", physical: "index.html.gz", enc: "gzip"})
		}
	}

	for _, c := range cands {
		f, err := spa.fs.Open(c.physical) // open EXACT path so we don't accidentally get SPA fallback
		if err != nil {
			continue
		}
		defer f.Close()

		// Cache headers
		if strings.HasSuffix(c.logical, ".html") {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		if ct := mimeByExt(path.Ext(c.logical)); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		w.Header().Set("Content-Encoding", c.enc)
		w.Header().Add("Vary", "Accept-Encoding")

		info, _ := f.Stat()
		modTime := time.Now()
		if info != nil {
			modTime = info.ModTime()
		}

		// Serve the precompressed bytes
		http.ServeContent(w, r, c.physical, modTime, file{f})
		return true
	}
	return false
}

func serveIfExists(w http.ResponseWriter, r *http.Request, spa spaFileSystem, filePath, ext, encoding string) bool {
	cf := filePath + ext
	f, err := spa.Open(cf)
	if err != nil {
		return false
	}
	defer f.Close()

	// Set caching headers
	if strings.HasSuffix(filePath, ".html") {
		w.Header().Set("Cache-Control", "no-cache")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
	// Preserve original content type by extension of *uncompressed* file
	if ct := mimeByExt(path.Ext(filePath)); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Content-Encoding", encoding)

	info, _ := f.Stat()
	modTime := time.Now()
	if info != nil {
		modTime = info.ModTime()
	}

	// Serve the compressed bytes as an io.ReadSeeker if possible
	http.ServeContent(w, r, cf, modTime, file{f})
	return true
}

func mimeByExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css; charset=utf-8"
	case ".json":
		return "application/json"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	default:
		return "" // let Go sniff if empty
	}
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

package api

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Recoverer turns a panic into a 500 and logs it as one structured event.
//
// chi's middleware.Recoverer writes a pretty stack straight to stderr, which
// lands in the pod log as an unstructured blob with no request_id — so the
// stack and the 500 that zeroLogMiddleware records are two separate things you
// have to correlate by eye. This keeps them together.
//
// http.ErrAbortHandler is re-panicked rather than swallowed: net/http uses it
// to signal a deliberately aborted response, and the server expects to see it.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler {
				panic(rec)
			}

			log.Error().
				Interface("panic", rec).
				Str("request_id", middleware.GetReqID(r.Context())).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Bytes("stack", debug.Stack()).
				Msg("http_panic")

			// The response may already be partially written, in which case this
			// header write is a no-op and the client sees a truncated body. That
			// is unavoidable, and the log above is the record that matters.
			w.WriteHeader(http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	})
}

package api

import (
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/config"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HSTS (enable only over TLS/behind HTTPS)
		// HSTS only when not in dev and over TLS/behind a proxy that terminates TLS
		if !config.IsDev() {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=(), interest-cohort=()")

		if config.IsDev() {
			// --- Relaxed CSP for Vite dev server & Google Fonts ---
			// Allows inline/eval for React Refresh preamble, HMR websocket, and fonts.
			// Tighten these as you move to prod or self-host fonts.
			w.Header().Set("Content-Security-Policy", strings.Join([]string{
				"default-src 'self'",
				"base-uri 'self'",
				"form-action 'self'",
				// Vite dev & inline preamble/eval:
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' http://localhost:5173",
				// allow dev style + Google Fonts
				"style-src 'self' 'unsafe-inline' http://localhost:5173 https://fonts.googleapis.com",
				"img-src 'self' data: blob:",
				// Google font files
				"font-src 'self' data: https://fonts.gstatic.com",
				// HMR connections
				"connect-src 'self' http://localhost:5173 ws://localhost:5173 ws://localhost:8080 https://api.github.com",
				"frame-ancestors 'none'",
			}, "; "))
		} else {
			// --- Strict CSP for production ---
			// If you keep using Google Fonts in prod, add:
			//   style-src ... https://fonts.googleapis.com
			//   font-src  ... https://fonts.gstatic.com
			// Recommended: self-host fonts in prod and keep these tight.
			w.Header().Set("Content-Security-Policy", strings.Join([]string{
				"default-src 'self'",
				"base-uri 'self'",
				"form-action 'self'",
				"script-src 'self' 'unsafe-inline'",
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
				"img-src 'self' data: blob:",
				"font-src 'self' data: https://fonts.gstatic.com",
				"connect-src 'self' ws://localhost:8080 https://api.github.com",
				"frame-ancestors 'none'",
			}, "; "))
		}

		next.ServeHTTP(w, r)
	})
}

package config

import (
	"net"
	"net/url"
	"os"
	"strings"
)

// defaultDevOrigins are the browser origins allowed when CORS_ALLOWED_ORIGINS
// is unset. Loopback aliases are added by AllowedOrigins.
var defaultDevOrigins = []string{
	"http://localhost:5173",
	"http://localhost:8080",
}

// AllowedOrigins returns the browser origins permitted to talk to the API,
// with loopback aliases expanded.
//
// Browsers treat http://localhost:8080 and http://127.0.0.1:8080 as entirely
// different origins, so anything origin-scoped (CORS, postMessage) has to name
// both or it silently fails. Note this does NOT make the two interchangeable
// for cookies: cookies are host-scoped, so a session set on localhost is never
// sent from 127.0.0.1. Pick one host and stay on it.
func AllowedOrigins() []string {
	var raw []string
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		for _, p := range strings.Split(v, ",") {
			if s := strings.TrimSpace(p); s != "" {
				raw = append(raw, s)
			}
		}
	}
	if len(raw) == 0 {
		raw = defaultDevOrigins
	}
	return withLoopbackAliases(raw)
}

// withLoopbackAliases returns origins plus the equivalent loopback spelling of
// each (localhost <-> 127.0.0.1), preserving order and dropping duplicates.
func withLoopbackAliases(origins []string) []string {
	seen := make(map[string]bool, len(origins)*2)
	out := make([]string, 0, len(origins)*2)

	add := func(o string) {
		if o == "" || seen[o] {
			return
		}
		seen[o] = true
		out = append(out, o)
	}

	for _, o := range origins {
		add(o)
		if alias := loopbackAlias(o); alias != "" {
			add(alias)
		}
	}
	return out
}

// loopbackAlias returns the sibling loopback spelling of an origin, or "" when
// the origin is not a loopback host.
func loopbackAlias(origin string) string {
	u, err := url.Parse(origin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}

	host := strings.ToLower(u.Hostname())
	var other string
	switch host {
	case "localhost":
		other = "127.0.0.1"
	case "127.0.0.1":
		other = "localhost"
	default:
		return ""
	}

	if port := u.Port(); port != "" {
		other = net.JoinHostPort(other, port)
	}
	return (&url.URL{Scheme: u.Scheme, Host: other}).String()
}

// IsAllowedOrigin reports whether origin is in AllowedOrigins.
func IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, a := range AllowedOrigins() {
		if strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}

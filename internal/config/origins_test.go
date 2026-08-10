package config

import (
	"slices"
	"testing"
)

func TestAllowedOriginsExpandsLoopbackAliases(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	got := AllowedOrigins()
	for _, want := range []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	} {
		if !slices.Contains(got, want) {
			t.Errorf("AllowedOrigins() missing %q; got %v", want, got)
		}
	}
}

func TestAllowedOriginsHonoursEnvAndStillAliases(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://127.0.0.1:3000, https://app.example.com")

	got := AllowedOrigins()

	if !slices.Contains(got, "http://localhost:3000") {
		t.Errorf("expected localhost alias for 127.0.0.1:3000; got %v", got)
	}
	if !slices.Contains(got, "https://app.example.com") {
		t.Errorf("expected the configured origin to survive; got %v", got)
	}
	// A non-loopback host must not grow an alias.
	if slices.Contains(got, "https://localhost") {
		t.Errorf("non-loopback origin was aliased; got %v", got)
	}
	// The defaults must not leak in once the env is set.
	if slices.Contains(got, "http://localhost:5173") {
		t.Errorf("defaults leaked despite CORS_ALLOWED_ORIGINS being set; got %v", got)
	}
}

func TestIsAllowedOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:8080")

	cases := map[string]bool{
		"http://localhost:8080":  true,
		"http://127.0.0.1:8080":  true, // the alias this whole file exists for
		"http://localhost:9999":  false,
		"https://localhost:8080": false, // scheme must match
		"http://evil.example":    false,
		"":                       false,
	}
	for origin, want := range cases {
		if got := IsAllowedOrigin(origin); got != want {
			t.Errorf("IsAllowedOrigin(%q) = %v, want %v", origin, got, want)
		}
	}
}

func TestLoopbackAliasIgnoresNonLoopback(t *testing.T) {
	if got := loopbackAlias("https://app.example.com"); got != "" {
		t.Errorf("loopbackAlias(non-loopback) = %q, want empty", got)
	}
	if got := loopbackAlias("not a url"); got != "" {
		t.Errorf("loopbackAlias(garbage) = %q, want empty", got)
	}
	if got := loopbackAlias("http://localhost"); got != "http://127.0.0.1" {
		t.Errorf("loopbackAlias(no port) = %q, want http://127.0.0.1", got)
	}
}

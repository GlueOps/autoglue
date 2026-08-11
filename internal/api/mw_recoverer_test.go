package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// captureLog swaps the global logger for one writing to a buffer, so the test
// can assert on what actually reaches the pod log.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Logger
	log.Logger = zerolog.New(&buf)
	t.Cleanup(func() { log.Logger = prev })
	return &buf
}

func TestRecoverer_LogsPanicAsOneStructuredEvent(t *testing.T) {
	buf := captureLog(t)

	h := Recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/explode", nil))

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rr.Code)
	}

	var ev map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &ev); err != nil {
		t.Fatalf("log line is not structured JSON: %v (%s)", err, buf.String())
	}

	if ev["message"] != "http_panic" {
		t.Errorf("message = %v, want http_panic", ev["message"])
	}
	if ev["panic"] != "boom" {
		t.Errorf("panic = %v, want boom", ev["panic"])
	}
	if ev["path"] != "/explode" {
		t.Errorf("path = %v, want /explode", ev["path"])
	}
	// The stack is the whole point: without it this is no better than the 500
	// the request logger already records.
	if s, _ := ev["stack"].(string); s == "" {
		t.Error("stack missing from the panic event")
	}
	if ev["level"] != "error" {
		t.Errorf("level = %v, want error", ev["level"])
	}
}

func TestRecoverer_RepanicsAbortHandler(t *testing.T) {
	// net/http uses ErrAbortHandler to signal a deliberately aborted response.
	// Swallowing it would turn an intentional abort into a bogus 500.
	h := Recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	defer func() {
		if rec := recover(); rec != http.ErrAbortHandler {
			t.Fatalf("recovered %v, want ErrAbortHandler to propagate", rec)
		}
	}()

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	t.Fatal("ErrAbortHandler did not propagate")
}

func TestRecoverer_PassesThroughWhenNoPanic(t *testing.T) {
	h := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	if rr.Code != http.StatusTeapot {
		t.Errorf("status = %d, want the handler's own 418", rr.Code)
	}
}

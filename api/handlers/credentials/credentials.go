package credentials

import (
	"encoding/json"
	"net/http"
	"time"
)

const timeRFC3339 = time.RFC3339

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

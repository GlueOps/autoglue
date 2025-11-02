package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is a simple, reusable error payload.
// swagger:model ErrorResponse
type ErrorResponse struct {
	// A machine-readable error code, e.g. "validation_error"
	// example: validation_error
	Code string `json:"code"`
	// Human-readable message
	// example: slug is required
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, ErrorResponse{Code: code, Message: msg})
}

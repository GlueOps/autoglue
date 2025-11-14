package api

import (
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func mountMetaRoutes(r chi.Router) {
	// Versioned JWKS for swagger
	r.Get("/.well-known/jwks.json", handlers.JWKSHandler)
	r.Get("/healthz", handlers.HealthCheck)
	r.Get("/version", handlers.Version)
}

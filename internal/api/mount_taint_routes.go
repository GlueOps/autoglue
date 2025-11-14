package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountTaintRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/taints", func(t chi.Router) {
		t.Use(authOrg)
		t.Get("/", handlers.ListTaints(db))
		t.Post("/", handlers.CreateTaint(db))
		t.Get("/{id}", handlers.GetTaint(db))
		t.Patch("/{id}", handlers.UpdateTaint(db))
		t.Delete("/{id}", handlers.DeleteTaint(db))
	})
}

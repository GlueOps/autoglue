package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountLabelRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/labels", func(l chi.Router) {
		l.Use(authOrg)
		l.Get("/", handlers.ListLabels(db))
		l.Post("/", handlers.CreateLabel(db))
		l.Get("/{id}", handlers.GetLabel(db))
		l.Patch("/{id}", handlers.UpdateLabel(db))
		l.Delete("/{id}", handlers.DeleteLabel(db))
	})
}

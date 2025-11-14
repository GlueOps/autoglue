package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountAnnotationRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/annotations", func(a chi.Router) {
		a.Use(authOrg)
		a.Get("/", handlers.ListAnnotations(db))
		a.Post("/", handlers.CreateAnnotation(db))
		a.Get("/{id}", handlers.GetAnnotation(db))
		a.Patch("/{id}", handlers.UpdateAnnotation(db))
		a.Delete("/{id}", handlers.DeleteAnnotation(db))
	})
}

package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountServerRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/servers", func(s chi.Router) {
		s.Use(authOrg)
		s.Get("/", handlers.ListServers(db))
		s.Post("/", handlers.CreateServer(db))
		s.Get("/{id}", handlers.GetServer(db))
		s.Patch("/{id}", handlers.UpdateServer(db))
		s.Delete("/{id}", handlers.DeleteServer(db))
		s.Post("/{id}/reset-hostkey", handlers.ResetServerHostKey(db))
	})
}

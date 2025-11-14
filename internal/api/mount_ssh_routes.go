package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountSSHRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/ssh", func(s chi.Router) {
		s.Use(authOrg)
		s.Get("/", handlers.ListPublicSshKeys(db))
		s.Post("/", handlers.CreateSSHKey(db))
		s.Get("/{id}", handlers.GetSSHKey(db))
		s.Delete("/{id}", handlers.DeleteSSHKey(db))
		s.Get("/{id}/download", handlers.DownloadSSHKey(db))
	})
}

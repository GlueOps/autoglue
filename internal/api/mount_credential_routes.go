package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountCredentialRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/credentials", func(c chi.Router) {
		c.Use(authOrg)
		c.Get("/", handlers.ListCredentials(db))
		c.Post("/", handlers.CreateCredential(db))
		c.Get("/{id}", handlers.GetCredential(db))
		c.Patch("/{id}", handlers.UpdateCredential(db))
		c.Delete("/{id}", handlers.DeleteCredential(db))
		c.Post("/{id}/reveal", handlers.RevealCredential(db))
	})
}

package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountMeRoutes(r chi.Router, db *gorm.DB, authUser func(http.Handler) http.Handler) {
	r.Route("/me", func(me chi.Router) {
		me.Use(authUser)

		me.Get("/", handlers.GetMe(db))
		me.Patch("/", handlers.UpdateMe(db))

		me.Get("/api-keys", handlers.ListUserAPIKeys(db))
		me.Post("/api-keys", handlers.CreateUserAPIKey(db))
		me.Delete("/api-keys/{id}", handlers.DeleteUserAPIKey(db))
	})
}

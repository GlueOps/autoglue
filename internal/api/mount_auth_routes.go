package api

import (
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountAuthRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/auth", func(a chi.Router) {
		a.Post("/{provider}/start", handlers.AuthStart(db))
		a.Get("/{provider}/callback", handlers.AuthCallback(db))
		a.Post("/refresh", handlers.Refresh(db))
		a.Post("/logout", handlers.Logout(db))
	})
}

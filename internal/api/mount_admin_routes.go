package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountAdminRoutes(r chi.Router, db *gorm.DB, authUser func(http.Handler) http.Handler) {
	r.Route("/admin", func(admin chi.Router) {
		admin.Route("/actions", func(action chi.Router) {
			action.Use(authUser)
			action.Use(httpmiddleware.RequirePlatformAdmin())

			action.Get("/", handlers.ListActions(db))
			action.Post("/", handlers.CreateAction(db))

			action.Get("/{actionID}", handlers.GetAction(db))
			action.Patch("/{actionID}", handlers.UpdateAction(db))
			action.Delete("/{actionID}", handlers.DeleteAction(db))
		})
	})
}

package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountOrgRoutes(r chi.Router, db *gorm.DB, authUser, authOrg func(http.Handler) http.Handler) {
	r.Route("/orgs", func(o chi.Router) {
		o.Use(authUser)
		o.Get("/", handlers.ListMyOrgs(db))
		o.Post("/", handlers.CreateOrg(db))

		o.Group(func(og chi.Router) {
			og.Use(authOrg)

			og.Get("/{id}", handlers.GetOrg(db))
			og.Patch("/{id}", handlers.UpdateOrg(db))
			og.Delete("/{id}", handlers.DeleteOrg(db))

			// members
			og.Get("/{id}/members", handlers.ListMembers(db))
			og.Post("/{id}/members", handlers.AddOrUpdateMember(db))
			og.Delete("/{id}/members/{user_id}", handlers.RemoveMember(db))

			// org-scoped key/secret pair
			og.Get("/{id}/api-keys", handlers.ListOrgKeys(db))
			og.Post("/{id}/api-keys", handlers.CreateOrgKey(db))
			og.Delete("/{id}/api-keys/{key_id}", handlers.DeleteOrgKey(db))
		})
	})
}

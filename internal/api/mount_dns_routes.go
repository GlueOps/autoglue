package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountDNSRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/dns", func(d chi.Router) {
		d.Use(authOrg)

		d.Get("/domains", handlers.ListDomains(db))
		d.Post("/domains", handlers.CreateDomain(db))
		d.Get("/domains/{id}", handlers.GetDomain(db))
		d.Patch("/domains/{id}", handlers.UpdateDomain(db))
		d.Delete("/domains/{id}", handlers.DeleteDomain(db))

		d.Get("/domains/{domain_id}/records", handlers.ListRecordSets(db))
		d.Post("/domains/{domain_id}/records", handlers.CreateRecordSet(db))
		d.Get("/records/{id}", handlers.GetRecordSet(db))
		d.Patch("/records/{id}", handlers.UpdateRecordSet(db))
		d.Delete("/records/{id}", handlers.DeleteRecordSet(db))
	})
}

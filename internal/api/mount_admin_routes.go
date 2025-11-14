package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountAdminRoutes(r chi.Router, db *gorm.DB, jobs *bg.Jobs, authUser func(http.Handler) http.Handler) {
	r.Route("/admin", func(admin chi.Router) {
		admin.Route("/archer", func(archer chi.Router) {
			archer.Use(authUser)
			archer.Use(httpmiddleware.RequirePlatformAdmin())

			archer.Get("/jobs", handlers.AdminListArcherJobs(db))
			archer.Post("/jobs", handlers.AdminEnqueueArcherJob(db, jobs))
			archer.Post("/jobs/{id}/retry", handlers.AdminRetryArcherJob(db))
			archer.Post("/jobs/{id}/cancel", handlers.AdminCancelArcherJob(db))
			archer.Get("/queues", handlers.AdminListArcherQueues(db))
		})
	})
}

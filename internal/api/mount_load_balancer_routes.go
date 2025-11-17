package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountLoadBalancerRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/load-balancers", func(l chi.Router) {
		l.Use(authOrg)
		l.Get("/", handlers.ListLoadBalancers(db))
		l.Post("/", handlers.CreateLoadBalancer(db))
		l.Get("/{id}", handlers.GetLoadBalancer(db))
		l.Patch("/{id}", handlers.UpdateLoadBalancer(db))
		l.Delete("/{id}", handlers.DeleteLoadBalancer(db))
	})
}

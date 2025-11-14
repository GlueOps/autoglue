package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountNodePoolRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/node-pools", func(n chi.Router) {
		n.Use(authOrg)
		n.Get("/", handlers.ListNodePools(db))
		n.Post("/", handlers.CreateNodePool(db))
		n.Get("/{id}", handlers.GetNodePool(db))
		n.Patch("/{id}", handlers.UpdateNodePool(db))
		n.Delete("/{id}", handlers.DeleteNodePool(db))

		// Servers
		n.Get("/{id}/servers", handlers.ListNodePoolServers(db))
		n.Post("/{id}/servers", handlers.AttachNodePoolServers(db))
		n.Delete("/{id}/servers/{serverId}", handlers.DetachNodePoolServer(db))

		// Taints
		n.Get("/{id}/taints", handlers.ListNodePoolTaints(db))
		n.Post("/{id}/taints", handlers.AttachNodePoolTaints(db))
		n.Delete("/{id}/taints/{taintId}", handlers.DetachNodePoolTaint(db))

		// Labels
		n.Get("/{id}/labels", handlers.ListNodePoolLabels(db))
		n.Post("/{id}/labels", handlers.AttachNodePoolLabels(db))
		n.Delete("/{id}/labels/{labelId}", handlers.DetachNodePoolLabel(db))

		// Annotations
		n.Get("/{id}/annotations", handlers.ListNodePoolAnnotations(db))
		n.Post("/{id}/annotations", handlers.AttachNodePoolAnnotations(db))
		n.Delete("/{id}/annotations/{annotationId}", handlers.DetachNodePoolAnnotation(db))
	})
}

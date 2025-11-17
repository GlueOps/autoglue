package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountClusterRoutes(r chi.Router, db *gorm.DB, authOrg func(http.Handler) http.Handler) {
	r.Route("/clusters", func(c chi.Router) {
		c.Use(authOrg)
		c.Get("/", handlers.ListClusters(db))
		c.Post("/", handlers.CreateCluster(db))

		c.Get("/{clusterID}", handlers.GetCluster(db))
		c.Patch("/{clusterID}", handlers.UpdateCluster(db))
		c.Delete("/{clusterID}", handlers.DeleteCluster(db))

		c.Post("/{clusterID}/captain-domain", handlers.AttachCaptainDomain(db))
		c.Delete("/{clusterID}/captain-domain", handlers.DetachCaptainDomain(db))

		c.Post("/{clusterID}/control-plane-record-set", handlers.AttachControlPlaneRecordSet(db))
		c.Delete("/{clusterID}/control-plane-record-set", handlers.DetachControlPlaneRecordSet(db))

		c.Post("/{clusterID}/apps-load-balancer", handlers.AttachAppsLoadBalancer(db))
		c.Delete("/{clusterID}/apps-load-balancer", handlers.DetachAppsLoadBalancer(db))
		c.Post("/{clusterID}/glueops-load-balancer", handlers.AttachGlueOpsLoadBalancer(db))
		c.Delete("/{clusterID}/glueops-load-balancer", handlers.DetachGlueOpsLoadBalancer(db))

		c.Post("/{clusterID}/bastion", handlers.AttachBastionServer(db))
		c.Delete("/{clusterID}/bastion", handlers.DetachBastionServer(db))

		c.Post("/{clusterID}/kubeconfig", handlers.SetClusterKubeconfig(db))
		c.Delete("/{clusterID}/kubeconfig", handlers.ClearClusterKubeconfig(db))

	})
}

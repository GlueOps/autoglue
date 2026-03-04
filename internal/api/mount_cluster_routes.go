package api

import (
	"net/http"

	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func mountClusterRoutes(r chi.Router, db *gorm.DB, cfg config.Config, jobs *bg.Jobs, authOrg func(http.Handler) http.Handler) {
	r.Route("/clusters", func(c chi.Router) {
		c.Use(authOrg)
		c.Get("/", handlers.ListClusters(db, cfg))
		c.Post("/", handlers.CreateCluster(db, cfg))

		c.Get("/{clusterID}", handlers.GetCluster(db, cfg))
		c.Patch("/{clusterID}", handlers.UpdateCluster(db, cfg))
		c.Delete("/{clusterID}", handlers.DeleteCluster(db))

		c.Post("/{clusterID}/captain-domain", handlers.AttachCaptainDomain(db, cfg))
		c.Delete("/{clusterID}/captain-domain", handlers.DetachCaptainDomain(db, cfg))

		c.Post("/{clusterID}/control-plane-record-set", handlers.AttachControlPlaneRecordSet(db, cfg))
		c.Delete("/{clusterID}/control-plane-record-set", handlers.DetachControlPlaneRecordSet(db, cfg))

		c.Post("/{clusterID}/apps-load-balancer", handlers.AttachAppsLoadBalancer(db, cfg))
		c.Delete("/{clusterID}/apps-load-balancer", handlers.DetachAppsLoadBalancer(db, cfg))
		c.Post("/{clusterID}/glueops-load-balancer", handlers.AttachGlueOpsLoadBalancer(db, cfg))
		c.Delete("/{clusterID}/glueops-load-balancer", handlers.DetachGlueOpsLoadBalancer(db, cfg))

		c.Post("/{clusterID}/bastion", handlers.AttachBastionServer(db, cfg))
		c.Delete("/{clusterID}/bastion", handlers.DetachBastionServer(db, cfg))

		c.Post("/{clusterID}/kubeconfig", handlers.SetClusterKubeconfig(db, cfg))
		c.Delete("/{clusterID}/kubeconfig", handlers.ClearClusterKubeconfig(db, cfg))

		c.Post("/{clusterID}/node-pools", handlers.AttachNodePool(db, cfg))
		c.Delete("/{clusterID}/node-pools/{nodePoolID}", handlers.DetachNodePool(db, cfg))

		c.Get("/{clusterID}/runs", handlers.ListClusterRuns(db))
		c.Get("/{clusterID}/runs/{runID}", handlers.GetClusterRun(db))
		c.Post("/{clusterID}/actions/{actionID}/runs", handlers.RunClusterAction(db, jobs))
	})
}

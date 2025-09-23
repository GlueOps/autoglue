package api

import (
	httpPprof "net/http/pprof"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/handlers/annotations"
	"github.com/glueops/autoglue/internal/handlers/authn"
	"github.com/glueops/autoglue/internal/handlers/clusters"
	"github.com/glueops/autoglue/internal/handlers/health"
	"github.com/glueops/autoglue/internal/handlers/jobs"
	"github.com/glueops/autoglue/internal/handlers/labels"
	"github.com/glueops/autoglue/internal/handlers/nodepools"
	"github.com/glueops/autoglue/internal/handlers/orgs"
	"github.com/glueops/autoglue/internal/handlers/servers"
	"github.com/glueops/autoglue/internal/handlers/ssh"
	"github.com/glueops/autoglue/internal/handlers/taints"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/ui"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/api", func(api chi.Router) {
		api.Get("/healthz", health.Check)

		api.Route("/v1", func(v1 chi.Router) {
			secret := viper.GetString("authentication.jwt_secret")
			authMW := middleware.AuthMiddleware(secret)

			v1.Route("/admin", func(ad chi.Router) {
				ad.Use(authMW)
				ad.Get("/users", authn.AdminListUsers)
				ad.Post("/users", authn.AdminCreateUser)
				ad.Patch("/users/{userId}", authn.AdminUpdateUser)
				ad.Delete("/users/{userId}", authn.AdminDeleteUser)
			})

			v1.Route("/jobs", func(j chi.Router) {
				j.Use(authMW)
				j.Get("/kpi", jobs.GetKPI)
				j.Get("/queues", jobs.GetQueues)
				j.Get("/active", jobs.GetActive)
				j.Get("/failures", jobs.GetFailures)
				j.Post("/{id}/retry", jobs.RetryNow)
				j.Post("/{id}/cancel", jobs.Cancel)
				j.Post("/{id}/enqueue", jobs.Enqueue)
			})

			v1.Route("/auth", func(a chi.Router) {
				a.Post("/login", authn.Login)
				a.Post("/register", authn.Register)
				a.Post("/introspect", authn.Introspect)
				a.Post("/password/forgot", authn.RequestPasswordReset)
				a.Post("/password/reset", authn.ConfirmPasswordReset)
				a.Get("/verify", authn.VerifyEmail)
				a.Post("/verify/resend", authn.ResendVerification)

				a.Group(func(pr chi.Router) {
					pr.Use(authMW)
					pr.Post("/refresh", authn.Refresh)
					pr.Post("/logout", authn.Logout)
					pr.Post("/logout_all", authn.LogoutAll)
					pr.Get("/me", authn.Me)
					pr.Post("/password/change", authn.ChangePassword)
					pr.Post("/refresh/rotate", authn.RotateRefreshToken)
				})
			})

			v1.Route("/annotations", func(a chi.Router) {
				a.Use(authMW)
				a.Get("/", annotations.ListAnnotations)
				a.Post("/", annotations.CreateAnnotation)
				a.Get("/{id}", annotations.GetAnnotation)
				a.Patch("/{id}", annotations.UpdateAnnotation)
				a.Delete("/{id}", annotations.DeleteAnnotation)
				a.Get("/{id}/node_pools", annotations.ListNodePoolsWithAnnotation)
				a.Post("/{id}/node_pools", annotations.AddAnnotationToNodePools)
				a.Delete("/{id}/node_pools/{poolId}", annotations.RemoveAnnotationFromNodePool)
			})

			v1.Route("/orgs", func(o chi.Router) {
				o.Use(authMW)
				o.Post("/", orgs.CreateOrganization)
				o.Get("/", orgs.ListOrganizations)
				o.Post("/invite", orgs.InviteMember)
				o.Get("/members", orgs.ListMembers)
				o.Delete("/members/{userId}", orgs.DeleteMember)
				o.Patch("/{orgId}", orgs.UpdateOrganization)
				o.Delete("/{orgId}", orgs.DeleteOrganization)
			})

			v1.Route("/ssh", func(s chi.Router) {
				s.Use(authMW)
				s.Get("/", ssh.ListPublicKeys)
				s.Post("/", ssh.CreateSSHKey)
				s.Get("/{id}", ssh.GetSSHKey)
				s.Delete("/{id}", ssh.DeleteSSHKey)
				s.Get("/{id}/download", ssh.DownloadSSHKey)
			})

			v1.Route("/servers", func(s chi.Router) {
				s.Use(authMW)
				s.Get("/", servers.ListServers)
				s.Post("/", servers.CreateServer)
				s.Get("/{id}", servers.GetServer)
				s.Patch("/{id}", servers.UpdateServer)
				s.Delete("/{id}", servers.DeleteServer)
			})

			v1.Route("/node-pools", func(np chi.Router) {
				np.Use(authMW)
				np.Get("/", nodepools.ListNodePools)
				np.Post("/", nodepools.CreateNodePool)
				np.Get("/{id}", nodepools.GetNodePool)
				np.Patch("/{id}", nodepools.UpdateNodePool)
				np.Delete("/{id}", nodepools.DeleteNodePool)

				// servers
				np.Get("/{id}/servers", nodepools.ListNodePoolServers)
				np.Post("/{id}/servers", nodepools.AttachNodePoolServers)
				np.Delete("/{id}/servers/{serverId}", nodepools.DetachNodePoolServer)

				// taints
				np.Get("/{id}/taints", nodepools.ListNodePoolTaints)
				np.Post("/{id}/taints", nodepools.AttachNodePoolTaints)
				np.Delete("/{id}/taints/{taintId}", nodepools.DetachNodePoolTaint)

				// labels
				np.Get("/{id}/labels", nodepools.ListNodePoolLabels)
				np.Post("/{id}/labels", nodepools.AttachNodePoolLabels)
				np.Delete("/{id}/labels/{labelId}", nodepools.DetachNodePoolLabel)

				// annotations
				np.Get("/{id}/annotations", nodepools.ListNodePoolAnnotations)
				np.Post("/{id}/annotations", nodepools.AttachNodePoolAnnotations)
				np.Delete("/{id}/annotations/{annotationId}", nodepools.DetachNodePoolAnnotation)
			})

			v1.Route("/taints", func(t chi.Router) {
				t.Use(authMW)
				t.Get("/", taints.ListTaints)
				t.Post("/", taints.CreateTaint)
				t.Get("/{id}", taints.GetTaint)
				t.Patch("/{id}", taints.UpdateTaint)
				t.Delete("/{id}", taints.DeleteTaint)
				t.Post("/{id}/node_pools", taints.AddTaintToNodePool)
				t.Get("/{id}/node_pools", taints.ListNodePoolsWithTaint)
				t.Delete("/{id}/node_pools/{poolId}", taints.RemoveTaintFromNodePool)
			})

			v1.Route("/labels", func(l chi.Router) {
				l.Use(authMW)
				l.Get("/", labels.ListLabels)
				l.Post("/", labels.CreateLabel)
				l.Get("/{id}", labels.GetLabel)
				l.Patch("/{id}", labels.UpdateLabel)
				l.Delete("/{id}", labels.DeleteLabel)
				l.Get("/{id}/node_pools", labels.ListNodePoolsWithLabel)
				l.Post("/{id}/node_pools", labels.AddLabelToNodePool)
				l.Delete("/{id}/node_pools/{poolId}", labels.RemoveLabelFromNodePool)
			})

			v1.Route("/clusters", func(c chi.Router) {
				c.Use(authMW)
				c.Get("/", clusters.ListClusters)
				c.Post("/", clusters.CreateCluster)

				c.Get("/{id}", clusters.GetCluster)
				c.Patch("/{id}", clusters.UpdateCluster)
				c.Delete("/{id}", clusters.DeleteCluster)

				c.Get("/{id}/node_pools", clusters.ListClusterNodePools)
				c.Post("/{id}/node_pools", clusters.AttachNodePools)
				c.Delete("/{id}/node_pools/{poolId}", clusters.DetachNodePool)

				c.Get("/{id}/bastion", clusters.GetBastion)
				c.Post("/{id}/bastion", clusters.PutBastion)
				c.Delete("/{id}/bastion", clusters.DeleteBastion)
			})
		})
	})

	r.Route("/debug/pprof", func(pr chi.Router) {
		pr.Get("/", httpPprof.Index)
		pr.Get("/cmdline", httpPprof.Cmdline)
		pr.Get("/profile", httpPprof.Profile)
		pr.Get("/symbol", httpPprof.Symbol)
		pr.Get("/trace", httpPprof.Trace)

		pr.Handle("/allocs", httpPprof.Handler("allocs"))
		pr.Handle("/block", httpPprof.Handler("block"))
		pr.Handle("/goroutine", httpPprof.Handler("goroutine"))
		pr.Handle("/heap", httpPprof.Handler("heap"))
		pr.Handle("/mutex", httpPprof.Handler("mutex"))
		pr.Handle("/threadcreate", httpPprof.Handler("threadcreate"))
	})

	if config.IsUIDev() {
		if h, err := ui.DevProxy("http://localhost:5173"); err == nil {
			r.NotFound(h.ServeHTTP)
		}
	} else {
		if h, err := ui.SPAHandler(); err == nil {
			r.NotFound(h.ServeHTTP)
		}
	}
}

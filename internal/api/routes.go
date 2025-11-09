package api

import (
	"fmt"
	"net/http"
	httpPprof "net/http/pprof"
	"os"
	"time"

	"github.com/glueops/autoglue/docs"
	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/handlers"
	"github.com/glueops/autoglue/internal/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	"gorm.io/gorm"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(db *gorm.DB, jobs *bg.Jobs) http.Handler {
	zerolog.TimeFieldFormat = time.RFC3339

	l := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"})
	log.Logger = l

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(zeroLogMiddleware())
	r.Use(middleware.Recoverer)
	r.Use(SecurityHeaders)
	r.Use(requestBodyLimit(10 << 20))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	allowed := getAllowedOrigins()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowed,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Org-ID",
			"X-API-KEY",
			"X-ORG-KEY",
			"X-ORG-SECRET",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           600,
	}))

	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/.well-known/jwks.json", handlers.JWKSHandler)
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			authUser := httpmiddleware.AuthMiddleware(db, false)
			authOrg := httpmiddleware.AuthMiddleware(db, true)

			// Also serving a versioned JWKS for swagger, which uses BasePath
			v1.Get("/.well-known/jwks.json", handlers.JWKSHandler)

			v1.Get("/healthz", handlers.HealthCheck)
			v1.Get("/version", handlers.Version)

			v1.Route("/auth", func(a chi.Router) {
				a.Post("/{provider}/start", handlers.AuthStart(db))
				a.Get("/{provider}/callback", handlers.AuthCallback(db))
				a.Post("/refresh", handlers.Refresh(db))
				a.Post("/logout", handlers.Logout(db))
			})

			v1.Route("/admin", func(admin chi.Router) {
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

			v1.Route("/me", func(me chi.Router) {
				me.Use(authUser)

				me.Get("/", handlers.GetMe(db))
				me.Patch("/", handlers.UpdateMe(db))

				me.Get("/api-keys", handlers.ListUserAPIKeys(db))
				me.Post("/api-keys", handlers.CreateUserAPIKey(db))
				me.Delete("/api-keys/{id}", handlers.DeleteUserAPIKey(db))
			})

			v1.Route("/orgs", func(o chi.Router) {
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

			v1.Route("/credentials", func(c chi.Router) {
				c.Use(authOrg)
				c.Get("/", handlers.ListCredentials(db))
				c.Post("/", handlers.CreateCredential(db))
				c.Get("/{id}", handlers.GetCredential(db))
				c.Patch("/{id}", handlers.UpdateCredential(db))
				c.Delete("/{id}", handlers.DeleteCredential(db))
				c.Post("/{id}/reveal", handlers.RevealCredential(db))
			})

			v1.Route("/ssh", func(s chi.Router) {
				s.Use(authOrg)
				s.Get("/", handlers.ListPublicSshKeys(db))
				s.Post("/", handlers.CreateSSHKey(db))
				s.Get("/{id}", handlers.GetSSHKey(db))
				s.Delete("/{id}", handlers.DeleteSSHKey(db))
				s.Get("/{id}/download", handlers.DownloadSSHKey(db))
			})

			v1.Route("/servers", func(s chi.Router) {
				s.Use(authOrg)
				s.Get("/", handlers.ListServers(db))
				s.Post("/", handlers.CreateServer(db))
				s.Get("/{id}", handlers.GetServer(db))
				s.Patch("/{id}", handlers.UpdateServer(db))
				s.Delete("/{id}", handlers.DeleteServer(db))
			})

			v1.Route("/taints", func(s chi.Router) {
				s.Use(authOrg)
				s.Get("/", handlers.ListTaints(db))
				s.Post("/", handlers.CreateTaint(db))
				s.Get("/{id}", handlers.GetTaint(db))
				s.Patch("/{id}", handlers.UpdateTaint(db))
				s.Delete("/{id}", handlers.DeleteTaint(db))
			})

			v1.Route("/labels", func(l chi.Router) {
				l.Use(authOrg)
				l.Get("/", handlers.ListLabels(db))
				l.Post("/", handlers.CreateLabel(db))
				l.Get("/{id}", handlers.GetLabel(db))
				l.Patch("/{id}", handlers.UpdateLabel(db))
				l.Delete("/{id}", handlers.DeleteLabel(db))
			})

			v1.Route("/annotations", func(a chi.Router) {
				a.Use(authOrg)
				a.Get("/", handlers.ListAnnotations(db))
				a.Post("/", handlers.CreateAnnotation(db))
				a.Get("/{id}", handlers.GetAnnotation(db))
				a.Patch("/{id}", handlers.UpdateAnnotation(db))
				a.Delete("/{id}", handlers.DeleteAnnotation(db))
			})

			v1.Route("/node-pools", func(n chi.Router) {
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
		})
	})
	if config.IsDebug() {
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
	}

	if config.IsSwaggerEnabled() {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("swagger.json"),
		))
		r.Get("/swagger/swagger.json", serveSwaggerFromEmbed(docs.SwaggerJSON, "application/json"))
		r.Get("/swagger/swagger.yaml", serveSwaggerFromEmbed(docs.SwaggerYAML, "application/x-yaml"))
	}

	if config.IsUIDev() {
		fmt.Println("Running in development mode")
		// Dev: isolate proxy from chi middlewares so WS upgrade can hijack.
		proxy, err := web.DevProxy("http://localhost:5173")
		if err != nil {
			log.Error().Err(err).Msg("dev proxy init failed")
			return r // fallback
		}

		mux := http.NewServeMux()
		// Send API/Swagger/pprof to chi
		mux.Handle("/api/", r)
		mux.Handle("/api", r)
		mux.Handle("/swagger/", r)
		mux.Handle("/debug/pprof/", r)
		// Everything else (/, /brand-preview, assets) → proxy (no middlewares)
		mux.Handle("/", proxy)

		return mux
	} else {
		fmt.Println("Running in production mode")
		if h, err := web.SPAHandler(); err == nil {
			r.NotFound(h.ServeHTTP)
		} else {
			log.Error().Err(err).Msg("spa handler init failed")
		}
	}

	return r
}

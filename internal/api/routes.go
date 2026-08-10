package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
)

// RouterOpts carries the optional handlers the API mounts alongside the
// versioned routes. Both are nil unless explicitly enabled.
type RouterOpts struct {
	// Studio is the pgweb handler mounted at /db-studio.
	Studio http.Handler
	// RiverUI is the River dashboard mounted at RiverUIPrefix.
	RiverUI http.Handler
}

func NewRouter(db *gorm.DB, jobs *bg.Client, cfg config.Config, opts RouterOpts) http.Handler {
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
	r.Use(httprate.LimitByIP(1000, 1*time.Minute))
	r.Use(middleware.StripSlashes)

	allowed := config.AllowedOrigins()
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

	r.Use(middleware.Maybe(
		middleware.AllowContentType("application/json"),
		func(r *http.Request) bool {
			// return true  => run AllowContentType
			// return false => skip AllowContentType for this request
			return !strings.HasPrefix(r.URL.Path, "/db-studio")
		}))
	//r.Use(middleware.AllowContentType("application/json"))

	// Unversioned, non-auth endpoints
	r.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	// Versioned API
	mountAPIRoutes(r, db, cfg, jobs)

	// Optional DB studio
	if opts.Studio != nil {
		r.Group(func(gr chi.Router) {
			authUser := httpmiddleware.AuthMiddleware(db, false)
			adminOnly := httpmiddleware.RequirePlatformAdmin()
			gr.Use(authUser, adminOnly)
			gr.Mount("/db-studio", opts.Studio)
		})
	}

	// River dashboard. Same admin gate as the DB studio: AuthMiddleware falls
	// back to the ag_jwt cookie, so a plain browser navigation authenticates.
	if opts.RiverUI != nil {
		r.Group(func(gr chi.Router) {
			authUser := httpmiddleware.AuthMiddleware(db, false)
			adminOnly := httpmiddleware.RequirePlatformAdmin()
			gr.Use(authUser, adminOnly)
			gr.Mount(RiverUIPrefix, opts.RiverUI)
		})
	}

	// pprof
	if config.IsDebug() {
		mountPprofRoutes(r)
	}

	// Swagger
	if config.IsSwaggerEnabled() {
		mountSwaggerRoutes(r)
	}

	// UI dev/prod
	if config.IsUIDev() {
		fmt.Println("Running in development mode")
		proxy, err := web.DevProxy("http://localhost:5173")
		if err != nil {
			log.Error().Err(err).Msg("dev proxy init failed")
			return r // fallback
		}

		mux := http.NewServeMux()
		mux.Handle("/api/", r)
		mux.Handle("/api", r)
		mux.Handle("/swagger", r)
		mux.Handle("/swagger/", r)
		mux.Handle("/db-studio/", r)
		mux.Handle(RiverUIPrefix+"/", r)
		mux.Handle("/debug/pprof/", r)
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

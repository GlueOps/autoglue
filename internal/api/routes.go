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

func NewRouter(db *gorm.DB, jobs *bg.Jobs, studio http.Handler) http.Handler {
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
	r.Use(middleware.StripSlashes)

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
	mountAPIRoutes(r, db, jobs)

	// Optional DB studio
	if studio != nil {
		r.Group(func(gr chi.Router) {
			authUser := httpmiddleware.AuthMiddleware(db, false)
			adminOnly := httpmiddleware.RequirePlatformAdmin()
			gr.Use(authUser, adminOnly)
			gr.Mount("/db-studio", studio)
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

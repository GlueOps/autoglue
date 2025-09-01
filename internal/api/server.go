package api

import (
	"net/http"
	"time"

	"github.com/glueops/autoglue/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "content-type", "Authorization", "authorization", "X-Org-ID", "x-org-id"},
		AllowCredentials: true,
		// OptionsPassthrough: false, // default; Chi will auto 200 OPTIONS
		// MaxAge: 300,               // optional
	}))

	RegisterRoutes(r)

	r.Mount("/swagger", httpSwagger.WrapHandler)
	r.Get("/swagger/swagger.json", serveSwaggerFromEmbed(docs.SwaggerJSON, "application/json"))
	r.Get("/swagger/swagger.yaml", serveSwaggerFromEmbed(docs.SwaggerYAML, "application/x-yaml"))
	return r
}

func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      NewRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func serveSwaggerFromEmbed(data []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

package api

import (
	"github.com/glueops/autoglue/api/handlers/health"
	"github.com/glueops/autoglue/api/handlers/pprof"
	"github.com/glueops/autoglue/docs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s %s", r.Method, r.URL.String(), r.RemoteAddr, r.Proto)
		next.ServeHTTP(w, r)
		log.Printf("completed in %v", time.Since(start))
	})
}

func NewRouter() http.Handler {
	router := mux.NewRouter()
	router.UseEncodedPath()
	router.StrictSlash(true)

	// Middleware
	router.Use(LoggingMiddleware)

	// Routes

	router.HandleFunc("/healthz", health.Check).Methods("GET")

	RegisterRoutes(router)
	pprof.RegisterPprofRoutes(router)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/swagger/swagger.json", serveSwaggerFromEmbed(docs.SwaggerJSON, "application/json"))
	router.HandleFunc("/swagger/swagger.yaml", serveSwaggerFromEmbed(docs.SwaggerYAML, "application/x-yaml"))

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Org-ID",
		}),
	)

	router.PathPrefix("/").Handler(StaticHandler())

	return corsHandler(router)
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
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

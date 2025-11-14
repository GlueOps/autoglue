package api

import (
	"github.com/glueops/autoglue/docs"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func mountSwaggerRoutes(r chi.Router) {
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger.json"),
	))
	r.Get("/swagger/swagger.json", serveSwaggerFromEmbed(docs.SwaggerJSON, "application/json"))
	r.Get("/swagger/swagger.yaml", serveSwaggerFromEmbed(docs.SwaggerYAML, "application/x-yaml"))
}

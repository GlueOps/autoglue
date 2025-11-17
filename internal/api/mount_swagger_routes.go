package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/glueops/autoglue/docs"
	"github.com/go-chi/chi/v5"
)

func mountSwaggerRoutes(r chi.Router) {
	r.Get("/swagger", RapidDocHandler("/swagger/swagger.yaml"))
	r.Get("/swagger/index.html", RapidDocHandler("/swagger/swagger.yaml"))
	r.Get("/swagger/swagger.json", serveSwaggerFromEmbed(docs.SwaggerJSON, "application/json"))
	r.Get("/swagger/swagger.yaml", serveSwaggerFromEmbed(docs.SwaggerYAML, "application/x-yaml"))
}

var rapidDocTmpl = template.Must(template.New("redoc").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>AutoGlue API Docs</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta charset="utf-8">
  <style>
    body { margin: 0; padding: 0; }
    .redoc-container { height: 100vh; }
  </style>
</head>
<body>
  <rapi-doc
    id="autoglue-docs"
    spec-url="{{.SpecURL}}"
    render-style="read"
    theme="dark"
    show-header="false"
    persist-auth="true"
	allow-advanced-search="true"
	schema-description-expanded="true"                                                                                                                                   
	allow-schema-description-expand-toggle="false"                                                                                                                       
	allow-spec-file-download="true"                                                                                                                                      
	allow-spec-file-load="false"                                                                                                                                         
	allow-spec-url-load="false"                       
    allow-try="true"
    schema-style="tree"
	fetch-credentials="include"
    default-api-server="{{.DefaultServer}}"
 	api-key-name="X-ORG-ID"
  	api-key-location="header"
  	api-key-value=""
  />
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
  <script>
    window.addEventListener('DOMContentLoaded', () => {
      const rd = document.getElementById('autoglue-docs');
      if (!rd) return;

      const storedOrg = localStorage.getItem('autoglue.org');
      if (storedOrg) {
        rd.setAttribute('api-key-value', storedOrg);
      }
    }
  </script>
</body>
</html>`))

func RapidDocHandler(specURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		host := r.Host
		defaultServer := fmt.Sprintf("%s://%s/api/v1", scheme, host)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := rapidDocTmpl.Execute(w, map[string]string{
			"SpecURL":       specURL,
			"DefaultServer": defaultServer,
		}); err != nil {
			http.Error(w, "failed to render docs", http.StatusInternalServerError)
			return
		}
	}
}

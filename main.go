package main

import (
	"os"
	"strings"

	"github.com/glueops/autoglue/cmd/cli"
	"github.com/glueops/autoglue/docs"
	"github.com/glueops/autoglue/internal/config"
)

// @title       AutoGlue API
// @version     1.0
// @description API for managing K3s clusters across cloud providers
// @BasePath    /
// @schemes     https http
// @host        autoglue.apps.nonprod.earth.onglueops.rocks

// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
func main() {
	if h := os.Getenv("SWAG_HOST"); h != "" {
		docs.SwaggerInfo.Host = h
	}
	if s := os.Getenv("SWAG_SCHEMES"); s != "" {
		// e.g. "http,https" or "https"
		docs.SwaggerInfo.Schemes = splitCSV(s)
	}

	config.Load()
	cli.Execute()
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

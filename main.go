package main

import (
	"github.com/glueops/autoglue/cmd"
)

// @title AutoGlue API
// @version 1.0
// @description API for managing K3s clusters across cloud providers
// @contact.name GlueOps

// @servers.url https://autoglue.onglueops.rocks/api/v1
// @servers.description Production API
// @servers.url https://autoglue.apps.nonprod.earth.onglueops.rocks/api/v1
// @servers.description Staging API
// @servers.url http://localhost:8080/api/v1
// @servers.description Local dev

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
// @description	User API key

// @securityDefinitions.apikey OrgKeyAuth
// @in header
// @name X-ORG-KEY
// @description Org-level key/secret authentication

// @securityDefinitions.apikey OrgSecretAuth
// @in header
// @name X-ORG-SECRET
// @description	Org-level secret

func main() {
	cmd.Execute()
}

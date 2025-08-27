package main

import (
	"github.com/glueops/autoglue/cmd"
	"github.com/glueops/autoglue/internal/config"
)

// @title AutoGlue API
// @version 1.0
// @description API for managing K3s clusters across cloud providers
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.Load()
	cmd.Execute()
}

package assets

import "embed"

//go:embed "deployment"
var EmbeddedFS embed.FS

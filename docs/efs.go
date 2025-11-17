package docs

import _ "embed"

//go:embed openapi.json
var SwaggerJSON []byte

//go:embed openapi.yaml
var SwaggerYAML []byte

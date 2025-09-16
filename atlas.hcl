data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./models",
    "--dialect", "postgresql",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgresql://dev"
}
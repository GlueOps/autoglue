terraform {
  required_version = ">= 1.5.0"

  required_providers {
    autoglue = {
      source = "glueops/autoglue/autoglue"
    }
    http = {
      source = "hashicorp/http"
    }
    local = {
      source = "hashicorp/local"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

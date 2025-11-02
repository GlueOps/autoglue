terraform {
  required_version = ">= 1.5.0"

  required_providers {
    autoglue = {
      source  = "glueops/autoglue/autoglue"
      version = "0.0.1" # matches your dev install VER
    }
    http = {
      source  = "hashicorp/http"
      version = ">= 3.4.0"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 2.5.1"
    }
  }
}
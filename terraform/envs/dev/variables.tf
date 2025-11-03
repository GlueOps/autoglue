variable "addr" {
  description = "Base URL to the Autoglue API, e.g. http://localhost:8080/api/v1"
  type        = string
  default     = "http://localhost:8080/api/v1"
}

variable "org_key" {
  description = "Org key for machine auth (sent as X-ORG-KEY)"
  type        = string
  sensitive   = true
}

variable "org_secret" {
  description = "Org secret for machine auth (sent as X-ORG-SECRET)"
  type        = string
  sensitive   = true
}

variable "ssh_keys" {
  description = "Map of SSH key specs"
  type = map(object({
    name            = string
    comment         = string
    type            = string
    bits            = optional(number)
    enable_download = optional(bool, true)
    download_part   = optional(string, "both")
    download_dir    = optional(string, "out")
  }))
}


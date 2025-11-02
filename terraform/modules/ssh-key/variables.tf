variable "addr" {
  type = string
}

variable "org_key"    {
  type = string
  sensitive = true
  default = null
}

variable "org_secret" {
  type = string
  sensitive = true
  default = null
}

variable "name" {
  type = string
}

variable "comment" {
  type = string
}

variable "type" {
  type = string
}

variable "enable_download" {
  type    = bool
  default = false
}

variable "download_part" {
  type    = string
  default = "both"
}

variable "download_dir" {
  type    = string
  default = "ssh_artifacts"
}

variable "bits" {
  type    = number
  default = null  # null for ed25519
}
variable "servers" {
  description = <<-EOT
    Map of servers to create. Example shape:
    {
      bastion = {
        hostname           = "bastion-01"
        private_ip_address = "10.0.0.10"
        public_ip_address  = "54.12.34.56" # required when role = "bastion"
        role               = "bastion"
        ssh_user           = "ubuntu"
        ssh_key_ref        = "bastionKey"  # OR set ssh_key_id instead
        # ssh_key_id       = "uuid-string"
        # status           = "pending|provisioning|ready|failed"
      }
    }
  EOT
  type = map(object({
    hostname           = optional(string)
    private_ip_address = string
    public_ip_address  = optional(string)
    role               = string
    ssh_user           = string
    ssh_key_ref        = optional(string) # name to look up in var.ssh_key_ids
    ssh_key_id         = optional(string) # direct UUID (overrides ssh_key_ref if set)
    status             = optional(string) # pending|provisioning|ready|failed
  }))
  default = {}
}

variable "ssh_key_ids" {
  description = "Map of SSH key IDs you can reference via servers[*].ssh_key_ref."
  type        = map(string)
  default     = {}
}

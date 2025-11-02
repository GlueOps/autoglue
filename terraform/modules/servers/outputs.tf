output "ids" {
  description = "Map of server IDs by key."
  value       = { for k, r in autoglue_server.this : k => r.id }
}

output "statuses" {
  description = "Map of server statuses by key."
  value       = { for k, r in autoglue_server.this : k => r.status }
}

output "details" {
  description = "Selected attributes for convenience."
  value = {
    for k, r in autoglue_server.this : k => {
      id                  = r.id
      organization_id     = r.organization_id
      hostname            = r.hostname
      private_ip_address  = r.private_ip_address
      public_ip_address   = r.public_ip_address
      role                = r.role
      ssh_user            = r.ssh_user
      ssh_key_id          = r.ssh_key_id
      status              = r.status
      created_at          = r.created_at
      updated_at          = r.updated_at
    }
  }
}

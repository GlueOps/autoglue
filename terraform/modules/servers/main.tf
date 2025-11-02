locals {
  # Resolve the SSH key ID for each server:
  # Prefer explicit ssh_key_id, otherwise look up by ssh_key_ref in var.ssh_key_ids.
  resolved_ssh_key_ids = {
    for name, spec in var.servers :
    name => coalesce(
      try(spec.ssh_key_id, null),
      try(var.ssh_key_ids[spec.ssh_key_ref], null)
    )
  }
}

resource "autoglue_server" "this" {
  for_each = var.servers

  hostname           = try(each.value.hostname, null)
  private_ip_address = each.value.private_ip_address
  public_ip_address  = try(each.value.public_ip_address, null)
  role               = lower(each.value.role)
  ssh_user           = each.value.ssh_user
  ssh_key_id         = local.resolved_ssh_key_ids[each.key]
  status             = try(each.value.status, null)

  # Client-side guards to match your API rules
  lifecycle {
    precondition {
      condition     = local.resolved_ssh_key_ids[each.key] != null && local.resolved_ssh_key_ids[each.key] != ""
      error_message = "Provide either ssh_key_id or ssh_key_ref (and pass ssh_key_ids to the module)."
    }
    precondition {
      condition     = lower(each.value.role) != "bastion" ? true : (try(each.value.public_ip_address, "") != "")
      error_message = "public_ip_address is required when role == \"bastion\"."
    }
    precondition {
      condition     = try(each.value.status, "") == "" || contains(["pending", "provisioning", "ready", "failed"], lower(each.value.status))
      error_message = "status must be one of: pending, provisioning, ready, failed (or omitted)."
    }
  }
}

# Use the module N times with for_each
module "ssh" {
  source   = "../../modules/ssh-key"
  for_each = var.ssh_keys

  # Pass through inputs
  addr            = var.addr                     # used by HTTP download URL
  name            = each.value.name
  comment         = each.value.comment
  type            = each.value.type
  bits            = try(each.value.bits, null)
  enable_download = try(each.value.enable_download, true)
  download_part   = try(each.value.download_part, "both")
  download_dir    = try(each.value.download_dir, "out/${each.key}")

  org_key = var.org_key
  org_secret = var.org_secret
}

# Example: aggregate outputs by key
output "ssh_ids" {
  value = { for k, m in module.ssh : k => m.id }
}
output "ssh_public_keys" {
  value = { for k, m in module.ssh : k => m.public_key }
}
output "ssh_written_files" {
  value = { for k, m in module.ssh : k => m.written_files }
}

module "servers" {
  source = "../../modules/servers"
  # Wire the SSH key IDs so servers can reference them by name
  ssh_key_ids = { for k, m in module.ssh : k => m.id }

  servers = {
    bastion = {
      hostname           = "bastion-01"
      private_ip_address = "10.0.0.10"
      public_ip_address = "65.109.95.175" # required for role=bastion
      role               = "bastion"
      ssh_user           = "root"
      ssh_key_ref = "bastionKey"  # points to module.ssh["bastionKey"].id
      status             = "pending"
    }

    manager1 = {
      hostname           = "k3s-mgr-01"
      private_ip_address = "10.0.1.11"
      role               = "master"
      ssh_user           = "ubuntu"
      ssh_key_ref        = "clusterKey"
      status             = "pending"
    }

    agent1 = {
      hostname           = "k3s-agent-01"
      private_ip_address = "10.0.2.21"
      role               = "worker"
      ssh_user           = "ubuntu"
      ssh_key_ref        = "clusterKey"
      status             = "pending"
    }
  }
}

output "server_ids" {
  value = module.servers.ids
}

output "server_statuses" {
  value = module.servers.statuses
}
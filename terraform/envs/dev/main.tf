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

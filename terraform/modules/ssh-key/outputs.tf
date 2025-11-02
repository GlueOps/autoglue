output "id"          { value = autoglue_ssh_key.this.id }
output "public_key"  { value = autoglue_ssh_key.this.public_key }
output "fingerprint" { value = autoglue_ssh_key.this.fingerprint }
output "created_at"  { value = autoglue_ssh_key.this.created_at }

output "written_files" {
  value = compact(concat(
      var.enable_download && var.download_part == "public"  ? [local_file.public_key[0].filename] : [],
      var.enable_download && var.download_part == "private" ? [local_sensitive_file.private_key[0].filename] : [],
      var.enable_download && var.download_part == "both"    ? [local_sensitive_file.zip[0].filename] : []
  ))
}

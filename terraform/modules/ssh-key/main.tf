locals { is_rsa = var.type == "rsa" }

# 1) Create key
resource "autoglue_ssh_key" "this" {
  name    = var.name
  comment = var.comment
  type    = var.type
  bits    = local.is_rsa ? var.bits : null
}

# 2) Optionally download via HTTP (mode=json)
data "http" "download" {
  count = var.enable_download ? 1 : 0

  url = "${var.addr}/ssh/${autoglue_ssh_key.this.id}/download?part=${var.download_part}&mode=json"

  # Inherit org_key/org_secret via provider headers — we’re not configuring http headers here
  # because your API auth for downloads is via X-ORG-KEY / X-ORG-SECRET.
  # If you require those headers here, add request_headers and pass them from root as inputs.
  # For org key/secret auth on download, uncomment and add module inputs:
  request_headers = {
    "X-ORG-KEY"    = var.org_key
    "X-ORG-SECRET" = var.org_secret
    "Accept"       = "application/json"
  }
}

locals {
  dl      = var.enable_download ? jsondecode(one(data.http.download[*].response_body)) : null
  zip_b64 = coalesce(try(local.dl.zipBase64, null), try(local.dl.zip_base64, null))
}

resource "null_resource" "mkdirs" {
  count = var.enable_download ? 1 : 0
  provisioner "local-exec" { command = "mkdir -p ${var.download_dir}" }
}

# public only
resource "local_file" "public_key" {
  count           = var.enable_download && var.download_part == "public" ? 1 : 0
  filename        = "${var.download_dir}/${try(local.dl.filenames[0], "id_rsa.pub")}"
  content         = try(local.dl.publicKey, "")
  file_permission = "0644"
  depends_on      = [null_resource.mkdirs]
}

# private only
resource "local_sensitive_file" "private_key" {
  count     = var.enable_download && var.download_part == "private" ? 1 : 0
  filename  = "${var.download_dir}/${try(local.dl.filenames[0], "id_rsa.pem")}"
  content   = try(local.dl.privatePEM, "")
  depends_on = [null_resource.mkdirs]
}

# both -> zip
resource "local_sensitive_file" "zip" {
  count          = var.enable_download && var.download_part == "both" ? 1 : 0
  filename       = "${var.download_dir}/${try(local.dl.filenames[0], "ssh_key.zip")}"
  content_base64 = local.zip_b64
  depends_on     = [null_resource.mkdirs]

  lifecycle {
    postcondition {
      condition     = length(try(local.zip_b64, "")) > 0
      error_message = "API did not return a zip payload for part=both."
    }
  }
}

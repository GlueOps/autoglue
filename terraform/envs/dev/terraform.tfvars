org_key = "org_lnJwmyyWH7JC-JgZo5v3Kw"
org_secret = "fqd9yebGMfK6h5HSgWn4sXrwr9xlFbvbIYtNylRElMQ"

ssh_keys = {
  key1 = {
    name            = "CI deploy key 1"
    comment         = "deploy1@autoglue"
    type            = "rsa"
    bits            = 4096
    enable_download = true
    download_part   = "both"
    download_dir    = "out/key1"
  }
  key2 = {
    name    = "CI deploy key 2"
    comment = "deploy2@autoglue"
    type    = "ed25519"           # bits ignored
    enable_download = true
    download_part   = "both"
    download_dir    = "out/key2"
  }
}


org_key = "org_lnJwmyyWH7JC-JgZo5v3Kw"
org_secret = "fqd9yebGMfK6h5HSgWn4sXrwr9xlFbvbIYtNylRElMQ"

ssh_keys = {
  bastionKey = {
    name            = "Bastion Key"
    comment         = "deploy@autoglue"
    type            = "rsa"
    bits            = 4096
    enable_download = true
    download_part   = "both"
    download_dir    = "out/bastionKey"
  }
  clusterKey = {
    name    = "Cluster Key"
    comment = "bastion@autoglue"
    type    = "ed25519"           # bits ignored
    enable_download = true
    download_part   = "both"
    download_dir    = "out/clusterKey"
  }
}
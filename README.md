# Autoglue

## Setup Env
create .env file:
```
DB_USER=autoglue
DB_PASSWORD=autoglue
DB_NAME=autoglue

BIND_ADDRESS=127.0.0.1
BIND_PORT=8080

DATABASE_URL=postgres://autoglue:autoglue@localhost:5432/autoglue

# Used to switch between dev proxy and spa handler
# dev proxy - proxy binds port 5173 (vite) to BIND_PORT at /
# spa handler - requires ui build and serves embedded from build result
UI_DEV=true

# in development - removes hsts in security middleware, any other value - hsts gets enabled
ENV=development

# Enables or disables /debug/pprof
DEBUG=true

# Enables or disables embedded swagger ui
SWAGGER=true

JWT_ISSUER=http://localhost:8080
JWT_AUDIENCE=autoglue-ui
# optional - 32 bytes base64url (no padding); use only if you want to encrypt private PEM in DB
# openssl rand -base64 32 | tr '+/' '-_' | tr -d '='
JWT_PRIVATE_ENC_KEY=IQp53NFkJvbfgLgZc-wo1CNLYRndJA-g5u13HFAn4rI

OAUTH_REDIRECT_BASE=http://localhost:8080
GOOGLE_CLIENT_ID=742299367132-na2i6spu7doukkqtk52vcqd6p94ijupi.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
#GITHUB_CLIENT_ID=...
#GITHUB_CLIENT_SECRET=...
```

## Bring up Database:
```bash
docker compose up -d
```

## Generate JWT tokens used for auth in the DB
Private Key is encrypted by `JWT_PRIVATE_ENC_KEY`
If this is not set, the private key is stored in plain text in the DB - its never exposed at any rate

```bash
go run . keys generate
```

## Generate master encryption key 
The master encryption key is used to generate the org encryption keys - without it there will be failures
```bash
go run . encrypt create-master
```

## Ensure your swagger and SDKs are up to date with the api spec
```bash
make swagger
make sdk-all
```

## Build terraform provider
Currently, the terraform provider consumes the go sdk through an alias
Once the SDK is published to its own repo, the alias can be removed. but this is why its imperative to ensure the SDK is up to date

The command below builds the provider binary, and installs it where terraform expects it to be placed if it were downloaded from a registry
```bash
cd terraform-provider-autoglue
make dev
```

## See UI & terraform in action
From the project root
## UI & API - required for the terraform
Start the API & UI (the env embeds it with a dev proxy)

if you witness a failure here, run `make ui`

This is most likely the SPA handler trying to embed ui files that dont exist
```bash
go run .
```

From your GLUEOPS profiled browser - http://localhost:8080
Login - this is restricted to glueops.dev at the minute (in google workspace settings - outside of the API)

Create your org (http://localhost:8080/me) - you should be redirected here after initial login

Once you have an org - create a set of api keys for your org:
They will be in the format of:
```text
Org Key: org_lnJwmyyWH7JC-JgZo5v3Kw
Org Secret: fqd9yebGMfK6h5HSgWn4sXrwr9xlFbvbIYtNylRElMQ
```

use them in terraform/envs/dev/terraform.tfvars

in my example here, i also create ssh keys in my example:
```terraform
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

```

explore `main.tf` for how the module ssh-keys module is used
also you will see there how to create servers using the servers module

in `terraform/envs/dev`
```bash
rm -rf .terraform*
tofu init -upgrade

tofu plan

tofu apply -auto-approve
```

If everything went to plan, you'll have an `out` directory containting 2 zip file - one for each of the ssh keys

In the UI you will also see the SSH Keys on its page,
you will also see the servers created on its page.

## <span style="color:red">WARNING</span>
<span style="color:red">!!!!Terraform destroy deletes the keys from the api as well as deletes the local files!!!!</span>
```bash
tofu destroy -auto-approve
```



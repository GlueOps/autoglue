# Autoglue

## Setup Env
create .env file:
```bash
cp .env.example .env
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
go run . serve
```

The API and the background workers are separate processes. `serve` enqueues
jobs but never runs them, so anything asynchronous — bastion bootstrap, DNS
reconcile, cluster actions, backups, key sweeps — stays queued until you also
start a worker in a second terminal:

```bash
go run . worker
```

Bare `go run .` is still an alias for `serve`, so on its own it gives you an API
whose jobs never execute. If a cluster action sits in `queued` forever, this is
why.

Multiple workers are safe to run at once: River leases each job to exactly one
of them, and the periodic schedule runs only on the elected leader.

From your GLUEOPS profiled browser - http://localhost:8080
Login - this is restricted to glueops.dev at the minute (in google workspace settings - outside of the API)

Create your org (http://localhost:8080/me) - you should be redirected here after initial login

## Background jobs

Jobs run on [River](https://riverqueue.com), backed by the same Postgres
database. The schema (`river_job` and friends) is migrated automatically at
startup by either process.

Queues:

| Queue         | Concurrency | Work                                            |
| ------------- | ----------- | ----------------------------------------------- |
| `clusters`    | 30          | `cluster_action`, `bootstrap_bastion`           |
| `maintenance` | 2           | `dns_reconcile`, `db_backup_s3`, `org_key_sweeper`, `tokens_cleanup`, `job_logs_cleanup`, `vacuum` |
| `default`     | 10          | unassigned work                                 |

Long-running cluster work is kept off `maintenance` so a multi-hour bootstrap
cannot starve the hourly sweepers.

River's own dashboard is mounted at `/admin/river/` behind the platform-admin
gate, and replaces the old hand-rolled jobs admin page. Retention of finished
jobs is handled by River itself (`river.completed_retain_days` and friends in
config), not by a cleanup job.

### Housekeeping schedule

| Job                | When            | What it does                                                  |
| ------------------ | --------------- | ------------------------------------------------------------- |
| `tokens_cleanup`   | daily 03:45     | Deletes revoked and expired `refresh_tokens`                   |
| `job_logs_cleanup` | daily 04:15     | Prunes `job_logs` older than `job_logs.retain_days` (def. 45)  |
| `vacuum`           | 1st of month 02:30 | `VACUUM (ANALYZE)` over the high-churn tables               |

`vacuum` runs plain `VACUUM`, never `VACUUM FULL`. `FULL` rewrites the table
under an `ACCESS EXCLUSIVE` lock, which on `river_job` would block the very
workers fetching from it. Plain `VACUUM` takes only `SHARE UPDATE EXCLUSIVE`, so
it makes dead space reusable without stopping traffic — the right trade for
tables that immediately refill that space. The `ANALYZE` half matters just as
much: the daily cleanups above leave the planner's row estimates badly stale.

Job transcripts outlive the jobs that wrote them on purpose. A River job row is
a status; `job_logs` is the evidence, and it is what someone comes back for
weeks after a bootstrap failed — long after the row itself has been expired.

Once you have an org - create a set of api keys for your org:
They will be in the format of:
Example values only; these are not real secrets.
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



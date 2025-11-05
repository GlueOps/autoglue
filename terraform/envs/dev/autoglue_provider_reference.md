# glueops/autoglue/autoglue – Reference (generated)

_Generated from providers schema JSON._

## Provider Configuration

| Name | Type | Flags | Description |
|---|---|---|---|
| `addr` | string | optional | Base URL to the autoglue API (e.g. https://autoglue.example.com/api/v1). Defaults to http://localhost:8080/api/v1. |
| `api_key` | string | optional, sensitive | User API key for key-only auth. |
| `bearer` | string | optional, sensitive | Bearer token (user access token). |
| `org_id` | string | optional | Organization ID (UUID). Required for user/bearer and user API key auth unless single-org membership. Omitted for org key/secret (derived server-side). |
| `org_key` | string | optional, sensitive | Org-scoped key for machine auth. |
| `org_secret` | string | optional, sensitive | Org-scoped secret for machine auth. |


### Basic usage

```hcl
terraform {
  required_providers {
    autoglue = {
      source  = "glueops/autoglue/autoglue"
      # version = ">= 0.0.0"
    }
  }
}

provider "autoglue" {
  # addr = "..."
  # api_key = "..."
  # bearer = "..."
  # org_id = "..."
  # org_key = "..."
  # org_secret = "..."
}
```
## Provider Functions

_No provider-defined functions._

## Resources

### `autoglue_annotation`

Create and manage a annotation (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `created_at` | string | computed |  |
| `id` | string | computed | ID (UUID). |
| `key` | string | required | Key. |
| `organization_id` | string | computed |  |
| `raw` | string | computed | Full server JSON from API. |
| `updated_at` | string | computed |  |
| `value` | string | required | Value. |


**Example**

```hcl
resource "autoglue_annotation" "example" {
  key = "..."
  value = "..."
}
```
### `autoglue_label`

Create and manage a label (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `created_at` | string | computed |  |
| `id` | string | computed | Server ID (UUID). |
| `key` | string | required | Key. |
| `organization_id` | string | computed |  |
| `raw` | string | computed | Full server JSON from API. |
| `updated_at` | string | computed |  |
| `value` | string | required | Value. |


**Example**

```hcl
resource "autoglue_label" "example" {
  key = "..."
  value = "..."
}
```
### `autoglue_server`

Create and manage a server (org-scoped). Mirrors API validation for role/status/ssh_key_id.

| Name | Type | Flags | Description |
|---|---|---|---|
| `created_at` | string | computed |  |
| `hostname` | string | required | Hostname. |
| `id` | string | computed | Server ID (UUID). |
| `organization_id` | string | computed |  |
| `private_ip_address` | string | required | Private IP address (required). |
| `public_ip_address` | string | optional | Public IP address (required when role = bastion). |
| `raw` | string | computed | Full server JSON from API. |
| `role` | string | required | Server role (e.g., agent/manager/bastion). Lowercased by the provider. |
| `ssh_key_id` | string | required | SSH key ID (UUID) that belongs to the org. |
| `ssh_user` | string | required | SSH username (required). |
| `status` | string | optional, computed | Status (pending|provisioning|ready|failed). Lowercased by the provider. |
| `updated_at` | string | computed |  |


**Example**

```hcl
resource "autoglue_server" "example" {
  hostname = "..."
  private_ip_address = "..."
  role = "..."
  ssh_key_id = "..."
  ssh_user = "..."
}
```
### `autoglue_ssh_key`

| Name | Type | Flags | Description |
|---|---|---|---|
| `bits` | number | optional | RSA key size (2048/3072/4096). Ignored for ed25519. |
| `comment` | string | required | Comment appended to authorized key |
| `created_at` | string | computed | Creation time (RFC3339, UTC) |
| `fingerprint` | string | computed | SHA256 fingerprint |
| `id` | string | computed | SSH key ID (UUID) |
| `name` | string | required | Display name |
| `private_key_pem` | string | computed, sensitive | Private key PEM (resource doesn’t reveal; stays empty). |
| `public_key` | string | computed | OpenSSH authorized key |
| `type` | string | optional | Key type: rsa or ed25519 (default rsa) |
| `updated_at` | string | computed | Update time (RFC3339, UTC) |


**Example**

```hcl
resource "autoglue_ssh_key" "example" {
  comment = "..."
  name = "..."
}
```
### `autoglue_taint`

Create and manage a taint (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `created_at` | string | computed |  |
| `effect` | string | required | Effect. |
| `id` | string | computed | Server ID (UUID). |
| `key` | string | required | Key. |
| `organization_id` | string | computed |  |
| `raw` | string | computed | Full server JSON from API. |
| `updated_at` | string | computed |  |
| `value` | string | required | Value. |


**Example**

```hcl
resource "autoglue_taint" "example" {
  effect = "..."
  key = "..."
  value = "..."
}
```
## Data Sources

### `autoglue_annotations`

List annotations for the organization (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `items` | (block) | computed | Annotations returned by the API. |


**Example**

```hcl
data "autoglue_annotations" "all" {}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {
#   value = try(data.autoglue_annotations.all.items[0].raw, null)
# }
```
### `autoglue_labels`

List labels for the organization (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `items` | (block) | computed | Labels returned by the API. |


**Example**

```hcl
data "autoglue_labels" "all" {}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {
#   value = try(data.autoglue_labels.all.items[0].raw, null)
# }
```
### `autoglue_servers`

List servers for the organization (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `items` | (block) | computed | Servers returned by the API. |
| `role` | string | optional | Filter by role. |
| `status` | string | optional | Filter by status (pending|provisioning|ready|failed). |


**Example**

```hcl
data "autoglue_servers" "all" {}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {
#   value = try(data.autoglue_servers.all.items[0].raw, null)
# }
```
### `autoglue_ssh_keys`

| Name | Type | Flags | Description |
|---|---|---|---|
| `fingerprint` | string | optional | Filter by exact fingerprint (client-side). |
| `keys` | (block) | computed | SSH keys |
| `name_contains` | string | optional | Filter by substring of name (client-side). |


**Example**

```hcl
data "autoglue_ssh_keys" "all" {}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {
#   value = try(data.autoglue_ssh_keys.all.items[0].raw, null)
# }
```
### `autoglue_taints`

List taints for the organization (org-scoped).

| Name | Type | Flags | Description |
|---|---|---|---|
| `items` | (block) | computed | Taints returned by the API. |


**Example**

```hcl
data "autoglue_taints" "all" {}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {
#   value = try(data.autoglue_taints.all.items[0].raw, null)
# }
```
# Autoglue Repository Architecture & Agents

## Overview

Autoglue is a Kubernetes cluster management platform built with Go that manages the lifecycle of K3s clusters across GlueOps-supported cloud providers. It provides a REST API for cluster provisioning, configuration, and management, along with a web UI and Terraform provider.

## Repository Structure

```
autoglue/
├── cmd/                    # CLI commands
├── internal/              # Internal packages
│   ├── api/              # HTTP API routes and middleware
│   ├── app/              # Application setup
│   ├── auth/             # Authentication logic
│   ├── bg/               # Background job workers
│   ├── common/           # Common utilities
│   ├── config/           # Configuration management
│   ├── db/               # Database operations
│   ├── handlers/         # HTTP request handlers
│   ├── keys/             # Cryptographic key management
│   ├── mapper/           # Data mapping utilities
│   ├── models/           # Database models
│   ├── utils/            # Utility functions
│   ├── version/          # Version information
│   └── web/              # Web UI integration
├── sdk/                   # Generated SDKs
│   └── ts/               # TypeScript SDK
├── ui/                    # Frontend application (React)
├── docs/                  # OpenAPI/Swagger documentation
├── postgres/              # PostgreSQL configuration
├── main.go                # Application entry point
├── schema.sql             # Database schema
└── docker-compose.yml     # Development environment

```

## Core Components

### 1. API Layer (`internal/api`)

The API layer provides RESTful endpoints for managing cloud resources:

- **Authentication** (`mount_auth_routes.go`): OAuth/OIDC integration, JWT tokens
- **Clusters** (`mount_cluster_routes.go`): Kubernetes cluster management
- **Servers** (`mount_server_routes.go`): Server resource management
- **SSH Keys** (`mount_ssh_routes.go`): SSH key generation and management
- **DNS** (`mount_dns_routes.go`): DNS record management
- **Load Balancers** (`mount_load_balancer_routes.go`): Load balancer configuration
- **Node Pools** (`mount_node_pool_routes.go`): Worker node pool management
- **Credentials** (`mount_credential_routes.go`): Cloud provider credentials
- **Organizations** (`mount_org_routes.go`): Multi-tenant organization management

**Middleware:**
- Request logging with zerolog
- Rate limiting (1000 requests/minute per IP)
- CORS handling
- Security headers
- Authentication/Authorization
- Request body size limits (10MB)

### 2. Background Jobs (`internal/bg`)

Autoglue uses the [Archer](https://github.com/dyaksa/archer) job queue system with PostgreSQL-backed job persistence.

**Active Job Workers:**

| Worker | Purpose | Timeout |
|--------|---------|---------|
| `bootstrap_bastion` | Provision and configure bastion host servers | Configurable (default 60s) |
| `archer_cleanup` | Clean up old job records | 5 minutes |
| `tokens_cleanup` | Purge expired refresh tokens | 5 minutes |
| `db_backup_s3` | Backup database to S3 | 15 minutes |
| `dns_reconcile` | Synchronize DNS records with Route53 | 2 minutes |
| `org_key_sweeper` | Remove expired organization API keys | 5 minutes |
| `cluster_action` | Execute cluster lifecycle actions | Configurable |

**Planned Job Workers (Currently Disabled):**

The following workers exist in the codebase but are currently commented out:
- `prepare_cluster` - Prepare infrastructure for cluster deployment
- `cluster_setup` - Initial cluster configuration
- `cluster_bootstrap` - Full Kubernetes cluster bootstrapping process

**Configuration:**
- `archer.instances`: Number of worker instances (default: 1)
- `archer.timeoutSec`: Job timeout in seconds (default: 60)
- `archer.cleanup_retain_days`: Job retention period (default: 7 days)

### 3. Data Models (`internal/models`)

**Core Models:**

- `User` - User accounts with OAuth integration
- `Organization` - Multi-tenant organizations
- `Membership` - User-organization relationships
- `ApiKey` - API authentication tokens (user and org-level)
- `OrganizationKey` - Organization-level credentials with auto-expiry
- `Cluster` - Kubernetes cluster definitions
- `NodePool` - Worker node group configurations
- `Server` - Individual server instances
- `SshKey` - SSH keypair management with encryption
- `LoadBalancer` - Load balancer configurations
- `Domain` - DNS domain management
- `Credential` - Cloud provider API credentials (AWS, etc.)
- `Job` - Background job queue records
- `SigningKey` - JWT signing keys with rotation
- `RefreshToken` - OAuth refresh token storage
- `MasterKey` - Master encryption key for data at rest
- `Label`, `Annotation`, `Taint` - Kubernetes resource metadata

### 4. Handlers (`internal/handlers`)

Request handlers implement business logic for API endpoints:

- `auth.go` - OAuth flows, token issuance
- `clusters.go` - Cluster CRUD operations
- `servers.go` - Server provisioning
- `ssh_keys.go` - SSH key generation with Ed25519/RSA support
- `dns.go` - DNS record management via Route53
- `load_balancers.go` - Load balancer configuration
- `node_pools.go` - Node pool management with labels/annotations/taints
- `credentials.go` - Cloud credential storage
- `orgs.go` - Organization management
- `me.go` - Current user information
- `me_keys.go` - User API key management
- `health.go` - Health check endpoints
- `version.go` - Version information

### 5. Security & Encryption

**Cryptography:**
- **Master Key**: AES-256-GCM encryption for root secrets
- **Organization Keys**: Per-org encryption keys derived from master key
- **SSH Keys**: Secure generation and encrypted storage
- **JWT Tokens**: RS256 signing with key rotation
- **API Keys**: Argon2id hashing for token storage
- **At-Rest Encryption**: All sensitive data (kubeconfigs, credentials, SSH keys)

**Authentication Methods:**
1. OAuth/OIDC (Google Workspace integration)
2. Bearer tokens (JWT)
3. Organization Key/Secret pairs
4. User API keys

### 6. CLI Commands (`cmd`)

- `serve` - Start the API server (default command)
- `keys generate` - Generate JWT signing keys
- `encrypt create-master` - Create master encryption key
- `db` - Database management utilities
- `version` - Display version information

### 7. Integration Points

**Cloud Providers:**
- AWS (Route53 for DNS, S3 for backups)
- Support for multi-cloud credentials

**External Services:**
- PostgreSQL (primary data store)
- S3-compatible storage (backups)
- OAuth providers (Google)

**SDKs:**
- TypeScript SDK (`sdk/ts/`) - Generated from OpenAPI spec
- Go SDK (consumed via module alias) - Used by external integrations

**External Integrations:**
- Terraform Provider - Separate repository providing IaC support for Autoglue resources

## Development Workflow

### Prerequisites
- Go 1.25.4+
- Docker & Docker Compose
- PostgreSQL (via docker-compose)
- Node.js (for UI development)

### Setup
```bash
# 1. Configure environment
cp .env.example .env

# 2. Start database
docker compose up -d

# 3. Generate JWT keys
go run . keys generate

# 4. Create master encryption key
go run . encrypt create-master

# 5. Update OpenAPI docs and SDKs
make swagger
make sdk-all

# 6. Start API server with embedded UI
go run .
```

### Build & Test
```bash
# Build application
go build -o autoglue .

# Run tests
go test ./...

# Build UI
make ui
```

**Note:** The Terraform provider is maintained in a separate repository.

## API Architecture

### Request Flow
```
Client → CORS → Rate Limit → Logger → Auth → Handler → DB/Jobs → Response
```

### Authentication Flow
1. User logs in via OAuth (Google)
2. Backend validates token with provider
3. JWT access token issued (short-lived)
4. Refresh token stored in DB
5. Organization context from `X-Org-ID` header

### Job Execution Flow
1. Handler enqueues job via `Jobs.Enqueue()`
2. Archer worker picks up job from PostgreSQL
3. Worker executes task with timeout
4. Result stored in `jobs` table
5. Retries on failure (configurable)

## Database Schema

**Key Tables:**
- `users` - User accounts
- `accounts` - OAuth provider linkage
- `organizations` - Tenant isolation
- `memberships` - User-org relationships
- `api_keys` - Authentication tokens
- `clusters` - K8s cluster definitions
- `node_pools` - Worker node groups
- `servers` - Compute instances
- `ssh_keys` - SSH keypair storage
- `load_balancers` - LB configurations
- `domains` - DNS domains
- `credentials` - Cloud API credentials
- `jobs` - Background job queue
- `signing_keys` - JWT key rotation
- `refresh_tokens` - OAuth token storage
- `master_keys` - Encryption key hierarchy

## Configuration

Environment variables (`.env`):
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_PRIVATE_ENC_KEY` - JWT private key encryption
- `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` - OAuth
- `ALLOWED_ORIGINS` - CORS configuration
- `archer.*` - Job queue settings
- AWS credentials for Route53/S3

## Deployment

### Docker
```bash
docker build -t autoglue .
docker run -p 8080:8080 --env-file .env autoglue
```

### Production Considerations
- Database connection pooling
- Rate limiting configuration
- CORS allowed origins
- JWT key rotation schedule
- Backup retention policies
- Worker instance scaling
- Monitoring and alerting

## API Documentation

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `docs/openapi.yaml`
- **SDK Documentation**: `sdk/ts/README.md`

## Testing

The repository includes:
- Unit tests for handlers (`*_test.go`)
- Test utilities (`internal/testutil/`)
- Integration tests with embedded PostgreSQL

Run tests:
```bash
go test ./internal/handlers/
go test -v ./...
```

## Key Features

1. **Multi-tenancy**: Organization-based resource isolation
2. **Encryption at Rest**: All sensitive data encrypted per-org
3. **Async Job Processing**: Background tasks with retry logic
4. **API Key Management**: Multiple authentication methods
5. **SSH Key Generation**: Automated keypair creation (RSA/Ed25519)
6. **DNS Automation**: Route53 integration for DNS records
7. **Kubernetes Management**: Cluster lifecycle automation
8. **Terraform Provider**: Infrastructure-as-Code support
9. **Web UI**: React-based management interface
10. **OpenAPI/Swagger**: Auto-generated API documentation

## Architecture Patterns

- **Repository Pattern**: Data access abstraction via GORM
- **Dependency Injection**: Dependencies passed to handlers
- **Middleware Chain**: Request processing pipeline
- **Job Queue**: Async processing with Archer
- **Multi-tenant**: Organization-scoped data isolation
- **Encryption**: Key hierarchy (master → org → resource)

## Future Enhancements

Based on commented code and structure:
- Full cluster provisioning automation
- Additional cloud provider support
- Enhanced monitoring and observability
- Cluster backup and restore
- Advanced RBAC controls
- Custom resource definitions

## Resources

- **GitHub**: https://github.com/GlueOps/autoglue
- **Production API**: https://autoglue.glueopshosted.com/api/v1
- **Pre-prod API**: https://autoglue.glueopshosted.rocks/api/v1
- **Staging API**: https://autoglue.apps.nonprod.earth.onglueops.rocks/api/v1

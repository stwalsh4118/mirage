# PBI-17: HashiCorp Vault Secret Management

[View in Backlog](../backlog.md#user-content-17)

## Overview

Implement HashiCorp Vault integration for secure, centralized secret management, including per-user Railway API token storage, secret rotation support, and migration from environment variable-based token storage.

## Problem Statement

Currently, Mirage uses a single Railway API token stored as an environment variable (`RAILWAY_API_TOKEN`). This approach has several critical limitations:

1. **No per-user tokens**: All users share the same Railway token, making it impossible to:
   - Attribute Railway actions to specific users
   - Implement per-user Railway project access controls
   - Revoke access for individual users without affecting everyone

2. **Insecure storage**: Environment variables are:
   - Visible in process listings and logs
   - Included in error reports and stack traces
   - Potentially exposed in CI/CD pipelines
   - Not encrypted at rest

3. **No secret rotation**: Changing tokens requires:
   - Application restart
   - Deployment to update environment variables
   - No graceful transition or zero-downtime rotation

4. **Limited secret types**: Only Railway tokens are managed; no support for:
   - GitHub tokens for private repository access
   - Docker registry credentials
   - Database credentials
   - Other third-party API keys

5. **No audit trail**: No way to track:
   - When secrets were accessed
   - Which user accessed which secrets
   - Secret lifecycle events (creation, rotation, deletion)

## User Stories

- As a platform engineer, I want to securely store my Railway API token so it's encrypted and not visible in logs or environment variables
- As a platform engineer, I want to use my own Railway API token so my environments are created under my Railway account
- As a platform engineer, I want to rotate my Railway token without restarting the application
- As a system admin, I want secrets encrypted at rest and in transit so they're protected from unauthorized access
- As a system admin, I want an audit trail of secret access so I can track who accessed what and when
- As a developer, I want to store additional secrets (GitHub tokens, Docker credentials) so I can access private repositories and registries

## Technical Approach

### Vault Setup & Configuration

#### Development Environment
- Use Vault in dev mode (in-memory storage) for local development
- Docker Compose service for easy local setup:
  ```yaml
  vault:
    image: hashicorp/vault:latest
    ports:
      - "8200:8200"
    environment:
      VAULT_DEV_ROOT_TOKEN_ID: dev-root-token
      VAULT_DEV_LISTEN_ADDRESS: 0.0.0.0:8200
    cap_add:
      - IPC_LOCK
  ```
- Unsealed automatically with known token for dev convenience
- No persistence (in-memory) acceptable for local testing

#### Production Environment - Self-Hosted on Railway

**Vault Server Service:**
- Deploy Vault as a Railway service using official Docker image
- Use integrated storage (Raft) for simplicity (no external dependencies)
- Mount Railway volume for persistent storage at `/vault/data`
- Configure TLS using Railway's automatic HTTPS (reverse proxy)
- Set resource limits appropriately (512MB-1GB RAM minimum)

**Auto-Unseal Configuration:**
- Option 1: Manual unseal (store unseal keys securely, requires manual intervention after restart)
- Option 2: Transit auto-unseal (use another Vault cluster as KMS - requires two Vault instances)
- Option 3: Cloud KMS auto-unseal (AWS KMS, GCP KMS) - requires cloud provider setup
- **Recommendation for MVP**: Manual unseal with documented process, migrate to auto-unseal later

**High Availability (Future):**
- Multiple Vault instances with integrated Raft storage
- Load balancer distributing requests across instances
- Automatic leader election and failover
- **Recommendation**: Start with single instance, add HA when needed

**Backup & Disaster Recovery:**
- Regular snapshots of Vault data using `vault operator raft snapshot`
- Store snapshots in Railway volume or external storage (S3, GCS)
- Automated backup cron job (daily at minimum)
- Documented recovery procedure
- Test recovery process regularly

**Monitoring & Alerting:**
- Vault health check endpoint monitoring
- Alert on seal status changes
- Alert on failed authentication attempts
- Log aggregation for audit logs
- Metrics export (optional: Prometheus integration)

#### Vault Configuration

**KV Secrets Engine v2:**
- Mount at path `/mirage`
- Namespace secrets by user: `/mirage/users/{user_id}/`
- Secret paths and types:
  - `/mirage/users/{user_id}/railway` - Railway API token (string)
  - `/mirage/users/{user_id}/github` - GitHub PAT for private repos (string)
  - `/mirage/users/{user_id}/docker/{registry}` - Docker registry credentials (JSON: username, password, email)
  - `/mirage/users/{user_id}/env_vars/{environment_id}` - Environment-specific secrets (JSON map)
  - `/mirage/users/{user_id}/custom/{key}` - Generic user secrets
- Enable versioning (keep last 10 versions)
- Set check-and-set required to prevent concurrent overwrites

**Audit Logging:**
- Enable file audit backend logging to `/vault/logs/audit.log`
- Log format: JSON for easy parsing
- Rotation: Daily rotation with 30-day retention
- Logged events: All secret operations (read, write, delete), auth attempts, policy changes
- Sensitive data masked in logs (actual secret values never logged)

**Authentication Methods:**
- Dev: Token auth with root token
- Prod: AppRole auth for Mirage backend service
  - Create `mirage-backend` role with policy to access all user secrets
  - Bind by Railway service IP or use secret-id wrapping
  - Token TTL: 1 hour with auto-renewal
  - Max TTL: 24 hours (requires re-authentication after)

**Policies:**
- `user-secrets-policy`: Template policy allowing access to `/mirage/users/{user_id}/*`
- `admin-policy`: Full access to all paths for admin operations
- `backup-policy`: Read-only access for backup operations

**Secret Metadata Tracking:**
- Custom metadata on all secrets:
  - `created_by`: User ID who created the secret
  - `created_at`: ISO timestamp
  - `last_validated`: When secret was last tested/validated
  - `secret_type`: "railway_token", "github_pat", "docker_credentials", etc.
  - `tags`: User-defined tags for organization

### Backend Vault Integration

#### Vault Client
- Install `github.com/hashicorp/vault/api` Go package
- Create `vault` package with client initialization:
  - Support multiple auth methods (token, AppRole for production)
  - Connection pooling and retry logic
  - Automatic token renewal
  - Health checking and circuit breaking

#### Secret Storage Interface
```go
type SecretMetadata struct {
    CreatedBy      string
    CreatedAt      time.Time
    LastValidated  *time.Time
    SecretType     string
    Tags           []string
    Version        int
}

type Secret struct {
    Key      string
    Value    string
    Metadata SecretMetadata
}

type DockerCredentials struct {
    Registry string
    Username string
    Password string
    Email    string
}

type SecretStore interface {
    // Railway token management
    StoreRailwayToken(ctx context.Context, userID, token string) error
    GetRailwayToken(ctx context.Context, userID string) (string, error)
    DeleteRailwayToken(ctx context.Context, userID string) error
    RotateRailwayToken(ctx context.Context, userID, newToken string) error
    ValidateRailwayToken(ctx context.Context, userID string) error
    
    // GitHub PAT management
    StoreGitHubToken(ctx context.Context, userID, token string) error
    GetGitHubToken(ctx context.Context, userID string) (string, error)
    DeleteGitHubToken(ctx context.Context, userID string) error
    
    // Docker credentials management
    StoreDockerCredentials(ctx context.Context, userID string, creds DockerCredentials) error
    GetDockerCredentials(ctx context.Context, userID, registry string) (DockerCredentials, error)
    ListDockerRegistries(ctx context.Context, userID string) ([]string, error)
    DeleteDockerCredentials(ctx context.Context, userID, registry string) error
    
    // Environment-specific secrets (key-value pairs per environment)
    StoreEnvironmentSecret(ctx context.Context, userID, envID, key, value string) error
    GetEnvironmentSecret(ctx context.Context, userID, envID, key string) (string, error)
    GetAllEnvironmentSecrets(ctx context.Context, userID, envID string) (map[string]string, error)
    DeleteEnvironmentSecret(ctx context.Context, userID, envID, key string) error
    
    // Generic secret management
    StoreSecret(ctx context.Context, userID, key, value string, metadata SecretMetadata) error
    GetSecret(ctx context.Context, userID, key string) (Secret, error)
    GetSecretValue(ctx context.Context, userID, key string) (string, error)
    DeleteSecret(ctx context.Context, userID, key string) error
    ListSecrets(ctx context.Context, userID string) ([]SecretMetadata, error)
    
    // Version management
    GetSecretVersion(ctx context.Context, userID, key string, version int) (Secret, error)
    ListSecretVersions(ctx context.Context, userID, key string) ([]int, error)
    RollbackSecret(ctx context.Context, userID, key string, toVersion int) error
    
    // Metadata management
    UpdateSecretMetadata(ctx context.Context, userID, key string, metadata SecretMetadata) error
    GetSecretMetadata(ctx context.Context, userID, key string) (SecretMetadata, error)
    
    // Health and status
    HealthCheck(ctx context.Context) error
    GetVaultStatus(ctx context.Context) (map[string]interface{}, error)
}
```

#### Vault Implementation
- Implement `SecretStore` interface using Vault KV v2
- Use versioned secrets for rotation support
- Implement caching layer for frequently accessed secrets (with TTL)
- Add metrics and logging for secret operations

#### Migration Strategy
- Keep environment variable fallback for backward compatibility
- Implement feature flag to enable/disable Vault
- Gradual migration path:
  1. Deploy with Vault optional, env var still works
  2. UI prompt for users to migrate their tokens to Vault
  3. Eventually deprecate env var support

### Railway Client Updates

Currently, Railway client is initialized once with a single token:
```go
rw := railway.NewClient(endpoint, cfg.RailwayAPIToken, httpc)
```

New approach - user-specific clients:
```go
// Create Railway client for specific user
func GetUserRailwayClient(ctx context.Context, userID string, vault SecretStore) (*railway.Client, error) {
    token, err := vault.GetRailwayToken(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get railway token: %w", err)
    }
    return railway.NewClient(endpoint, token, httpc), nil
}
```

Update all Railway API calls to use per-user clients:
- Extract user ID from request context (set by auth middleware)
- Fetch Railway token from Vault
- Create Railway client with user's token
- Execute Railway API call

### API Endpoints for Secret Management

#### Railway Token Management
- `POST /api/v1/secrets/railway` - Store/update Railway token
  - Request: `{"token": "railway-token-here"}`
  - Validates token by making test Railway API call
  - Stores in Vault under user's namespace
  - Response: `{"success": true, "validated": true, "stored_at": "timestamp"}`
  
- `GET /api/v1/secrets/railway/status` - Check if user has Railway token configured
  - Returns: `{"configured": true, "last_validated": "timestamp", "needs_rotation": false}`
  
- `POST /api/v1/secrets/railway/validate` - Test Railway token validity
  - Makes test API call to Railway
  - Updates last_validated timestamp
  - Response: `{"valid": true, "permissions": ["read", "write"]}`
  
- `DELETE /api/v1/secrets/railway` - Delete Railway token
  - Removes from Vault
  - Response: `{"success": true, "affected_environments": 5}`
  
- `POST /api/v1/secrets/railway/rotate` - Rotate Railway token
  - Request: `{"new_token": "new-railway-token"}`
  - Validates new token
  - Stores as new version in Vault
  - Keeps old version for rollback
  - Response: `{"success": true, "version": 2, "previous_version": 1}`

#### GitHub Token Management
- `POST /api/v1/secrets/github` - Store GitHub PAT
  - Request: `{"token": "ghp_..."}`
  - Optionally validates by calling GitHub API
  - Response: `{"success": true, "username": "detected-from-token"}`
  
- `GET /api/v1/secrets/github/status` - Check GitHub token status
  - Returns: `{"configured": true, "last_validated": "timestamp", "scopes": ["repo", "read:org"]}`
  
- `DELETE /api/v1/secrets/github` - Delete GitHub token

#### Docker Credentials Management
- `POST /api/v1/secrets/docker` - Store Docker registry credentials
  - Request: `{"registry": "docker.io", "username": "user", "password": "pass", "email": "user@example.com"}`
  - Stores credentials for specific registry
  - Response: `{"success": true, "registry": "docker.io"}`
  
- `GET /api/v1/secrets/docker` - List configured Docker registries
  - Returns: `{"registries": ["docker.io", "ghcr.io", "registry.gitlab.com"]}`
  
- `GET /api/v1/secrets/docker/{registry}` - Get credentials for specific registry
  - Returns: `{"registry": "docker.io", "username": "user", "email": "user@example.com"}`
  - Password not returned, only metadata
  
- `DELETE /api/v1/secrets/docker/{registry}` - Delete Docker credentials

#### Environment-Specific Secrets
- `POST /api/v1/environments/{env_id}/secrets` - Store environment secret
  - Request: `{"key": "DATABASE_URL", "value": "postgres://..."}`
  - Stores in Vault under environment-specific path
  - Response: `{"success": true, "key": "DATABASE_URL"}`
  
- `GET /api/v1/environments/{env_id}/secrets` - List environment secret keys
  - Returns: `{"secrets": ["DATABASE_URL", "API_KEY", "WEBHOOK_SECRET"]}`
  
- `GET /api/v1/environments/{env_id}/secrets/{key}` - Get specific environment secret
  - Returns: `{"key": "DATABASE_URL", "value": "postgres://...", "created_at": "..."}`
  
- `DELETE /api/v1/environments/{env_id}/secrets/{key}` - Delete environment secret
  
- `POST /api/v1/environments/{env_id}/secrets/bulk` - Store multiple secrets at once
  - Request: `{"secrets": {"KEY1": "value1", "KEY2": "value2"}}`
  - Atomic operation (all or nothing)
  - Response: `{"success": true, "stored": 2}`

#### Generic Secret Management
- `POST /api/v1/secrets` - Store generic secret
  - Request: `{"key": "custom_key", "value": "secret-value", "tags": ["prod", "important"]}`
  - Stores with metadata
  - Response: `{"success": true, "key": "custom_key", "version": 1}`
  
- `GET /api/v1/secrets` - List user's secret keys with metadata
  - Query params: `?type=railway_token&tag=prod`
  - Returns: `{"secrets": [{"key": "railway", "type": "railway_token", "created_at": "...", "version": 2}]}`
  
- `GET /api/v1/secrets/{key}` - Get secret with full metadata
  - Returns: `{"key": "custom_key", "value": "secret-value", "metadata": {...}}`
  
- `DELETE /api/v1/secrets/{key}` - Delete secret
  - Soft delete (versions preserved)
  - Response: `{"success": true, "versions_archived": 3}`

#### Version Management
- `GET /api/v1/secrets/{key}/versions` - List all versions of a secret
  - Returns: `{"versions": [{"version": 2, "created_at": "..."}, {"version": 1, "created_at": "..."}]}`
  
- `GET /api/v1/secrets/{key}/versions/{version}` - Get specific version
  - Returns: `{"key": "railway", "value": "old-token", "version": 1, "created_at": "..."}`
  
- `POST /api/v1/secrets/{key}/rollback` - Rollback to previous version
  - Request: `{"to_version": 1}`
  - Creates new version with old value
  - Response: `{"success": true, "new_version": 3, "restored_from": 1}`

#### Secret Health & Status
- `GET /api/v1/secrets/health` - Check Vault connectivity and health
  - Returns: `{"vault_reachable": true, "sealed": false, "version": "1.15.0"}`
  
- `GET /api/v1/secrets/audit` - Get user's secret access audit log
  - Query params: `?secret_key=railway&action=read&start_date=2024-01-01`
  - Returns: `{"logs": [{"action": "read", "secret": "railway", "timestamp": "...", "ip": "..."}]}`

### Security Considerations

1. **Secret Transmission**:
   - All API calls over HTTPS only
   - Secrets never logged in plain text
   - Secrets masked in error messages

2. **Secret Storage**:
   - Vault encrypts all secrets at rest
   - Encryption keys managed by Vault's barrier
   - Automatic key rotation supported

3. **Access Control**:
   - Users can only access their own secrets
   - Vault policies enforce namespace isolation
   - Audit log tracks all access attempts

4. **Secret Lifecycle**:
   - TTLs can be configured per secret
   - Automatic cleanup of expired secrets
   - Version history maintained for auditing

## UX/UI Considerations

### First-Time Setup Flow
1. User signs up and logs in via Clerk
2. Dashboard shows "Configure Railway Token" banner
3. User clicks "Add Railway Token" button
4. Modal appears with:
   - Instructions on getting Railway token
   - Link to Railway account settings
   - Input field for token (masked)
   - "Validate & Save" button
5. Backend validates token by testing Railway API
6. Success message: "Railway token configured successfully"
7. User can now create environments

### Credentials Management UI (Settings Page)

**Main Credentials Page:**
- Tabbed interface:
  - Railway Token tab
  - GitHub Token tab  
  - Docker Registries tab
  - Custom Secrets tab
  
**Railway Token Tab:**
- Status card showing:
  - Connection status (connected, not configured, invalid)
  - Last validated timestamp
  - Token expiry warning (if detectable)
- Actions:
  - "Configure Token" / "Update Token" button
  - "Test Connection" button (calls Railway API)
  - "Rotate Token" button (with confirmation)
  - "Remove Token" button (shows affected resources count)
- Modal for adding/updating token:
  - Masked input field
  - "How to get Railway token" help link
  - Live validation as you type
  - "Save & Validate" button

**GitHub Token Tab:**
- Similar to Railway, showing:
  - Username detected from token
  - Scopes/permissions granted
  - Repos accessible count
- Test connection button validates with GitHub API

**Docker Registries Tab:**
- Table of configured registries:
  - Registry URL
  - Username
  - Last used timestamp
  - Actions (test, edit, remove)
- "Add Registry" button opens modal:
  - Registry URL (dropdown with common ones + custom)
  - Username
  - Password (masked)
  - Email
  - "Test Login" button

**Custom Secrets Tab:**
- Table of user-defined secrets:
  - Key name
  - Type/category tags
  - Created date
  - Last accessed
  - Version count
  - Actions (view, edit, delete, rollback)
- "Add Secret" button
- Filter by tags
- Search by key name

**Version History Modal:**
- Timeline view of secret versions
- Each version showing:
  - Version number
  - Created timestamp
  - Created by (if metadata available)
  - "View" and "Rollback to this version" buttons
- Current version highlighted

**Environment Secrets Management:**
- In environment detail page, new "Secrets" tab
- Table of environment-specific secrets:
  - Key
  - Description (user-provided)
  - Last updated
  - Actions (edit, delete)
- "Add Secret" button
- "Bulk Import" button (paste key=value pairs)
- "Export" button (masked values for security)
- Warning about Railway syncing these secrets

### Error Handling
- Clear error messages when token is missing/invalid
- Guidance on how to configure token
- Link to documentation
- Graceful degradation (show UI but disable creation actions)

## Acceptance Criteria

1. **Vault Infrastructure**:
   - Vault server deployed on Railway with persistent storage
   - Vault health checks and monitoring configured
   - Audit logging enabled and accessible
   - Automated backups configured and tested
   - Recovery procedure documented and tested
   - Dev environment uses Docker Compose Vault in dev mode

2. **Vault Client Integration**:
   - Go Vault client successfully connects and authenticates
   - AppRole authentication works for production
   - Automatic token renewal prevents expiration
   - Connection failures handled gracefully with retries and circuit breaking
   - Caching layer reduces Vault calls without sacrificing security
   - Health check endpoint reports Vault status

3. **Railway Token Management**:
   - Users can store/update/delete their Railway tokens
   - Token validation works by testing Railway API before storage
   - Railway client dynamically uses per-user tokens
   - Token rotation creates new versions while preserving history
   - Token status endpoint shows configuration state and last validation
   - Environments show clear error when user's Railway token missing/invalid

4. **GitHub Token Management**:
   - Users can store/delete GitHub PATs
   - Token validation detects username and scopes
   - GitHub scanner uses user's token for private repo access
   - Token status shows scopes and permissions

5. **Docker Credentials Management**:
   - Users can store credentials for multiple registries
   - Credentials properly formatted for Railway Docker authentication
   - List endpoint shows all configured registries
   - Test functionality validates credentials work

6. **Environment-Specific Secrets**:
   - Users can store key-value pairs per environment
   - Bulk operations work for importing multiple secrets
   - Environment secrets are properly namespaced by user and environment
   - Secrets can be synced to Railway environment variables (future integration)

7. **Generic Secret Management**:
   - Users can store arbitrary secrets with custom keys
   - Metadata (tags, types, timestamps) properly tracked
   - List/filter/search functionality works correctly
   - Secrets properly isolated between users

8. **Version Management**:
   - Secret versions properly tracked (keep last 10)
   - Users can view version history
   - Users can retrieve specific versions
   - Rollback creates new version with old value
   - Version metadata includes timestamps

9. **Secret Management UI**:
   - Credentials settings page has tabs for all secret types
   - Railway token tab shows status, allows add/update/test/remove
   - GitHub token tab shows username and scopes
   - Docker registries tab lists all registries with actions
   - Custom secrets tab shows searchable/filterable table
   - Environment detail page has secrets tab
   - Version history modal shows timeline
   - All secret inputs are masked
   - Clear success/error feedback for all operations

10. **Security & Audit**:
    - All secrets encrypted at rest in Vault
    - Secrets never appear in logs or error messages
    - Secret values masked in UI and API responses
    - Vault audit log tracks all secret access
    - Users can only access their own secrets (enforced by policies)
    - Audit endpoint allows users to view their secret access history
    - Admin users cannot directly access other users' secrets

11. **Error Handling & Resilience**:
    - Vault unavailability doesn't crash the application
    - Cached tokens used when Vault temporarily unreachable
    - Circuit breaker prevents cascade failures
    - Clear error messages guide users to fix issues
    - Railway operations fail gracefully when token missing
    - Retry logic with exponential backoff for transient failures

12. **Migration & Backward Compatibility**:
    - Feature flag controls Vault usage
    - Environment variable fallback works when Vault disabled
    - Existing deployments continue working during migration
    - Migration guide documents upgrade path
    - Database tracks Vault vs env var token source

## Dependencies

- HashiCorp Vault server (Docker for dev, managed service for prod)
- Environment variables for Vault configuration:
  - `VAULT_ADDR` - Vault server address
  - `VAULT_TOKEN` - Vault authentication token (or AppRole credentials)
  - `VAULT_NAMESPACE` - Vault namespace (for Vault Enterprise)
  - `VAULT_SKIP_VERIFY` - Skip TLS verification (dev only)
- `github.com/hashicorp/vault/api` Go package
- PBI 16 (User Management) must be completed first to have user IDs

## Open Questions

1. **Vault Hosting**: ✅ ANSWERED - Self-host on Railway with integrated Raft storage

2. **Auto-Unseal Strategy**: How should we handle Vault unsealing in production?
   - **Option A**: Manual unseal (requires operator intervention after restart)
   - **Option B**: Cloud KMS auto-unseal (AWS/GCP KMS integration)
   - **Option C**: Transit auto-unseal (second Vault cluster as KMS)
   - **Recommendation**: Start with manual unseal + documented process, add cloud KMS later
   - **Question**: Which cloud provider KMS to prioritize if we add auto-unseal?

3. **Secret Rotation**: Should we implement automatic rotation?
   - **Proposal**: Manual rotation for MVP, automatic rotation in future PBI
   - **Consideration**: Railway doesn't expose token TTL, so we can't automate rotation easily

4. **Secret Versioning**: How many versions should we keep?
   - **Decision**: Keep last 10 versions per secret
   - **Rationale**: Balance between history and storage

5. **Caching Strategy**: How long should we cache secrets?
   - **Decision**: 5 minute TTL with explicit invalidation on update/delete
   - **Rationale**: Reduces Vault load while keeping secrets reasonably fresh
   - **Question**: Should cache be in-memory or Redis for multi-instance deployments?

6. **Fallback Strategy**: What if Vault is completely unavailable?
   - **Proposal**: 
     - Use cached tokens with warning banner
     - Prevent new operations requiring fresh tokens
     - Allow read-only operations
     - Alert admins of Vault outage
   - **Question**: Should we implement degraded mode with database-stored encrypted secrets as last resort?

7. **Admin Access**: Should admins be able to view/manage all users' secrets?
   - **Decision**: No direct access to secret values, but admins can:
     - Revoke/delete any user's secrets (forcing re-configuration)
     - View secret metadata (when created, last accessed, etc.)
     - View audit logs of secret access
   - **Rationale**: Maintains security while allowing admin intervention for issues

8. **Backup Frequency**: How often should we snapshot Vault data?
   - **Proposal**: Daily snapshots kept for 30 days
   - **Question**: Should we also do hourly snapshots kept for 7 days?

9. **Secret Limits**: Should we limit number of secrets per user?
   - **Proposal**: No hard limits for MVP, add soft limits (with warnings) if needed
   - **Consideration**: Vault KV v2 has no practical limits for our scale

10. **Environment Variable Sync**: Should environment-specific secrets automatically sync to Railway?
    - **Proposal**: Defer to future PBI, requires Railway API integration
    - **MVP**: Manual copy-paste or CLI tool

11. **Secret Templates**: Should we support secret templates/presets?
    - **Example**: "Node.js app" template with DATABASE_URL, REDIS_URL, etc.
    - **Proposal**: Defer to future PBI

12. **Secret Sharing**: Should users be able to share secrets with team members?
    - **Proposal**: Defer until we implement teams/organizations (future PBI)
    - **Security**: Requires careful access control design

## Related Tasks

This PBI will be broken down into tasks covering:

**Infrastructure & Dependencies:**
- Research and document HashiCorp Vault API (`vault-api-guide.md`)
- Setup Docker Compose for local Vault development
- Deploy Vault server on Railway with Raft storage
- Configure Railway volume for Vault data persistence
- Create Vault initialization and unsealing documentation
- Setup AppRole authentication for production
- Configure Vault policies (user-secrets, admin, backup)
- Enable and configure audit logging
- Setup automated backup cron job
- Document disaster recovery procedure

**Backend Vault Client:**
- Install and configure `github.com/hashicorp/vault/api` package
- Implement Vault client initialization with connection pooling
- Implement automatic token renewal
- Implement health checking and circuit breaker
- Implement secret caching layer with TTL
- Create Vault client mock for testing
- Add Vault connectivity health check endpoint

**Secret Storage Interface & Implementation:**
- Define SecretStore interface with all operations
- Implement Railway token management (store, get, delete, rotate, validate)
- Implement GitHub token management (store, get, delete, validate)
- Implement Docker credentials management (store, get, list, delete)
- Implement environment-specific secrets (store, get, get all, delete, bulk)
- Implement generic secret management (store, get, delete, list with filters)
- Implement version management (get version, list versions, rollback)
- Implement metadata management (update metadata, get metadata)
- Implement secret validation functions (test Railway, GitHub, Docker)
- Add metrics and logging for all secret operations

**Railway Client Integration:**
- Refactor Railway client to support per-user tokens
- Update all Railway operations to fetch user token from Vault
- Implement fallback to cached token on Vault unavailability
- Add clear error messages when user token missing/invalid
- Update Environment controller to use user-specific Railway client
- Update Service controller to use user-specific Railway client
- Update all Railway API calls with proper error handling

**API Endpoints - Railway Tokens:**
- Implement POST /api/v1/secrets/railway (store/update)
- Implement GET /api/v1/secrets/railway/status
- Implement POST /api/v1/secrets/railway/validate
- Implement DELETE /api/v1/secrets/railway
- Implement POST /api/v1/secrets/railway/rotate

**API Endpoints - GitHub Tokens:**
- Implement POST /api/v1/secrets/github (store)
- Implement GET /api/v1/secrets/github/status
- Implement DELETE /api/v1/secrets/github

**API Endpoints - Docker Credentials:**
- Implement POST /api/v1/secrets/docker (store)
- Implement GET /api/v1/secrets/docker (list registries)
- Implement GET /api/v1/secrets/docker/:registry
- Implement DELETE /api/v1/secrets/docker/:registry

**API Endpoints - Environment Secrets:**
- Implement POST /api/v1/environments/:id/secrets
- Implement GET /api/v1/environments/:id/secrets
- Implement GET /api/v1/environments/:id/secrets/:key
- Implement DELETE /api/v1/environments/:id/secrets/:key
- Implement POST /api/v1/environments/:id/secrets/bulk

**API Endpoints - Generic Secrets:**
- Implement POST /api/v1/secrets (store)
- Implement GET /api/v1/secrets (list with filters)
- Implement GET /api/v1/secrets/:key
- Implement DELETE /api/v1/secrets/:key
- Implement GET /api/v1/secrets/:key/versions
- Implement GET /api/v1/secrets/:key/versions/:version
- Implement POST /api/v1/secrets/:key/rollback
- Implement GET /api/v1/secrets/health
- Implement GET /api/v1/secrets/audit

**Database Schema:**
- Add VaultEnabled boolean to User model (track migration status)
- Add LastValidated timestamp fields for tracking token health
- Create migration to add new fields
- Add indexes for query performance

**Frontend - Credentials Settings Page:**
- Create main credentials/settings page layout
- Create tabbed interface (Railway, GitHub, Docker, Custom)
- Implement Railway Token tab with status card
- Implement Railway token add/update modal
- Implement Railway token test connection button
- Implement Railway token rotation with confirmation
- Implement Railway token removal with affected resources warning

**Frontend - GitHub & Docker Tabs:**
- Implement GitHub Token tab with status and validation
- Implement GitHub token add/remove functionality
- Implement Docker Registries tab with table
- Implement Docker registry add modal
- Implement Docker credential test functionality
- Implement Docker credential removal

**Frontend - Custom Secrets:**
- Implement Custom Secrets tab with searchable table
- Implement add custom secret modal
- Implement secret filtering by tags/type
- Implement secret editing
- Implement secret deletion with confirmation
- Implement version history modal
- Implement rollback functionality

**Frontend - Environment Secrets:**
- Add Secrets tab to environment detail page
- Implement environment secrets table
- Implement add environment secret modal
- Implement bulk import functionality
- Implement export (with masking) functionality
- Implement edit/delete environment secrets

**Migration & Backward Compatibility:**
- Implement feature flag for Vault enablement
- Implement fallback to env var RAILWAY_API_TOKEN
- Create migration guide documentation
- Implement UI prompts for users to migrate tokens
- Add admin dashboard showing migration status
- Create CLI tool for bulk migration (optional)

**Testing:**
- Unit tests for Vault client operations
- Unit tests for secret storage interface
- Unit tests for caching layer
- Integration tests with real Vault (Docker)
- Integration tests for Railway token management
- Integration tests for all secret types
- Integration tests for version management
- E2E tests for secret configuration flows
- Load testing for secret caching effectiveness
- Chaos testing for Vault unavailability scenarios

**Security & Operations:**
- Security audit of secret handling
- Penetration testing of Vault integration
- Document Vault backup and recovery procedures
- Document Vault unsealing procedures
- Create runbook for Vault operational issues
- Setup monitoring and alerting for Vault
- Create incident response guide for secret breaches
- Document secret rotation best practices

**Documentation:**
- User guide for configuring Railway tokens
- User guide for managing secrets
- Admin guide for Vault operations
- Developer guide for Vault integration
- API documentation for all secret endpoints
- Troubleshooting guide for common issues
- Security best practices documentation


# Environment Variables

This document lists all environment variables used by the Mirage application.

## Application Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_ENV` | No | `development` | Application environment (development, production) |
| `HTTP_PORT` | No | `8080` | HTTP server port |

## Database Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | No | - | Primary database connection URL |
| `DB_URL` | No | - | Alternative database URL (fallback) |

## Railway API Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `RAILWAY_API_TOKEN` | Yes* | - | Railway API token (required if VAULT_ENABLED=false) |
| `RAILWAY_PROJECT_ID` | No | - | Railway project ID |
| `RAILWAY_GRAPHQL_ENDPOINT` | No | - | Railway GraphQL endpoint override |

*When `VAULT_ENABLED=true`, user-specific Railway tokens are fetched from Vault instead.

## CORS Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ALLOWED_ORIGINS` | No | `http://localhost:3000,http://127.0.0.1:3000,http://localhost:3002` | Comma-separated list of allowed CORS origins |

## Status Poller Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `POLL_INTERVAL_SECONDS` | No | `30` | Status polling interval in seconds |
| `POLL_JITTER_FRACTION` | No | `0.2` | Jitter fraction for polling (0-1) |

## Clerk Authentication

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `CLERK_SECRET_KEY` | Yes | - | Clerk secret key for JWT verification |
| `CLERK_WEBHOOK_SECRET` | Yes | - | Clerk webhook signing secret |

## HashiCorp Vault Configuration

### Vault Enablement

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VAULT_ENABLED` | No | `false` | Enable Vault integration for secret management |

### Vault Connection

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VAULT_ADDR` | Yes* | - | Vault server address (e.g., `http://localhost:8200` or `https://vault.railway.app`) |
| `VAULT_NAMESPACE` | No | - | Vault namespace (Vault Enterprise only) |
| `VAULT_SKIP_VERIFY` | No | `false` | Skip TLS verification (**DEVELOPMENT ONLY**) |
| `VAULT_MOUNT_PATH` | No | `mirage` | Mount path for Mirage secrets in Vault |

*Required if `VAULT_ENABLED=true`

### Vault Authentication

For **development**, use token authentication:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VAULT_TOKEN` | Yes* | - | Vault authentication token |

For **production**, use AppRole authentication:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `VAULT_ROLE_ID` | Yes* | - | Vault AppRole role ID |
| `VAULT_SECRET_ID` | Yes* | - | Vault AppRole secret ID |

*Either `VAULT_TOKEN` or both `VAULT_ROLE_ID` and `VAULT_SECRET_ID` required when `VAULT_ENABLED=true`

## Example Configuration Files

### Development (.env file)

```bash
# Application
APP_ENV=development
HTTP_PORT=8080

# Database
DATABASE_URL=

# Clerk Authentication
CLERK_SECRET_KEY=your_clerk_secret_key
CLERK_WEBHOOK_SECRET=your_webhook_secret

# Vault (using docker-compose)
VAULT_ENABLED=true
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=dev-root-token
VAULT_SKIP_VERIFY=true
VAULT_MOUNT_PATH=mirage

# Railway (fallback when Vault disabled)
RAILWAY_API_TOKEN=your_railway_token
```

### Production (Railway environment variables)

```bash
# Application
APP_ENV=production
HTTP_PORT=8080

# Database
DATABASE_URL=postgresql://...

# Clerk Authentication
CLERK_SECRET_KEY=your_clerk_secret_key
CLERK_WEBHOOK_SECRET=your_webhook_secret

# Vault (production AppRole)
VAULT_ENABLED=true
VAULT_ADDR=https://vault-service.railway.app
VAULT_ROLE_ID=your_role_id
VAULT_SECRET_ID=your_secret_id
VAULT_SKIP_VERIFY=false
VAULT_MOUNT_PATH=mirage
```

## Vault Configuration Notes

1. **Development Mode:**
   - Use `VAULT_TOKEN=dev-root-token` with docker-compose
   - Set `VAULT_SKIP_VERIFY=true` (no TLS in dev mode)
   - Use `VAULT_ADDR=http://localhost:8200`

2. **Production Mode:**
   - Use AppRole authentication (`VAULT_ROLE_ID` + `VAULT_SECRET_ID`)
   - Set `VAULT_SKIP_VERIFY=false` (Railway handles TLS)
   - Use `VAULT_ADDR=https://your-vault-service.railway.app`

3. **Backward Compatibility:**
   - When `VAULT_ENABLED=false`, the system uses `RAILWAY_API_TOKEN` from environment
   - This allows gradual migration to Vault-based secret management

4. **Security:**
   - Never commit `.env` files to version control
   - Rotate `VAULT_SECRET_ID` regularly in production
   - Use Vault's token renewal to avoid token expiration



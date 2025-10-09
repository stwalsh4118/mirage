# HashiCorp Vault Setup Guide

This guide explains how to set up HashiCorp Vault for local development and production use with Mirage.

## Local Development Setup

### Prerequisites

- Docker and Docker Compose installed
- Mirage API repository cloned

### Step 1: Start Vault with Docker Compose

From the project root, run:

```bash
docker-compose up -d vault
```

This starts Vault in development mode with:
- **Address:** http://localhost:8200
- **Root Token:** `dev-root-token`
- **Auto-unsealed:** Yes (dev mode automatically unseals)
- **Storage:** In-memory (data is lost on restart)

### Step 2: Verify Vault is Running

Check Vault status:

```bash
docker-compose ps vault
```

You should see the vault service running.

Test the API:

```bash
curl http://localhost:8200/v1/sys/health
```

Expected response:
```json
{
  "initialized": true,
  "sealed": false,
  "standby": false,
  "version": "1.15.x"
}
```

### Step 3: Configure Environment Variables

Create a `.env` file in the `api/` directory (if it doesn't exist):

```bash
# Vault Configuration
VAULT_ENABLED=true
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=dev-root-token
VAULT_SKIP_VERIFY=true
VAULT_MOUNT_PATH=mirage
```

See [ENVIRONMENT_VARIABLES.md](./ENVIRONMENT_VARIABLES.md) for all available options.

### Step 4: Access Vault UI (Optional)

The Vault UI is accessible at: http://localhost:8200/ui

**Login Method:** Token  
**Token:** `dev-root-token`

From the UI, you can:
- Browse secrets
- View audit logs
- Manage policies
- Test API calls using the built-in OpenAPI explorer (type `api` in the console)

### Step 5: Initialize KV v2 Secrets Engine (First Time Only)

The Mirage application expects a KV v2 secrets engine mounted at `/mirage`. Initialize it:

```bash
# Using Docker exec
docker exec mirage-vault-dev vault secrets enable -path=mirage kv-v2

# Configure versioning (keep last 10 versions)
docker exec mirage-vault-dev vault kv metadata put -max-versions=10 mirage/config
```

Or using curl:

```bash
# Enable KV v2 at /mirage
curl -X POST \
  -H "X-Vault-Token: dev-root-token" \
  -d '{"type":"kv-v2"}' \
  http://localhost:8200/v1/sys/mounts/mirage
```

### Step 6: Start the Mirage API

With Vault running, start the Mirage API:

```bash
cd api
go run cmd/server/main.go
```

You should see a log message: `Vault client initialized successfully`

---

## Production Setup (Railway)

### Overview

For production, Vault runs as a separate Railway service with:
- **Storage:** Raft integrated storage (persisted to Railway volume)
- **TLS:** Handled by Railway's reverse proxy
- **Authentication:** AppRole (not root token)
- **High Availability:** Single instance for MVP (can scale later)

### Step 1: Deploy Vault Service on Railway

1. **Create a new Railway service**
   - Name: `mirage-vault`
   - Docker Image: `hashicorp/vault:latest`

2. **Add a Railway Volume**
   - Mount Path: `/vault/data`
   - Size: 10GB (minimum 1GB)

3. **Configure Environment Variables**
   ```
   VAULT_ADDR=http://127.0.0.1:8200
   ```

4. **Create Vault Configuration**
   
   You'll need to create a `vault-config.hcl` file (can be added via Railway's file editor or build process):
   
   ```hcl
   storage "raft" {
     path = "/vault/data"
   }
   
   listener "tcp" {
     address = "0.0.0.0:8200"
     tls_disable = true  # Railway handles TLS via reverse proxy
   }
   
   api_addr = "https://mirage-vault.railway.app"
   cluster_addr = "https://mirage-vault.railway.app:8201"
   ui = true
   ```

5. **Expose Port 8200**

### Step 2: Initialize Vault

After deploying, you need to initialize Vault **once**:

```bash
# Get Railway shell access
railway shell

# Initialize Vault
vault operator init

# Output will look like:
# Unseal Key 1: <key1>
# Unseal Key 2: <key2>
# Unseal Key 3: <key3>
# Unseal Key 4: <key4>
# Unseal Key 5: <key5>
#
# Initial Root Token: <root-token>
```

**IMPORTANT:** 
- Save the unseal keys and root token securely (password manager, secrets vault)
- You need 3 of 5 keys to unseal Vault after restarts
- **DO NOT** lose these keys - you cannot recover Vault data without them

### Step 3: Unseal Vault

Vault starts sealed after initialization or restart. Unseal it:

```bash
# You need to provide 3 different keys
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>
```

After 3 keys, Vault status should show `Sealed: false`

### Step 4: Enable KV v2 Secrets Engine

```bash
export VAULT_TOKEN=<root-token>

# Enable KV v2 at /mirage
vault secrets enable -path=mirage kv-v2

# Configure versioning
vault kv metadata put -max-versions=10 mirage/config
```

### Step 5: Configure AppRole Authentication

```bash
# Enable AppRole auth method
vault auth enable approle

# Create policy for Mirage backend
vault policy write mirage-backend-policy - <<EOF
path "mirage/data/users/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "mirage/metadata/users/*" {
  capabilities = ["read", "list"]
}
EOF

# Create AppRole
vault write auth/approle/role/mirage-backend \
  token_policies="mirage-backend-policy" \
  token_ttl=1h \
  token_max_ttl=24h \
  bind_secret_id=true

# Get Role ID (store in Railway env var: VAULT_ROLE_ID)
vault read auth/approle/role/mirage-backend/role-id

# Generate Secret ID (store in Railway env var: VAULT_SECRET_ID)
vault write -f auth/approle/role/mirage-backend/secret-id
```

### Step 6: Enable Audit Logging

```bash
vault audit enable file file_path=/vault/logs/audit.log
```

### Step 7: Configure Mirage API Service

In your Mirage API Railway service, add environment variables:

```
VAULT_ENABLED=true
VAULT_ADDR=https://mirage-vault.railway.app
VAULT_ROLE_ID=<role-id-from-step-5>
VAULT_SECRET_ID=<secret-id-from-step-5>
VAULT_SKIP_VERIFY=false
VAULT_MOUNT_PATH=mirage
```

### Step 8: Test the Connection

Deploy the Mirage API and check logs for:
```
authenticated with Vault using AppRole
Vault client initialized successfully
```

---

## Common Operations

### Check Vault Status

```bash
# Development
docker exec mirage-vault-dev vault status

# Production
railway run vault status
```

### Unseal Vault (Production only)

After Railway service restarts, Vault will be sealed:

```bash
railway shell
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>
```

### View Secrets

```bash
# Development
docker exec mirage-vault-dev vault kv get mirage/users/user-123/railway

# Production
railway run vault kv get mirage/users/user-123/railway
```

### List Secrets

```bash
docker exec mirage-vault-dev vault kv list mirage/users/user-123
```

### Backup Vault Data (Production)

```bash
# Take a snapshot
railway run vault operator raft snapshot save backup.snap

# Download from Railway
railway files download backup.snap
```

### Restore from Backup

```bash
railway run vault operator raft snapshot restore backup.snap
```

---

## Troubleshooting

### Vault is Sealed

**Symptom:** HTTP 503 responses, logs show "Vault is sealed"

**Solution:** Unseal Vault using 3 of 5 unseal keys:
```bash
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>
```

### Permission Denied

**Symptom:** HTTP 403 responses, "permission denied" errors

**Solution:** 
1. Check token is valid: `vault token lookup`
2. Verify policy grants access: `vault policy read mirage-backend-policy`
3. Ensure AppRole has correct policy attached

### Connection Refused

**Symptom:** "connection refused" errors

**Solution:**
1. Verify Vault is running: `docker-compose ps vault`
2. Check VAULT_ADDR environment variable matches running instance
3. Ensure port 8200 is accessible

### Secret Not Found

**Symptom:** HTTP 404 when reading secret

**Solution:**
1. Verify secret exists: `vault kv list mirage/users/<user-id>`
2. Check path includes `/data/`: `/v1/mirage/data/users/...`
3. Ensure user has created the secret first

### Token Expired

**Symptom:** "token expired" or "permission denied" after working previously

**Solution:**
- Development: Use `dev-root-token` (never expires)
- Production: Application should auto-renew tokens. Check logs for renewal failures.

---

## Security Best Practices

1. **Never commit secrets to version control**
   - Use `.gitignore` for `.env` files
   - Store unseal keys and root token securely

2. **Rotate Secret IDs regularly**
   ```bash
   vault write -f auth/approle/role/mirage-backend/secret-id
   ```
   Update `VAULT_SECRET_ID` environment variable with new value.

3. **Monitor audit logs**
   - Check `/vault/logs/audit.log` for suspicious activity
   - Set up alerts for failed authentication attempts

4. **Use least privilege policies**
   - AppRole should only access paths it needs
   - Avoid using root token in production

5. **Enable MFA for UI access** (Vault Enterprise)
   - Adds extra layer of security for human access

6. **Regular backups**
   - Automated daily snapshots
   - Test restoration process

---

## Resources

- **Vault HTTP API Guide:** [docs/delivery/17/17-1-vault-http-api-guide.md](../delivery/17/17-1-vault-http-api-guide.md)
- **Environment Variables:** [ENVIRONMENT_VARIABLES.md](./ENVIRONMENT_VARIABLES.md)
- **Official Vault Documentation:** https://developer.hashicorp.com/vault/docs
- **Vault API Reference:** https://developer.hashicorp.com/vault/api-docs



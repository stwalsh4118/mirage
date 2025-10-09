# HashiCorp Vault HTTP API Implementation Guide

**Research Date:** October 9, 2025  
**Vault Version:** 1.15+ (HTTP API v1)  
**Official Documentation:** https://developer.hashicorp.com/vault/api-docs

## Overview

This guide documents the HashiCorp Vault HTTP API endpoints needed for implementing secret management in Mirage. We're using the HTTP API directly rather than the Go SDK (which is in beta) for production stability and control.

## Base URL Structure

All Vault HTTP API endpoints follow this pattern:
```
http(s)://<vault-address>:<port>/v1/<path>
```

Example: `http://localhost:8200/v1/sys/health`

## Authentication

All authenticated requests require the `X-Vault-Token` header:
```http
X-Vault-Token: <your-token>
```

For Vault Enterprise with namespaces, also include:
```http
X-Vault-Namespace: <namespace-path>
```

---

## 1. Authentication Endpoints

### 1.1 AppRole Login

**Endpoint:** `POST /v1/auth/approle/login`

**Description:** Authenticate using AppRole credentials to obtain a client token.

**Request Body:**
```json
{
  "role_id": "your-role-id",
  "secret_id": "your-secret-id"
}
```

**Response (200 OK):**
```json
{
  "auth": {
    "client_token": "hvs.CAESIJ...",
    "accessor": "accessor-id",
    "policies": ["default", "mirage-backend-policy"],
    "token_policies": ["default", "mirage-backend-policy"],
    "metadata": {},
    "lease_duration": 3600,
    "renewable": true,
    "entity_id": "entity-id"
  }
}
```

**Usage:**
```bash
curl -X POST \
  http://localhost:8200/v1/auth/approle/login \
  -d '{"role_id":"xxx","secret_id":"yyy"}'
```

**Go Implementation:**
```go
type AppRoleLoginRequest struct {
    RoleID   string `json:"role_id"`
    SecretID string `json:"secret_id"`
}

type AppRoleLoginResponse struct {
    Auth struct {
        ClientToken   string   `json:"client_token"`
        Accessor      string   `json:"accessor"`
        Policies      []string `json:"policies"`
        LeaseDuration int      `json:"lease_duration"`
        Renewable     bool     `json:"renewable"`
    } `json:"auth"`
}
```

---

### 1.2 Token Authentication (Development)

**Endpoint:** Direct token usage via header

**Description:** For development, use a static token directly in the `X-Vault-Token` header.

**Usage:**
```bash
curl -H "X-Vault-Token: dev-root-token" \
  http://localhost:8200/v1/sys/health
```

---

### 1.3 Token Lookup (Self)

**Endpoint:** `GET /v1/auth/token/lookup-self`

**Description:** Get information about the current token.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Response (200 OK):**
```json
{
  "data": {
    "accessor": "accessor-id",
    "creation_time": 1699564800,
    "creation_ttl": 3600,
    "display_name": "token",
    "entity_id": "entity-id",
    "expire_time": "2025-10-09T15:00:00Z",
    "explicit_max_ttl": 0,
    "id": "hvs.CAESIJ...",
    "issue_time": "2025-10-09T14:00:00Z",
    "meta": null,
    "num_uses": 0,
    "orphan": false,
    "path": "auth/approle/login",
    "policies": ["default", "mirage-backend-policy"],
    "renewable": true,
    "ttl": 3000
  }
}
```

**Usage:**
```bash
curl -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/auth/token/lookup-self
```

---

### 1.4 Token Renewal

**Endpoint:** `POST /v1/auth/token/renew-self`

**Description:** Renew the current token's TTL.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Request Body (optional):**
```json
{
  "increment": "1h"
}
```

**Response (200 OK):**
```json
{
  "auth": {
    "client_token": "hvs.CAESIJ...",
    "policies": ["default", "mirage-backend-policy"],
    "lease_duration": 3600,
    "renewable": true
  }
}
```

**Usage:**
```bash
curl -X POST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -d '{"increment":"1h"}' \
  http://localhost:8200/v1/auth/token/renew-self
```

---

## 2. KV Secrets Engine v2

### Path Structure

KV v2 uses two path types:
- **Data path:** `/v1/{mount}/data/{secret-path}` - for secret values
- **Metadata path:** `/v1/{mount}/metadata/{secret-path}` - for metadata/versions

For Mirage, we'll use mount path `mirage`, so:
- Data: `/v1/mirage/data/users/{user_id}/{secret_type}`
- Metadata: `/v1/mirage/metadata/users/{user_id}/{secret_type}`

---

### 2.1 Write Secret

**Endpoint:** `POST /v1/{mount}/data/{path}`

**Description:** Create a new version of a secret or create a new secret.

**Headers:**
```http
X-Vault-Token: <your-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "data": {
    "token": "railway-api-token-here",
    "metadata": {
      "created_by": "user-123",
      "created_at": "2025-10-09T14:00:00Z",
      "secret_type": "railway_token"
    }
  },
  "options": {
    "cas": 0
  }
}
```

**Parameters:**
- `data` (required): Map containing the secret data
- `options.cas` (optional): Check-And-Set value for preventing concurrent updates (0 = create only if doesn't exist)

**Response (200 OK):**
```json
{
  "data": {
    "created_time": "2025-10-09T14:00:00Z",
    "custom_metadata": null,
    "deletion_time": "",
    "destroyed": false,
    "version": 1
  }
}
```

**Usage:**
```bash
curl -X POST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -H "Content-Type: application/json" \
  -d '{"data":{"token":"xxx"}}' \
  http://localhost:8200/v1/mirage/data/users/user-123/railway
```

**Go Implementation:**
```go
type WriteSecretRequest struct {
    Data    map[string]interface{} `json:"data"`
    Options *WriteOptions          `json:"options,omitempty"`
}

type WriteOptions struct {
    CAS int `json:"cas,omitempty"`
}

type WriteSecretResponse struct {
    Data struct {
        CreatedTime    string `json:"created_time"`
        CustomMetadata map[string]string `json:"custom_metadata"`
        DeletionTime   string `json:"deletion_time"`
        Destroyed      bool   `json:"destroyed"`
        Version        int    `json:"version"`
    } `json:"data"`
}
```

---

### 2.2 Read Secret

**Endpoint:** `GET /v1/{mount}/data/{path}?version={version}`

**Description:** Read the latest version of a secret, or a specific version.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Query Parameters:**
- `version` (optional): Specific version to read (defaults to latest)

**Response (200 OK):**
```json
{
  "data": {
    "data": {
      "token": "railway-api-token-here",
      "metadata": {
        "created_by": "user-123",
        "created_at": "2025-10-09T14:00:00Z",
        "secret_type": "railway_token"
      }
    },
    "metadata": {
      "created_time": "2025-10-09T14:00:00Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  }
}
```

**Response (404 Not Found):**
```json
{
  "errors": []
}
```

**Usage:**
```bash
# Read latest version
curl -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/data/users/user-123/railway

# Read specific version
curl -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/data/users/user-123/railway?version=2
```

**Go Implementation:**
```go
type ReadSecretResponse struct {
    Data struct {
        Data     map[string]interface{} `json:"data"`
        Metadata struct {
            CreatedTime    string            `json:"created_time"`
            CustomMetadata map[string]string `json:"custom_metadata"`
            DeletionTime   string            `json:"deletion_time"`
            Destroyed      bool              `json:"destroyed"`
            Version        int               `json:"version"`
        } `json:"metadata"`
    } `json:"data"`
}
```

---

### 2.3 Delete Secret (Soft Delete)

**Endpoint:** `DELETE /v1/{mount}/data/{path}`

**Description:** Mark specific versions of a secret as deleted (soft delete). Can be undeleted.

**Headers:**
```http
X-Vault-Token: <your-token>
Content-Type: application/json
```

**Request Body (optional):**
```json
{
  "versions": [1, 2]
}
```
If no versions specified, deletes the latest version.

**Response (204 No Content)**

**Usage:**
```bash
# Delete latest version
curl -X DELETE \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/data/users/user-123/railway

# Delete specific versions
curl -X DELETE \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -H "Content-Type: application/json" \
  -d '{"versions":[1,2]}' \
  http://localhost:8200/v1/mirage/data/users/user-123/railway
```

---

### 2.4 Destroy Secret (Permanent Delete)

**Endpoint:** `POST /v1/{mount}/destroy/{path}`

**Description:** Permanently destroy specific versions of a secret. Cannot be recovered.

**Headers:**
```http
X-Vault-Token: <your-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "versions": [1, 2]
}
```

**Response (204 No Content)**

**Usage:**
```bash
curl -X POST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -H "Content-Type: application/json" \
  -d '{"versions":[1]}' \
  http://localhost:8200/v1/mirage/destroy/users/user-123/railway
```

---

### 2.5 Undelete Secret

**Endpoint:** `POST /v1/{mount}/undelete/{path}`

**Description:** Restore deleted versions of a secret.

**Headers:**
```http
X-Vault-Token: <your-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "versions": [1, 2]
}
```

**Response (204 No Content)**

**Usage:**
```bash
curl -X POST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -H "Content-Type: application/json" \
  -d '{"versions":[1]}' \
  http://localhost:8200/v1/mirage/undelete/users/user-123/railway
```

---

### 2.6 List Secrets

**Endpoint:** `LIST /v1/{mount}/metadata/{path}`

**Description:** List secret names at a given path.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Response (200 OK):**
```json
{
  "data": {
    "keys": [
      "railway",
      "github",
      "docker/"
    ]
  }
}
```

**Usage:**
```bash
# Use LIST method
curl -X LIST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/metadata/users/user-123

# Alternative: GET with ?list=true
curl -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/metadata/users/user-123?list=true
```

**Note:** Keys ending with `/` are subpaths (folders).

---

### 2.7 Read Secret Metadata

**Endpoint:** `GET /v1/{mount}/metadata/{path}`

**Description:** Get metadata about a secret, including all versions.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Response (200 OK):**
```json
{
  "data": {
    "cas_required": false,
    "created_time": "2025-10-09T14:00:00Z",
    "current_version": 3,
    "custom_metadata": {
      "owner": "user-123",
      "environment": "production"
    },
    "delete_version_after": "0s",
    "max_versions": 10,
    "oldest_version": 0,
    "updated_time": "2025-10-09T16:00:00Z",
    "versions": {
      "1": {
        "created_time": "2025-10-09T14:00:00Z",
        "deletion_time": "",
        "destroyed": false
      },
      "2": {
        "created_time": "2025-10-09T15:00:00Z",
        "deletion_time": "2025-10-09T15:30:00Z",
        "destroyed": false
      },
      "3": {
        "created_time": "2025-10-09T16:00:00Z",
        "deletion_time": "",
        "destroyed": false
      }
    }
  }
}
```

**Usage:**
```bash
curl -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/metadata/users/user-123/railway
```

**Go Implementation:**
```go
type MetadataResponse struct {
    Data struct {
        CASRequired       bool              `json:"cas_required"`
        CreatedTime       string            `json:"created_time"`
        CurrentVersion    int               `json:"current_version"`
        CustomMetadata    map[string]string `json:"custom_metadata"`
        DeleteVersionAfter string           `json:"delete_version_after"`
        MaxVersions       int               `json:"max_versions"`
        OldestVersion     int               `json:"oldest_version"`
        UpdatedTime       string            `json:"updated_time"`
        Versions          map[string]VersionInfo `json:"versions"`
    } `json:"data"`
}

type VersionInfo struct {
    CreatedTime  string `json:"created_time"`
    DeletionTime string `json:"deletion_time"`
    Destroyed    bool   `json:"destroyed"`
}
```

---

### 2.8 Update Secret Metadata (Custom Metadata)

**Endpoint:** `POST /v1/{mount}/metadata/{path}`

**Description:** Update custom metadata for a secret (doesn't create new version).

**Headers:**
```http
X-Vault-Token: <your-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "custom_metadata": {
    "owner": "user-123",
    "environment": "production",
    "last_validated": "2025-10-09T16:00:00Z"
  },
  "max_versions": 10,
  "cas_required": false
}
```

**Response (204 No Content)**

**Usage:**
```bash
curl -X POST \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  -H "Content-Type: application/json" \
  -d '{"custom_metadata":{"last_validated":"2025-10-09T16:00:00Z"}}' \
  http://localhost:8200/v1/mirage/metadata/users/user-123/railway
```

---

### 2.9 Delete Secret Metadata (Permanent)

**Endpoint:** `DELETE /v1/{mount}/metadata/{path}`

**Description:** Permanently delete all versions and metadata for a secret.

**Headers:**
```http
X-Vault-Token: <your-token>
```

**Response (204 No Content)**

**Usage:**
```bash
curl -X DELETE \
  -H "X-Vault-Token: hvs.CAESIJ..." \
  http://localhost:8200/v1/mirage/metadata/users/user-123/railway
```

**Warning:** This is a permanent operation and deletes ALL versions of the secret.

---

## 3. System Endpoints

### 3.1 Health Check

**Endpoint:** `GET /v1/sys/health`

**Description:** Check the health status of the Vault server.

**Headers:** None required (unauthenticated)

**Query Parameters:**
- `standbyok` (optional): Return 200 if standby
- `perfstandbyok` (optional): Return 200 if performance standby
- `sealedcode` (optional): HTTP code to return if sealed (default: 503)
- `uninitcode` (optional): HTTP code to return if uninitialized (default: 501)

**Response (200 OK - Healthy):**
```json
{
  "initialized": true,
  "sealed": false,
  "standby": false,
  "performance_standby": false,
  "replication_performance_mode": "disabled",
  "replication_dr_mode": "disabled",
  "server_time_utc": 1699564800,
  "version": "1.15.0",
  "cluster_name": "vault-cluster-12345",
  "cluster_id": "cluster-id"
}
```

**Response (503 Service Unavailable - Sealed):**
```json
{
  "initialized": true,
  "sealed": true,
  "standby": false
}
```

**Response (501 Not Implemented - Uninitialized):**
```json
{
  "initialized": false,
  "sealed": true,
  "standby": false
}
```

**Usage:**
```bash
curl http://localhost:8200/v1/sys/health
```

**Go Implementation:**
```go
type HealthResponse struct {
    Initialized              bool   `json:"initialized"`
    Sealed                   bool   `json:"sealed"`
    Standby                  bool   `json:"standby"`
    PerformanceStandby       bool   `json:"performance_standby"`
    ReplicationPerformanceMode string `json:"replication_performance_mode"`
    ReplicationDRMode        string `json:"replication_dr_mode"`
    ServerTimeUTC            int64  `json:"server_time_utc"`
    Version                  string `json:"version"`
    ClusterName              string `json:"cluster_name"`
    ClusterID                string `json:"cluster_id"`
}
```

---

### 3.2 Seal Status

**Endpoint:** `GET /v1/sys/seal-status`

**Description:** Get current seal status without performing a health check.

**Headers:** None required (unauthenticated)

**Response (200 OK):**
```json
{
  "type": "shamir",
  "initialized": true,
  "sealed": false,
  "t": 3,
  "n": 5,
  "progress": 0,
  "nonce": "",
  "version": "1.15.0",
  "build_date": "2023-09-01T00:00:00Z",
  "migration": false,
  "cluster_name": "vault-cluster-12345",
  "cluster_id": "cluster-id",
  "recovery_seal": false,
  "storage_type": "raft"
}
```

**Usage:**
```bash
curl http://localhost:8200/v1/sys/seal-status
```

---

## 4. Error Responses

### Standard Error Format

All errors return JSON in this format:
```json
{
  "errors": [
    "error message 1",
    "error message 2"
  ]
}
```

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request succeeded with data |
| 204 | No Content | Request succeeded, no data returned |
| 400 | Bad Request | Invalid request, missing or invalid data |
| 403 | Forbidden | Authentication failed or insufficient permissions |
| 404 | Not Found | Path doesn't exist or permission denied to check |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Vault server error |
| 501 | Not Implemented | Vault is uninitialized |
| 503 | Service Unavailable | Vault is sealed or in standby mode |

### Common Error Examples

**403 Forbidden (Permission Denied):**
```json
{
  "errors": [
    "permission denied"
  ]
}
```

**404 Not Found (Secret Not Found):**
```json
{
  "errors": []
}
```
Note: Empty errors array for 404 on read operations.

**400 Bad Request (Invalid Data):**
```json
{
  "errors": [
    "invalid secret data provided"
  ]
}
```

---

## 5. Mirage Path Conventions

### Secret Organization

All Mirage secrets will be stored under the `mirage` mount with this structure:

```
/mirage/
  └── users/
      └── {user_id}/
          ├── railway              # Railway API token
          ├── github               # GitHub PAT
          ├── docker/
          │   ├── docker.io       # Docker Hub credentials
          │   ├── ghcr.io         # GitHub Container Registry
          │   └── custom-registry # Custom registry
          ├── env_vars/
          │   ├── {env_id_1}      # Environment-specific secrets
          │   └── {env_id_2}
          └── custom/
              ├── {key_1}         # User-defined secrets
              └── {key_2}
```

### Example Paths

- Railway token: `/mirage/users/user-123/railway`
- GitHub token: `/mirage/users/user-123/github`
- Docker Hub: `/mirage/users/user-123/docker/docker.io`
- Environment secrets: `/mirage/users/user-123/env_vars/env-456`
- Custom secret: `/mirage/users/user-123/custom/my-api-key`

---

## 6. Go HTTP Client Implementation Patterns

### Basic HTTP Client

```go
package vault

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    address    string
    token      string
    httpClient *http.Client
    namespace  string
}

func NewClient(address, token string) *Client {
    return &Client{
        address: address,
        token:   token,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) makeRequest(method, path string, body interface{}, result interface{}) error {
    var reqBody io.Reader
    if body != nil {
        data, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("marshal request: %w", err)
        }
        reqBody = bytes.NewReader(data)
    }
    
    req, err := http.NewRequest(method, c.address+path, reqBody)
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }
    
    req.Header.Set("X-Vault-Token", c.token)
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    if c.namespace != "" {
        req.Header.Set("X-Vault-Namespace", c.namespace)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        var errResp struct {
            Errors []string `json:"errors"`
        }
        json.NewDecoder(resp.Body).Decode(&errResp)
        return fmt.Errorf("vault error (%d): %v", resp.StatusCode, errResp.Errors)
    }
    
    if result != nil && resp.StatusCode != 204 {
        if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
            return fmt.Errorf("decode response: %w", err)
        }
    }
    
    return nil
}
```

### Example: Write Secret

```go
func (c *Client) WriteSecret(path string, data map[string]interface{}) error {
    reqBody := map[string]interface{}{
        "data": data,
    }
    
    var resp WriteSecretResponse
    err := c.makeRequest("POST", fmt.Sprintf("/v1/mirage/data/%s", path), reqBody, &resp)
    if err != nil {
        return err
    }
    
    log.Printf("Secret written: version=%d", resp.Data.Version)
    return nil
}
```

### Example: Read Secret

```go
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
    var resp ReadSecretResponse
    err := c.makeRequest("GET", fmt.Sprintf("/v1/mirage/data/%s", path), nil, &resp)
    if err != nil {
        return nil, err
    }
    
    return resp.Data.Data, nil
}
```

---

## 7. Common Gotchas and Best Practices

### 1. KV v2 Path Structure
- **Data operations:** Use `/v1/{mount}/data/{path}`
- **Metadata operations:** Use `/v1/{mount}/metadata/{path}`
- Don't forget the `/data/` or `/metadata/` segment!

### 2. Empty Errors Array on 404
When reading a non-existent secret, Vault returns:
```json
{"errors": []}
```
Check for `404` status code, not the errors array.

### 3. Version Management
- Versions start at 1, not 0
- Latest version is always returned unless you specify `?version=N`
- Deleted versions can be undeleted (unless destroyed)
- Configure `max_versions` to limit history (recommend 10)

### 4. Check-And-Set (CAS)
Use CAS to prevent concurrent updates:
```json
{
  "data": {"key": "value"},
  "options": {"cas": 2}
}
```
This only succeeds if current version is 2. Use `cas: 0` to only write if secret doesn't exist.

### 5. Token Renewal
- Monitor token TTL from `lookup-self`
- Renew before expiration (recommend at 50% of TTL)
- Handle renewal failures gracefully

### 6. Error Handling
- Always check HTTP status codes
- Parse error messages from `errors` array
- Implement exponential backoff for 429 and 500 errors
- Cache health status to avoid excessive health checks

### 7. List Operations
- Use `LIST` method or `GET` with `?list=true`
- Keys ending with `/` are subpaths
- List is not recursive - list each level separately

### 8. Performance
- Cache frequently accessed secrets (with appropriate TTL)
- Use connection pooling in HTTP client
- Batch operations when possible
- Monitor Vault server load

---

## 8. Summary of Key Endpoints for Mirage

### Authentication
- `POST /v1/auth/approle/login` - AppRole login
- `GET /v1/auth/token/lookup-self` - Check token info
- `POST /v1/auth/token/renew-self` - Renew token

### Secret Operations
- `POST /v1/mirage/data/{path}` - Write secret (creates new version)
- `GET /v1/mirage/data/{path}?version={n}` - Read secret (optionally specific version)
- `DELETE /v1/mirage/data/{path}` - Soft delete secret
- `POST /v1/mirage/destroy/{path}` - Permanently destroy secret
- `POST /v1/mirage/undelete/{path}` - Undelete secret
- `LIST /v1/mirage/metadata/{path}` - List secrets at path

### Metadata Operations
- `GET /v1/mirage/metadata/{path}` - Get secret metadata and version history
- `POST /v1/mirage/metadata/{path}` - Update custom metadata
- `DELETE /v1/mirage/metadata/{path}` - Permanently delete all versions

### System Operations
- `GET /v1/sys/health` - Health check
- `GET /v1/sys/seal-status` - Seal status

---

## References

- Official Vault API Documentation: https://developer.hashicorp.com/vault/api-docs
- KV v2 Secrets Engine: https://developer.hashicorp.com/vault/api-docs/secret/kv/kv-v2
- AppRole Auth: https://developer.hashicorp.com/vault/api-docs/auth/approle
- System Endpoints: https://developer.hashicorp.com/vault/api-docs/system
- OpenAPI Explorer: Available in Vault UI by typing `api` in console

---

**Last Updated:** October 9, 2025  
**Maintainer:** Mirage Development Team


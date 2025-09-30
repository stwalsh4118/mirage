# Railway API Guide: Docker Image Deployments

**Created:** 2025-09-30  
**Source:** Railway GraphQL API v2 (https://backboard.railway.com/graphql/v2)  
**Reference:** Railway API introspection and documentation

This guide documents Railway's GraphQL API support for creating services from Docker images, complementing the existing source repository deployment approach.

## Overview

Railway supports two primary service deployment sources:
1. **Source Repository** (`source: { repo: "..." }`) - Deploy from GitHub/GitLab repos
2. **Docker Image** (`source: { image: "..." }`) - Deploy from container registries

## ServiceCreate Mutation for Docker Images

### Mutation Signature

```graphql
mutation ServiceCreate($input: ServiceCreateInput!) {
    serviceCreate(input: $input) {
        id
        name
        projectId
        createdAt
        updatedAt
        icon
        featureFlags
        templateServiceId
        templateThreadSlug
        deletedAt
    }
}
```

### ServiceCreateInput Structure for Image Deployments

```graphql
input ServiceCreateInput {
    projectId: String!              # Required: Railway project ID
    environmentId: String           # Optional: Scope to specific environment
    name: String                    # Optional: Service name (auto-generated if not provided)
    source: ServiceSourceInput      # Source configuration (repo OR image)
    registryCredentials: RegistryCredentialsInput  # Optional: For private registries
    icon: String                    # Optional: Service icon
    branch: String                  # Optional: For repo deployments
    variables: EnvironmentVariables # Optional: Environment variables
    templateServiceId: String       # Optional: Clone from template
}

input ServiceSourceInput {
    repo: String   # For source repo deployments (e.g., "github.com/owner/repo")
    image: String  # For Docker image deployments (e.g., "nginx:latest", "ghcr.io/owner/image:v1.0")
}

input RegistryCredentialsInput {
    username: String  # Registry username (e.g., Docker Hub username, GitHub username)
    password: String  # Registry password or access token
}
```

## Supported Container Registries

Railway supports deploying from any Docker-compatible registry:

- **Docker Hub** (`docker.io` or no prefix)
  - Public: `nginx:latest`, `redis:7-alpine`
  - Private: `username/private-image:tag`
- **GitHub Container Registry (GHCR)** (`ghcr.io`)
  - Format: `ghcr.io/owner/repository:tag`
- **GitLab Container Registry** (`registry.gitlab.com`)
  - Format: `registry.gitlab.com/group/project:tag`
- **Quay.io** (`quay.io`)
  - Format: `quay.io/organization/repository:tag`
- **Custom registries** (any Docker-compatible registry)
  - Format: `registry.example.com/image:tag`

## Image Reference Formats

### Public Images

**Docker Hub (short form):**
```
nginx:latest
postgres:15
redis:7-alpine
```

**Docker Hub (full form):**
```
docker.io/library/nginx:latest
docker.io/username/custom-app:v1.0
```

**GHCR:**
```
ghcr.io/owner/repository:latest
ghcr.io/owner/repository:sha-abc123
```

### Using Digests (for reproducibility)

Instead of tags, use SHA256 digests for immutable deployments:
```
nginx@sha256:abcdef123456...
ghcr.io/owner/image@sha256:fedcba654321...
```

## Example Mutations

### Example 1: Deploy Public Docker Hub Image

```graphql
mutation CreateNginxService {
    serviceCreate(
        input: {
            projectId: "01234567-89ab-cdef-0123-456789abcdef"
            environmentId: "abcdef12-3456-7890-abcd-ef1234567890"
            name: "web-server"
            source: { image: "nginx:latest" }
        }
    ) {
        id
        name
        projectId
        createdAt
    }
}
```

### Example 2: Deploy from GHCR Public Image

```graphql
mutation CreateGHCRService {
    serviceCreate(
        input: {
            projectId: "01234567-89ab-cdef-0123-456789abcdef"
            environmentId: "abcdef12-3456-7890-abcd-ef1234567890"
            name: "api-service"
            source: { image: "ghcr.io/organization/api-image:v2.1.0" }
        }
    ) {
        id
        name
        projectId
    }
}
```

### Example 3: Deploy Private Docker Hub Image

```graphql
mutation CreatePrivateDockerHubService {
    serviceCreate(
        input: {
            projectId: "01234567-89ab-cdef-0123-456789abcdef"
            environmentId: "abcdef12-3456-7890-abcd-ef1234567890"
            name: "private-app"
            source: { image: "username/private-image:latest" }
            registryCredentials: {
                username: "dockerhub-username"
                password: "dckr_pat_abc123..."  # Docker Hub access token
            }
        }
    ) {
        id
        name
        projectId
    }
}
```

### Example 4: Deploy Private GHCR Image

```graphql
mutation CreatePrivateGHCRService {
    serviceCreate(
        input: {
            projectId: "01234567-89ab-cdef-0123-456789abcdef"
            environmentId: "abcdef12-3456-7890-abcd-ef1234567890"
            name: "private-ghcr-app"
            source: { image: "ghcr.io/owner/private-repo:latest" }
            registryCredentials: {
                username: "github-username"
                password: "ghp_abc123..."  # GitHub Personal Access Token with read:packages scope
            }
        }
    ) {
        id
        name
        projectId
    }
}
```

### Example 5: Deploy with Environment Variables

```graphql
mutation CreateServiceWithEnvVars {
    serviceCreate(
        input: {
            projectId: "01234567-89ab-cdef-0123-456789abcdef"
            environmentId: "abcdef12-3456-7890-abcd-ef1234567890"
            name: "configured-service"
            source: { image: "redis:7-alpine" }
            variables: {
                REDIS_PASSWORD: "super-secret-password"
                REDIS_PORT: "6379"
            }
        }
    ) {
        id
        name
        projectId
    }
}
```

## Security Considerations

### Credential Handling

**❌ NEVER:**
- Store registry credentials in your application database
- Log credentials in application logs
- Expose credentials in frontend code
- Commit credentials to version control

**✅ ALWAYS:**
- Use HTTPS for all API communication
- Pass credentials directly from frontend → backend → Railway
- Let Railway handle credential encryption and storage
- Use access tokens instead of passwords when possible

**Credential Flow:**
```
User Input (HTTPS)
    ↓
Your Backend API (HTTPS) - Pass-through only, no persistence
    ↓
Railway API (HTTPS) - Railway securely encrypts and stores credentials
```

### Access Token Best Practices

**Docker Hub:**
- Use Docker Hub Access Tokens (not account password)
- Create at: https://hub.docker.com/settings/security
- Scope: Read-only access is sufficient for pulling images

**GitHub (GHCR):**
- Use Personal Access Tokens (PAT) or fine-grained tokens
- Required scope: `read:packages`
- Create at: https://github.com/settings/tokens

**GitLab:**
- Use Deploy Tokens or Personal Access Tokens
- Required scope: `read_registry`

## Image Deployment Limitations

Based on Railway's infrastructure:

1. **Image Size:** Large images (>10GB) may experience slower deployment times
2. **Architecture:** Railway runs on `linux/amd64` - ensure your images support this architecture
3. **Multi-architecture Images:** Railway will automatically select the `amd64` variant
4. **Private Registries:** Require authentication on every service creation
5. **Image Pull Timeout:** Railway has timeouts for image pulls (typically 10-15 minutes)

## Error Handling

### Common Errors

**Image Not Found:**
```json
{
  "errors": [
    {
      "message": "Failed to pull image: image not found",
      "extensions": {
        "code": "IMAGE_PULL_ERROR"
      }
    }
  ]
}
```
**Solution:** Verify image name and tag exist in registry

**Authentication Failed:**
```json
{
  "errors": [
    {
      "message": "Failed to authenticate with registry",
      "extensions": {
        "code": "REGISTRY_AUTH_ERROR"
      }
    }
  ]
}
```
**Solution:** Verify credentials are correct and have necessary permissions

**Invalid Image Format:**
```json
{
  "errors": [
    {
      "message": "Invalid image reference format",
      "extensions": {
        "code": "INVALID_IMAGE_REFERENCE"
      }
    }
  ]
}
```
**Solution:** Ensure image follows format: `[registry/]repository[:tag|@digest]`

## Go Client Implementation Example

```go
// ServiceSourceInput can specify either repo or image
type ServiceSourceInput struct {
    Repo  *string `json:"repo,omitempty"`
    Image *string `json:"image,omitempty"`
}

// RegistryCredentials for private registries
type RegistryCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// CreateServiceInput extended for image deployments
type CreateServiceInput struct {
    ProjectID           string               `json:"projectId"`
    EnvironmentID       string               `json:"environmentId,omitempty"`
    Name                string               `json:"name,omitempty"`
    Source              *ServiceSourceInput  `json:"source,omitempty"`
    RegistryCredentials *RegistryCredentials `json:"registryCredentials,omitempty"`
    Variables           map[string]string    `json:"variables,omitempty"`
}

// Example: Create service from Docker Hub image
func CreateImageBasedService(client *railway.Client, projectID, envID string) error {
    image := "nginx:latest"
    input := CreateServiceInput{
        ProjectID:     projectID,
        EnvironmentID: envID,
        Name:          "nginx-server",
        Source: &ServiceSourceInput{
            Image: &image,
        },
    }
    
    result, err := client.CreateService(context.Background(), input)
    if err != nil {
        return fmt.Errorf("failed to create service: %w", err)
    }
    
    log.Printf("Created service: %s", result.ServiceID)
    return nil
}

// Example: Create service from private GHCR image
func CreatePrivateImageService(client *railway.Client, projectID, envID, username, token string) error {
    image := "ghcr.io/owner/private-image:v1.0"
    input := CreateServiceInput{
        ProjectID:     projectID,
        EnvironmentID: envID,
        Name:          "private-service",
        Source: &ServiceSourceInput{
            Image: &image,
        },
        RegistryCredentials: &RegistryCredentials{
            Username: username,
            Password: token,
        },
    }
    
    result, err := client.CreateService(context.Background(), input)
    if err != nil {
        return fmt.Errorf("failed to create service: %w", err)
    }
    
    log.Printf("Created service: %s", result.ServiceID)
    return nil
}
```

## Comparison: Source Repo vs Docker Image

| Aspect | Source Repo | Docker Image |
|--------|-------------|--------------|
| Build | Railway builds from source | Pre-built, Railway pulls image |
| Build Time | Varies (1-10+ min) | Fast (~30s-2min) |
| Dockerfile | Optional (Nixpacks auto-detects) | Pre-built, not used |
| Source Control | Required (GitHub/GitLab) | Not required |
| CI/CD Integration | Railway handles build | External CI builds image |
| Use Case | Source code in repo | Pre-built artifacts, third-party images |

## Credential Storage and Reuse

### How Railway Stores Credentials

When you provide `registryCredentials` in a `serviceCreate` mutation:

1. **Railway stores the credentials** - They are encrypted and associated with the service
2. **Automatic reuse** - Railway automatically uses these credentials for all image pulls (deployments, redeployments)
3. **No re-sending needed** - You only provide credentials once during service creation
4. **Scoped to service** - Each service has its own credential storage

### Using Shared Environment Variables

Railway also supports using project/environment-level variables for registry authentication:

**Set these variables in Railway Dashboard:**
- `DOCKER_USERNAME` and `DOCKER_PASSWORD` - For Docker Hub
- `GITHUB_TOKEN` - For GitHub Container Registry (GHCR)
- `GITLAB_TOKEN` - For GitLab Container Registry

**Then create services WITHOUT inline credentials:**
```graphql
serviceCreate(
    input: {
        projectId: "..."
        source: { image: "username/private-image:latest" }
        # Railway automatically detects and uses DOCKER_USERNAME/DOCKER_PASSWORD
    }
)
```

**Benefits of shared variables:**
- Set credentials once for all services
- Easy credential rotation (update env var, redeploy)
- No credentials in API calls
- Centralized management

**Trade-offs:**
- All services in project/environment share same credentials
- Requires manual setup in Railway UI first
- Less explicit than inline credentials

### Recommended Approach for Mirage

**For MVP:**
Use inline `registryCredentials` in the mutation:
- Explicit and predictable
- No manual Railway setup required
- Per-service credential control

**For Future Enhancement:**
Add option to use Railway shared variables:
- User sets credentials once in Railway
- Mirage references them for all private image deployments
- Simpler UX for repeated use

## Testing Checklist

- [ ] Create service from Docker Hub public image (`nginx:latest`)
- [ ] Create service from GHCR public image
- [ ] Create service from Docker Hub private image (with inline credentials)
- [ ] Verify credentials are stored (trigger manual redeploy, check it still works)
- [ ] Create service from GHCR private image (with inline credentials)
- [ ] Test using shared env vars (DOCKER_USERNAME/DOCKER_PASSWORD)
- [ ] Verify service deployment succeeds in Railway
- [ ] Test with image tag (`:latest`, `:v1.0`)
- [ ] Test with image digest (`@sha256:...`)
- [ ] Test error handling for non-existent images
- [ ] Test error handling for invalid credentials
- [ ] Test error handling for malformed image references
- [ ] Verify credentials persist across redeployments

## Additional Resources

- **Railway Documentation:** https://docs.railway.com/reference/services
- **Railway API Endpoint:** https://backboard.railway.com/graphql/v2
- **Docker Hub:** https://hub.docker.com
- **GHCR Documentation:** https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry

## Notes

1. Railway stores registry credentials securely - they are encrypted at rest
2. Credentials are scoped to the specific service and not shared across services
3. For repeated use of the same private image, consider storing credentials as Railway secrets/environment variables
4. Railway automatically retries failed image pulls with exponential backoff
5. Image deployments skip the build phase, making them faster than source deployments


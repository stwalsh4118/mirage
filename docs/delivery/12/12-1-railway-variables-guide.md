# Railway Service Variables Guide

**Created**: 2025-10-01  
**Task**: 12-1  
**Purpose**: Document how to set Railway service variables, specifically `RAILWAY_DOCKERFILE_PATH`, via GraphQL API

## Overview

Railway services support environment variables that control build behavior. The key variable for custom Dockerfile paths is `RAILWAY_DOCKERFILE_PATH`, which tells Railway where to find the Dockerfile in a repository.

## Source Documentation

- Railway Docs: [Docker Builds](https://docs.railway.com/guides/dockerfiles)
- User-provided example from Railway GraphQL API response

## Setting Variables During Service Creation

### Mutation Structure

The `serviceCreate` mutation accepts a `variables` field that is a simple map/object of string key-value pairs:

```graphql
mutation ServiceCreate($input: ServiceCreateInput!) {
  serviceCreate(input: $input) {
    id
    name
    projectId
    updatedAt
  }
}
```

### Variables Input

The `variables` field in `ServiceCreateInput` is a map where:
- **Keys**: Environment variable names (strings)
- **Values**: Environment variable values (strings)
- **Type**: `{ [key: string]: string }` or `map[string]string`

### Complete Example

```json
{
  "input": {
    "projectId": "01234567-89ab-cdef-0123-456789abcdef",
    "environmentId": "abcdef01-2345-6789-abcd-ef0123456789",
    "name": "api-service",
    "source": {
      "repo": "https://github.com/owner/monorepo"
    },
    "branch": "main",
    "variables": {
      "RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile",
      "NODE_ENV": "production",
      "BUILD_VERSION": "1.0.0"
    }
  }
}
```

### Go Implementation Pattern

```go
input := map[string]any{
    "projectId":     projectID,
    "environmentId": environmentID,
    "name":          serviceName,
    "source": map[string]any{
        "repo": repoURL,
    },
    "branch": branchName,
}

// Add variables if present
if len(variables) > 0 {
    input["variables"] = variables
}

vars := map[string]any{
    "input": input,
}
```

## RAILWAY_DOCKERFILE_PATH Specification

### Variable Name

```
RAILWAY_DOCKERFILE_PATH
```

### Path Format

**Relative paths from repository root:**

✅ **Valid Examples:**
- `Dockerfile` (root of repo)
- `services/api/Dockerfile` (subdirectory)
- `apps/backend/prod.dockerfile` (subdirectory with custom name)
- `packages/auth/Dockerfile` (nested structure)

❌ **Invalid Examples:**
- `/services/api/Dockerfile` (absolute path - remove leading slash)
- `../Dockerfile` (parent directory traversal)
- `C:\apps\Dockerfile` (absolute Windows path)

### Build Context Behavior

When `RAILWAY_DOCKERFILE_PATH` is set:
- Railway uses the **directory containing the Dockerfile** as the build context
- Example: If `RAILWAY_DOCKERFILE_PATH=services/api/Dockerfile`
  - Build context: `services/api/`
  - Dockerfile: `services/api/Dockerfile`

This means relative paths in the Dockerfile are relative to `services/api/`, not the repo root.

## Other Relevant Railway Build Variables

### Standard Railway Variables

Railway automatically provides these variables (read-only):
```
RAILWAY_PROJECT_ID
RAILWAY_PROJECT_NAME
RAILWAY_ENVIRONMENT_ID
RAILWAY_ENVIRONMENT_NAME
RAILWAY_SERVICE_ID
RAILWAY_SERVICE_NAME
RAILWAY_STATIC_URL
RAILWAY_PRIVATE_DOMAIN
```

### User-Configurable Build Variables

Variables you can set to customize builds:

| Variable | Purpose | Example |
|----------|---------|---------|
| `RAILWAY_DOCKERFILE_PATH` | Custom Dockerfile location | `services/api/Dockerfile` |
| `RAILWAY_HEALTHCHECK_TIMEOUT_SEC` | Health check timeout | `300` |
| `NIXPACKS_*` | Nixpacks build configuration | Various |

**Note**: As of October 2025, Railway does not have a separate `RAILWAY_DOCKERFILE_BUILD_CONTEXT` variable. The build context is automatically inferred from the Dockerfile's directory.

## Updating Variables on Existing Services

While this task focuses on service creation, variables can also be updated using the `variableUpsert` mutation (if needed in future tasks):

```graphql
mutation VariableUpsert($input: VariableUpsertInput!) {
  variableUpsert(input: $input)
}
```

Example input:
```json
{
  "input": {
    "projectId": "project-uuid",
    "environmentId": "env-uuid",
    "serviceId": "service-uuid",
    "name": "RAILWAY_DOCKERFILE_PATH",
    "value": "services/api/Dockerfile"
  }
}
```

## Variable Resolution Example

Given this mutation:
```json
{
  "variables": {
    "RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile",
    "PORT": "3000"
  }
}
```

Railway will:
1. Clone the repository
2. Look for Dockerfile at `services/api/Dockerfile`
3. Use `services/api/` as the build context
4. Build the image using the specified Dockerfile
5. Set `PORT=3000` as an environment variable in the running container

## Error Handling

### Common Errors

**Invalid Path:**
```json
{
  "error": "Dockerfile not found at specified path"
}
```
- Verify the path is correct relative to repo root
- Ensure the file exists in the repository

**Build Context Issues:**
```json
{
  "error": "COPY failed: file not found"
}
```
- Check that files referenced in Dockerfile exist relative to build context
- Remember build context is the Dockerfile's directory, not repo root

## Implementation Checklist

For task 12-3, ensure:
- [ ] `variables` field added to `CreateServiceInput` struct
- [ ] Variables are passed as `map[string]string` or equivalent
- [ ] Variables included in GraphQL mutation conditionally (only if non-empty)
- [ ] Dockerfile path validation (relative, no traversal)
- [ ] Path uses forward slashes (even on Windows)
- [ ] Empty variables map doesn't break mutation
- [ ] Logging shows which variables are being set

## Examples for Testing

### Test Case 1: Single Service with Dockerfile
```json
{
  "name": "api",
  "repo": "github.com/org/monorepo",
  "branch": "main",
  "variables": {
    "RAILWAY_DOCKERFILE_PATH": "services/api/Dockerfile"
  }
}
```

### Test Case 2: Multiple Variables
```json
{
  "name": "worker",
  "repo": "github.com/org/monorepo",
  "branch": "main",
  "variables": {
    "RAILWAY_DOCKERFILE_PATH": "services/worker/prod.dockerfile",
    "NODE_ENV": "production",
    "WORKER_CONCURRENCY": "10"
  }
}
```

### Test Case 3: No Custom Variables
```json
{
  "name": "web",
  "repo": "github.com/org/simple-app",
  "branch": "main"
  // No variables field - Railway uses default Dockerfile in root
}
```

## Additional Notes

### Variable Precedence

1. **Service-level variables** (what we're setting) override environment-level variables
2. **Environment-level variables** override project-level variables
3. **Project-level variables** are the base defaults

### Security Considerations

- Variables are stored encrypted by Railway
- Sensitive values (API keys, passwords) should use Railway's secrets management
- `RAILWAY_DOCKERFILE_PATH` is not sensitive - it's just a path

### Performance

- Setting variables during service creation is atomic - no separate API call needed
- Variables are immediately available to the build process
- No significant performance impact from additional variables

## References

- [Railway Dockerfile Guide](https://docs.railway.com/guides/dockerfiles)
- [Railway GraphQL API](https://docs.railway.com/reference/api)
- User-provided Railway API response examples


# Railway API Guide: Environment Variables Operations

**Created:** 2025-10-06  
**Task:** 15-1  
**Purpose:** Document Railway's GraphQL API for fetching and setting environment variables to support environment cloning (PBI 15)

## Overview

Railway's GraphQL API provides operations to manage environment variables at different scopes (project, environment, and service levels). For environment cloning, we need to:
1. **Fetch** all environment variables from a source environment
2. **Upsert** environment variables to a target environment

This guide documents the specific GraphQL operations and patterns needed for these operations.

## Source Documentation

- **Railway Public API**: https://docs.railway.com/guides/public-api
- **Railway API Endpoint**: https://backboard.railway.com/graphql/v2
- **Authentication**: Bearer token in `Authorization` header
- **Task Context**: PBI 15 (Environment Cloning) - PBI 13 established that env vars are NOT stored locally for security

## Railway Variables Hierarchy

Railway organizes environment variables in a hierarchical structure:

```
Project
├── Project-level variables (shared across all environments)
├── Environment
│   ├── Environment-level variables (override project vars)
│   └── Services
│       └── Service-level variables (override environment vars)
```

### Variable Precedence (Highest to Lowest)
1. **Service-level variables** - Specific to a service in an environment
2. **Environment-level variables** - Shared by all services in an environment
3. **Project-level variables** - Base defaults for all environments

For environment cloning, we primarily focus on **environment-level** and **service-level** variables.

## Fetching Environment Variables

### Query Structure

Railway provides a simple root-level `variables` query that returns a flat map of environment variables. The query accepts optional `projectId` and `environmentId` parameters to scope the results.

### Query: Fetch Variables

```graphql
query Variables($projectId: String, $environmentId: String, $serviceId: String) {
  variables(projectId: $projectId, environmentId: $environmentId, serviceId: $serviceId)
}
```

### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | String | No | Filter to project-level variables when provided alone |
| `environmentId` | String | No | Filter to environment-level variables (requires projectId) |
| `serviceId` | String | No | Filter to service-level variables (requires environmentId) |

### Parameter Combinations

1. **All variables for a project**:
   ```graphql
   variables(projectId: "proj-123", environmentId: null, serviceId: null)
   ```

2. **All variables for an environment** (environment + project variables):
   ```graphql
   variables(projectId: "proj-123", environmentId: "env-456", serviceId: null)
   ```

3. **All variables for a service** (service + environment + project variables):
   ```graphql
   variables(projectId: "proj-123", environmentId: "env-456", serviceId: "svc-789")
   ```

### Response Format

The query returns a flat JSON object where keys are variable names and values are variable values:

```json
{
  "data": {
    "variables": {
      "VARIABLE_NAME_1": "value1",
      "VARIABLE_NAME_2": "value2",
      "DATABASE_URL": "postgresql://localhost:5432/db",
      "PORT": "3000"
    }
  }
}
```

### Example Response (Railway System Variables)

```json
{
  "data": {
    "variables": {
      "RAILWAY_PUBLIC_DOMAIN": "myapp.example.com",
      "RAILWAY_PRIVATE_DOMAIN": "myapp.railway.internal",
      "RAILWAY_PROJECT_NAME": "my-project",
      "RAILWAY_ENVIRONMENT_NAME": "production",
      "RAILWAY_SERVICE_NAME": "api-service",
      "RAILWAY_PROJECT_ID": "01234567-89ab-cdef-0123-456789abcdef",
      "RAILWAY_ENVIRONMENT_ID": "abcdef01-2345-6789-abcd-ef0123456789",
      "RAILWAY_SERVICE_ID": "fedcba98-7654-3210-fedc-ba9876543210",
      "RAILWAY_STATIC_URL": "myapp.example.com",
      "RAILWAY_ENVIRONMENT": "production",
      "RAILWAY_SERVICE_API_SERVICE_URL": "myapp.example.com"
    }
  }
}
```

**Note**: The above example shows Railway system variables (RAILWAY_*) which are automatically injected. User-defined variables would appear in the same flat structure.

## Setting Environment Variables

### Mutation: variableCollectionUpsert

Railway provides a batch upsert mutation for setting multiple variables at once. This is the recommended approach for cloning.

```graphql
mutation UpsertVariables(
  $input: VariableCollectionUpsertInput!
) {
  variableCollectionUpsert(input: $input)
}
```

### VariableCollectionUpsertInput Structure

```graphql
input VariableCollectionUpsertInput {
  projectId: String!           # Required: Railway project ID
  environmentId: String!       # Required: Target environment ID
  serviceId: String            # Optional: Scope to specific service
  variables: VariableInput!    # Required: Variables to upsert
  replace: Boolean             # Optional: If true, replace all existing variables
}

# Variables is a JSON object/map: { "KEY": "value", ... }
scalar VariableInput
```

### Key Parameters

- **`projectId`**: The Railway project ID
- **`environmentId`**: The target environment where variables will be set
- **`serviceId`**: Optional - if provided, variables are service-scoped; if omitted, they're environment-scoped
- **`variables`**: A JSON object/map where keys are variable names and values are variable values
- **`replace`**: If `true`, removes all existing variables and replaces with provided set. If `false` (default), merges/updates existing variables

### Example 1: Upsert Environment-Level Variables

```graphql
mutation UpsertEnvironmentVariables {
  variableCollectionUpsert(
    input: {
      projectId: "proj-123-abc"
      environmentId: "env-456-def"
      variables: {
        DATABASE_URL: "postgresql://user:pass@host:5432/newdb"
        LOG_LEVEL: "debug"
        API_KEY: "secret-key-xyz"
      }
    }
  )
}
```

**Result**: Variables are set at environment level (all services can access them)

### Example 2: Upsert Service-Level Variables

```graphql
mutation UpsertServiceVariables {
  variableCollectionUpsert(
    input: {
      projectId: "proj-123-abc"
      environmentId: "env-456-def"
      serviceId: "svc-789-ghi"
      variables: {
        PORT: "8080"
        RAILWAY_DOCKERFILE_PATH: "services/api/Dockerfile"
        NODE_ENV: "production"
      }
    }
  )
}
```

**Result**: Variables are set at service level (only this service can access them)

### Example 3: Replace All Environment Variables

```graphql
mutation ReplaceAllEnvironmentVariables {
  variableCollectionUpsert(
    input: {
      projectId: "proj-123-abc"
      environmentId: "env-456-def"
      replace: true
      variables: {
        NEW_VAR_1: "value1"
        NEW_VAR_2: "value2"
      }
    }
  )
}
```

**Result**: All existing environment-level variables are removed, replaced with the new set

### Response

The mutation returns a `Boolean`:
- `true`: Variables were successfully upserted
- GraphQL error: Operation failed (e.g., invalid project/environment ID, authentication issues)

```json
{
  "data": {
    "variableCollectionUpsert": true
  }
}
```

## Alternative Mutation: variableUpsert (Single Variable)

For setting individual variables (less efficient for cloning, but useful for reference):

```graphql
mutation UpsertSingleVariable(
  $input: VariableUpsertInput!
) {
  variableUpsert(input: $input)
}
```

### VariableUpsertInput Structure

```graphql
input VariableUpsertInput {
  projectId: String!        # Required
  environmentId: String!    # Required
  serviceId: String         # Optional
  name: String!             # Variable name/key
  value: String!            # Variable value
}
```

### Example

```graphql
mutation {
  variableUpsert(
    input: {
      projectId: "proj-123"
      environmentId: "env-456"
      name: "DATABASE_URL"
      value: "postgresql://localhost:5432/db"
    }
  )
}
```

**Note**: For cloning with many variables, use `variableCollectionUpsert` instead to reduce API calls.

## Variable Types and Special Considerations

### System Variables (Read-Only)

Railway automatically provides these variables (cannot be overridden):

```
RAILWAY_PROJECT_ID
RAILWAY_PROJECT_NAME
RAILWAY_ENVIRONMENT_ID
RAILWAY_ENVIRONMENT_NAME
RAILWAY_SERVICE_ID
RAILWAY_SERVICE_NAME
RAILWAY_STATIC_URL
RAILWAY_PRIVATE_DOMAIN
RAILWAY_PUBLIC_DOMAIN
RAILWAY_GIT_COMMIT_SHA
RAILWAY_GIT_BRANCH
```

**Cloning Behavior**: Do NOT attempt to copy these during cloning - Railway will generate new values for the target environment.

### Build-Time vs Runtime Variables

All variables are available at both build time and runtime by default. Some special variables affect build behavior:

| Variable | Scope | Purpose |
|----------|-------|---------|
| `RAILWAY_DOCKERFILE_PATH` | Service | Custom Dockerfile location |
| `RAILWAY_HEALTHCHECK_TIMEOUT_SEC` | Service | Health check timeout |
| `NIXPACKS_*` | Service | Nixpacks build configuration |

### Multiline Variables

Variables can span multiple lines (e.g., SSH keys, certificates). These are stored with newline characters (`\n`) in the value string:

```json
{
  "SSH_PRIVATE_KEY": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCA...\n-----END RSA PRIVATE KEY-----",
  "PEM_CERTIFICATE": "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAw...\n-----END CERTIFICATE-----"
}
```

**Cloning Behavior**: Multiline variables are preserved automatically since they're just strings with embedded newlines. No special handling needed.

### Secret References

Railway supports secret references using the format `${{ secrets.SECRET_NAME }}`. These reference project-level secrets.

**Cloning Behavior**: 
- Secret references can be copied as-is IF the target project has the same secrets defined
- Consider warning users if secret references are detected during cloning
- May need to be updated to match target project's secret names

## Go Client Implementation Patterns

### Type Definitions

```go
package railway

// GetEnvironmentVariablesInput contains parameters for fetching variables
type GetEnvironmentVariablesInput struct {
    ProjectID     string
    EnvironmentID string
    ServiceID     *string // nil for environment-level, set for service-level
}

// GetEnvironmentVariablesResult contains the fetched variables as a flat map
type GetEnvironmentVariablesResult struct {
    Variables map[string]string // Variable name -> value
}

// UpsertVariablesInput contains parameters for upserting variables
type UpsertVariablesInput struct {
    ProjectID     string
    EnvironmentID string
    ServiceID     *string // nil for environment-level
    Variables     map[string]string
    Replace       bool // If true, replace all existing variables
}
```

### Fetching Variables

```go
// GetEnvironmentVariables fetches all variables for an environment or service
// Returns a flat map of variable name -> value
func (c *Client) GetEnvironmentVariables(
    ctx context.Context,
    in GetEnvironmentVariablesInput,
) (GetEnvironmentVariablesResult, error) {
    query := gqlGetEnvironmentVariables // embedded .graphql file
    
    vars := map[string]any{
        "projectId":     in.ProjectID,
        "environmentId": in.EnvironmentID,
    }
    
    // Add serviceId if requesting service-level variables
    if in.ServiceID != nil {
        vars["serviceId"] = *in.ServiceID
    } else {
        vars["serviceId"] = nil
    }
    
    var resp struct {
        Variables map[string]string `json:"variables"`
    }
    
    if err := c.execute(ctx, query, vars, &resp); err != nil {
        return GetEnvironmentVariablesResult{}, fmt.Errorf("fetch environment variables: %w", err)
    }
    
    return GetEnvironmentVariablesResult{
        Variables: resp.Variables,
    }, nil
}
```

### Upserting Variables

```go
// UpsertEnvironmentVariables sets or updates environment variables
func (c *Client) UpsertEnvironmentVariables(
    ctx context.Context,
    in UpsertVariablesInput,
) error {
    mutation := gqlUpsertVariables // embedded .graphql file
    
    input := map[string]any{
        "projectId":     in.ProjectID,
        "environmentId": in.EnvironmentID,
        "variables":     in.Variables,
    }
    
    // Add optional fields
    if in.ServiceID != nil {
        input["serviceId"] = *in.ServiceID
    }
    if in.Replace {
        input["replace"] = true
    }
    
    vars := map[string]any{
        "input": input,
    }
    
    var resp struct {
        VariableCollectionUpsert bool `json:"variableCollectionUpsert"`
    }
    
    if err := c.execute(ctx, mutation, vars, &resp); err != nil {
        return fmt.Errorf("upsert environment variables: %w", err)
    }
    
    if !resp.VariableCollectionUpsert {
        return fmt.Errorf("variableCollectionUpsert returned false")
    }
    
    return nil
}
```

### Example Usage: Cloning Variables

```go
func CloneEnvironmentVariables(
    client *railway.Client,
    sourceProjectID string,
    sourceEnvID string,
    targetProjectID string,
    targetEnvID string,
) error {
    ctx := context.Background()
    
    // Step 1: Fetch all variables from source environment
    log.Info().
        Str("source_env_id", sourceEnvID).
        Msg("fetching source environment variables")
    
    result, err := client.GetEnvironmentVariables(ctx, railway.GetEnvironmentVariablesInput{
        ProjectID:     sourceProjectID,
        EnvironmentID: sourceEnvID,
        ServiceID:     nil, // Get environment-level variables
    })
    if err != nil {
        return fmt.Errorf("fetch source variables: %w", err)
    }
    
    // Step 2: Filter out Railway system variables
    filteredVars := make(map[string]string)
    systemVarCount := 0
    
    for name, value := range result.Variables {
        // Skip Railway system variables (they'll be auto-generated in target)
        if strings.HasPrefix(name, "RAILWAY_") {
            systemVarCount++
            continue
        }
        filteredVars[name] = value
    }
    
    log.Info().
        Int("total_variables", len(result.Variables)).
        Int("system_variables_skipped", systemVarCount).
        Int("user_variables", len(filteredVars)).
        Msg("filtered environment variables")
    
    // Step 3: Upsert filtered variables to target environment
    if len(filteredVars) > 0 {
        log.Info().
            Int("count", len(filteredVars)).
            Str("target_env_id", targetEnvID).
            Msg("upserting environment-level variables")
        
        err = client.UpsertEnvironmentVariables(ctx, railway.UpsertVariablesInput{
            ProjectID:     targetProjectID,
            EnvironmentID: targetEnvID,
            ServiceID:     nil, // Environment-level
            Variables:     filteredVars,
            Replace:       false, // Merge with any existing variables
        })
        if err != nil {
            return fmt.Errorf("upsert environment variables: %w", err)
        }
    }
    
    log.Info().
        Int("variables_cloned", len(filteredVars)).
        Msg("environment variables cloned successfully")
    
    return nil
}

// CloneServiceVariables clones variables for a specific service
func CloneServiceVariables(
    client *railway.Client,
    sourceProjectID string,
    sourceEnvID string,
    sourceServiceID string,
    targetProjectID string,
    targetEnvID string,
    targetServiceID string,
) error {
    ctx := context.Background()
    
    // Fetch service-level variables from source
    result, err := client.GetEnvironmentVariables(ctx, railway.GetEnvironmentVariablesInput{
        ProjectID:     sourceProjectID,
        EnvironmentID: sourceEnvID,
        ServiceID:     &sourceServiceID,
    })
    if err != nil {
        return fmt.Errorf("fetch source service variables: %w", err)
    }
    
    // Filter out Railway system variables
    filteredVars := make(map[string]string)
    for name, value := range result.Variables {
        if !strings.HasPrefix(name, "RAILWAY_") {
            filteredVars[name] = value
        }
    }
    
    if len(filteredVars) == 0 {
        log.Info().
            Str("source_service_id", sourceServiceID).
            Msg("no user-defined service variables to clone")
        return nil
    }
    
    // Upsert to target service
    log.Info().
        Str("source_service_id", sourceServiceID).
        Str("target_service_id", targetServiceID).
        Int("count", len(filteredVars)).
        Msg("upserting service-level variables")
    
    err = client.UpsertEnvironmentVariables(ctx, railway.UpsertVariablesInput{
        ProjectID:     targetProjectID,
        EnvironmentID: targetEnvID,
        ServiceID:     &targetServiceID,
        Variables:     filteredVars,
        Replace:       false,
    })
    if err != nil {
        return fmt.Errorf("upsert service variables: %w", err)
    }
    
    return nil
}
```

## Error Handling

### Common Errors

**1. Environment Not Found**
```json
{
  "errors": [
    {
      "message": "Environment not found",
      "extensions": { "code": "NOT_FOUND" }
    }
  ]
}
```
**Solution**: Verify environment ID is correct

**2. Insufficient Permissions**
```json
{
  "errors": [
    {
      "message": "Insufficient permissions to access environment",
      "extensions": { "code": "FORBIDDEN" }
    }
  ]
}
```
**Solution**: Ensure API token has access to the project/environment

**3. Invalid Variable Format**
```json
{
  "errors": [
    {
      "message": "Variable names must not contain spaces or special characters",
      "extensions": { "code": "INVALID_INPUT" }
    }
  ]
}
```
**Solution**: Validate variable names before upserting (alphanumeric + underscore only)

**4. Project/Environment Mismatch**
```json
{
  "errors": [
    {
      "message": "Environment does not belong to specified project",
      "extensions": { "code": "INVALID_INPUT" }
    }
  ]
}
```
**Solution**: Ensure projectId and environmentId match

### Go Error Handling Pattern

```go
func (c *Client) UpsertEnvironmentVariables(
    ctx context.Context,
    in UpsertVariablesInput,
) error {
    // Validate inputs
    if in.ProjectID == "" {
        return fmt.Errorf("projectId is required")
    }
    if in.EnvironmentID == "" {
        return fmt.Errorf("environmentId is required")
    }
    if len(in.Variables) == 0 {
        return fmt.Errorf("variables map is empty")
    }
    
    // Validate variable names
    for name := range in.Variables {
        if !isValidVariableName(name) {
            return fmt.Errorf("invalid variable name: %q (must be alphanumeric + underscore)", name)
        }
    }
    
    // Execute mutation with retry logic (handled by c.execute)
    // ... implementation
    
    return nil
}

func isValidVariableName(name string) bool {
    if name == "" {
        return false
    }
    for _, ch := range name {
        if !((ch >= 'A' && ch <= 'Z') || 
             (ch >= 'a' && ch <= 'z') || 
             (ch >= '0' && ch <= '9') || 
             ch == '_') {
            return false
        }
    }
    return true
}
```

## Security Considerations

### 1. Never Store Variables Locally

**Critical**: Following PBI 13's design decision:
- ✅ Fetch variables from Railway API when needed
- ✅ Pass variables through backend → Railway API
- ❌ NEVER store variables in local database
- ❌ NEVER log variable values (log keys only)

### 2. Credentials in Transit

- All API communication must use HTTPS
- Variables pass through: `Frontend (HTTPS) → Backend (HTTPS) → Railway API`
- Backend should not persist variables at any point

### 3. Logging Variables

```go
// ❌ DON'T: Log variable values
log.Info().
    Str("DATABASE_URL", dbURL). // NEVER log sensitive values
    Msg("setting variables")

// ✅ DO: Log variable keys only
log.Info().
    Strs("variable_keys", variableKeys).
    Int("count", len(variables)).
    Msg("upserting environment variables")
```

### 4. Secret References

Variables containing secret references (`${{ secrets.SECRET_NAME }}`):
- Are safe to store/log (they're references, not values)
- May need adjustment when cloning to different projects
- Should be documented in clone operation

## Testing Strategy

### Unit Tests

```go
func TestUpsertEnvironmentVariables(t *testing.T) {
    tests := []struct {
        name      string
        input     UpsertVariablesInput
        expectErr bool
    }{
        {
            name: "valid environment variables",
            input: UpsertVariablesInput{
                ProjectID:     "proj-123",
                EnvironmentID: "env-456",
                ServiceID:     nil,
                Variables: map[string]string{
                    "DATABASE_URL": "postgresql://localhost:5432/db",
                    "PORT":         "3000",
                    "LOG_LEVEL":    "info",
                },
            },
            expectErr: false,
        },
        {
            name: "valid service variables",
            input: UpsertVariablesInput{
                ProjectID:     "proj-123",
                EnvironmentID: "env-456",
                ServiceID:     strPtr("svc-789"),
                Variables: map[string]string{
                    "RAILWAY_DOCKERFILE_PATH": "Dockerfile",
                    "PORT":                     "8080",
                },
            },
            expectErr: false,
        },
        {
            name: "multiline variable",
            input: UpsertVariablesInput{
                ProjectID:     "proj-123",
                EnvironmentID: "env-456",
                Variables: map[string]string{
                    "SSH_KEY": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAK...\n-----END RSA PRIVATE KEY-----",
                },
            },
            expectErr: false,
        },
        {
            name: "invalid variable name with spaces",
            input: UpsertVariablesInput{
                ProjectID:     "proj-123",
                EnvironmentID: "env-456",
                Variables: map[string]string{
                    "INVALID NAME": "value", // Contains space
                },
            },
            expectErr: true,
        },
        {
            name: "empty project ID",
            input: UpsertVariablesInput{
                ProjectID:     "",
                EnvironmentID: "env-456",
                Variables:     map[string]string{"KEY": "value"},
            },
            expectErr: true,
        },
        {
            name: "empty variables map",
            input: UpsertVariablesInput{
                ProjectID:     "proj-123",
                EnvironmentID: "env-456",
                Variables:     map[string]string{},
            },
            expectErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := client.UpsertEnvironmentVariables(context.Background(), tt.input)
            if (err != nil) != tt.expectErr {
                t.Errorf("expected error=%v, got error=%v", tt.expectErr, err)
            }
        })
    }
}

func strPtr(s string) *string {
    return &s
}
```

### Integration Tests

```go
func TestGetEnvironmentVariables_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    client := railway.NewFromConfig(config.Load())
    ctx := context.Background()
    
    // Use a real test environment ID
    projectID := os.Getenv("RAILWAY_TEST_PROJECT_ID")
    envID := os.Getenv("RAILWAY_TEST_ENV_ID")
    if projectID == "" || envID == "" {
        t.Skip("RAILWAY_TEST_PROJECT_ID or RAILWAY_TEST_ENV_ID not set")
    }
    
    result, err := client.GetEnvironmentVariables(ctx, railway.GetEnvironmentVariablesInput{
        ProjectID:     projectID,
        EnvironmentID: envID,
        ServiceID:     nil, // Get environment-level variables
    })
    
    if err != nil {
        t.Fatalf("GetEnvironmentVariables failed: %v", err)
    }
    
    // Verify structure (don't assert specific values)
    t.Logf("Found %d variables", len(result.Variables))
    
    // Count Railway system vs user variables
    systemVars := 0
    userVars := 0
    for name := range result.Variables {
        if strings.HasPrefix(name, "RAILWAY_") {
            systemVars++
        } else {
            userVars++
        }
    }
    
    t.Logf("System variables: %d", systemVars)
    t.Logf("User variables: %d", userVars)
    
    // Verify Railway system variables are present
    requiredSystemVars := []string{
        "RAILWAY_PROJECT_ID",
        "RAILWAY_ENVIRONMENT_ID",
        "RAILWAY_PROJECT_NAME",
        "RAILWAY_ENVIRONMENT_NAME",
    }
    
    for _, varName := range requiredSystemVars {
        if _, ok := result.Variables[varName]; !ok {
            t.Errorf("Expected system variable %q not found", varName)
        }
    }
}

func TestUpsertAndGetVariables_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    client := railway.NewFromConfig(config.Load())
    ctx := context.Background()
    
    projectID := os.Getenv("RAILWAY_TEST_PROJECT_ID")
    envID := os.Getenv("RAILWAY_TEST_ENV_ID")
    if projectID == "" || envID == "" {
        t.Skip("RAILWAY_TEST_PROJECT_ID or RAILWAY_TEST_ENV_ID not set")
    }
    
    // Test variables
    testVars := map[string]string{
        "TEST_VAR_1": "value1",
        "TEST_VAR_2": "value2",
        "TEST_MULTILINE": "line1\nline2\nline3",
    }
    
    // Upsert test variables
    err := client.UpsertEnvironmentVariables(ctx, railway.UpsertVariablesInput{
        ProjectID:     projectID,
        EnvironmentID: envID,
        ServiceID:     nil,
        Variables:     testVars,
        Replace:       false,
    })
    if err != nil {
        t.Fatalf("UpsertEnvironmentVariables failed: %v", err)
    }
    
    // Fetch and verify
    result, err := client.GetEnvironmentVariables(ctx, railway.GetEnvironmentVariablesInput{
        ProjectID:     projectID,
        EnvironmentID: envID,
        ServiceID:     nil,
    })
    if err != nil {
        t.Fatalf("GetEnvironmentVariables failed: %v", err)
    }
    
    // Verify test variables were set correctly
    for name, expectedValue := range testVars {
        actualValue, ok := result.Variables[name]
        if !ok {
            t.Errorf("Variable %q not found after upsert", name)
            continue
        }
        if actualValue != expectedValue {
            t.Errorf("Variable %q: expected %q, got %q", name, expectedValue, actualValue)
        }
    }
}
```

## Implementation Checklist for Tasks 15-2 and 15-3

### Task 15-2: Implement GetEnvironmentVariables
- [ ] Create `queries/queries/variables-get.graphql` file with variables query
- [ ] Add `//go:embed` directive in `api/internal/railway/variables.go`
- [ ] Implement `GetEnvironmentVariables()` method with simplified types
- [ ] Add type definitions: `GetEnvironmentVariablesInput`, `GetEnvironmentVariablesResult`
- [ ] Handle `projectId`, `environmentId`, and optional `serviceId` parameters
- [ ] Parse flat map response: `map[string]string`
- [ ] Add unit tests for input validation
- [ ] Add integration test with real Railway environment
- [ ] Document usage in `api/internal/railway/variables.go`
- [ ] Export types for use in other packages

### Task 15-3: Implement UpsertEnvironmentVariables
- [ ] Create `queries/mutations/variables-upsert.graphql` file with variableCollectionUpsert mutation
- [ ] Add `//go:embed` directive in `api/internal/railway/variables.go`
- [ ] Implement `UpsertEnvironmentVariables()` method
- [ ] Add input validation (variable name format: alphanumeric + underscore only)
- [ ] Validate required fields: `projectId`, `environmentId`, non-empty `variables`
- [ ] Support optional `serviceId` for service-level variables
- [ ] Support optional `replace` flag (default false)
- [ ] Add helper function to filter Railway system variables (RAILWAY_* prefix)
- [ ] Add unit tests with various scenarios (env-level, service-level, validation)
- [ ] Add integration test: upsert variables and verify via GET
- [ ] Add error handling for common Railway API errors
- [ ] Document usage with examples in `api/internal/railway/variables.go`

## Summary

This guide provides the foundation for implementing environment variable operations in tasks 15-2 and 15-3:

- **Fetching**: Use root-level `variables(projectId, environmentId, serviceId)` query that returns a flat `map[string]string`
- **Setting**: Use `variableCollectionUpsert` mutation with `projectId`, `environmentId`, optional `serviceId`, and `variables` map
- **Response Format**: Simple flat JSON object, not relay-style connections - much simpler than initially assumed
- **Scoping**: Variables can be scoped to project, environment, or service level via query parameters
- **System Variables**: Railway automatically injects `RAILWAY_*` variables - filter these out when cloning
- **Security**: Never store variables locally; Railway API is the source of truth (per PBI 13 design)
- **Cloning**: Fetch from source → Filter out RAILWAY_* → Upsert to target
- **Go Patterns**: Follow existing Railway client patterns with embed, retry, and structured error handling

## References

- **Railway Public API**: https://docs.railway.com/guides/public-api
- **Railway API Endpoint**: https://backboard.railway.com/graphql/v2
- **Related PBI**: PBI 13 (Service Build Configuration Management)
- **Related PBI**: PBI 15 (Environment Cloning)
- **Related Tasks**: 15-2 (Implement GET), 15-3 (Implement UPSERT)


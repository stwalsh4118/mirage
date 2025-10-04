# PBI-15: Environment Cloning

[View in Backlog](../backlog.md#user-content-15)

## Overview
Enable users to clone an existing environment with all its services, configuration, and metadata. This allows teams to quickly create replicas of environments for different purposes (e.g., create a staging environment from dev, spin up a test environment from prod), significantly reducing setup time and ensuring consistency.

## Problem Statement
Creating similar environments is currently a manual, time-consuming, and error-prone process. Users must:
- Recreate each service individually
- Manually copy environment variables
- Reconstruct service relationships
- Risk configuration drift between environments
- Spend significant time on repetitive setup

Environment cloning eliminates this toil by automating the duplication process while allowing customization for the target environment type.

## User Stories
- As a developer, I want to clone my dev environment to create a staging environment so I can test in production-like conditions
- As a platform engineer, I want to clone environments to different Railway projects for organizational separation
- As a developer, I want to customize the cloned environment's configuration (env type, TTL, env vars) before creation
- As a developer, I want cloning to preserve service relationships and dependencies
- As a developer, I want to see which environment a clone was created from for traceability

## Technical Approach

### Backend (Go)

#### Clone Operation
- **Deep Copy Strategy**: Create new entities with references to original
  - Generate new IDs for environment and services
  - Copy all configuration from `EnvironmentMetadata` (PBI 13)
  - Preserve service relationships via `ServiceDependency` table
  - Set `ClonedFromEnvID` in metadata to track lineage

#### Clone Process Steps
1. **Validation**: Verify source environment exists and is accessible
2. **Metadata Retrieval**: Fetch complete environment metadata (PBI 13)
3. **Environment Variables Fetching**:
   - Fetch env vars from source environment via Railway API
   - **Note**: Env vars are NOT stored in local database for security (see PBI 13)
   - Railway API is the source of truth for all environment variables
   - User can preview and modify env vars before applying to target
4. **Configuration Transformation**:
   - Apply target environment type (dev → staging, etc.)
   - Apply user overrides to environment variables
   - Modify resource allocations if specified
   - Update TTL settings
5. **Railway Provisioning**:
   - Create Railway project (if target is new project)
   - Create Railway environment
   - Create all services with cloned configuration
   - Apply environment variables via Railway API
   - Establish service connections
6. **Metadata Storage**:
   - Save new environment to database
   - Record clone relationship (`ClonedFromEnvID`)
   - Store build configurations
   - Save service dependencies
   - **Note**: Do NOT store environment variables locally

#### API Endpoints
- `POST /api/environments/:id/clone` - Clone environment
  - Request body: target project, environment type, customizations
- `GET /api/environments/:id/clone-preview` - Preview clone operation
  - Returns what will be created without executing
- `GET /api/environments/:id/clones` - List environments cloned from this one
- `GET /api/environments/:id/lineage` - Get clone lineage (ancestors and descendants)

#### Error Handling
- Rollback on partial failure (atomic operation)
- Clear error messages for Railway API failures
- Validation errors before starting clone operation
- Track clone operation progress for long-running operations

### Frontend (Next.js/React)

#### Clone Wizard
- **Trigger**: "Clone Environment" button on environment detail page
- **Step 1: Clone Source**:
  - Display source environment details
  - Show services that will be cloned
  - Display current configuration summary
- **Step 2: Target Configuration**:
  - Project selection: existing or new
  - Environment name input
  - Environment type selector (dev, staging, prod, ephemeral)
  - TTL configuration
- **Step 3: Customization**:
  - Service selection (include/exclude specific services)
  - Environment variable overrides
  - Resource allocation adjustments
  - Service name transformations (optional suffixes/prefixes)
- **Step 4: Review**:
  - Side-by-side comparison of source vs. target config
  - Diff view for changed values
  - Estimated cost (if available)
  - Confirmation prompt
- **Step 5: Execution**:
  - Progress indicator (similar to creation wizard)
  - Per-service status updates
  - Success/failure feedback
  - Link to new environment

#### Environment Lineage Visualization
- **Clone Badge**: Show badge on cloned environments
- **Lineage View**: Visual tree showing clone relationships
  - "Cloned from" link on environment detail page
  - "Clones" section showing derived environments
- **Diff View**: Compare configuration between original and clone

### Integration with PBI 13
- Depends on `EnvironmentMetadata` table for full config storage
- Uses `ServiceBuildConfig` for service-level settings
- Leverages `ServiceDependency` graph for relationship preservation

## UX/UI Considerations

### Clone Button Placement
- Primary action button on environment detail page
- Disabled state if environment is not in healthy state
- Tooltip explaining what will be cloned

### Clone Wizard Design
- Similar visual style to creation wizard (consistency)
- Clear progress indicators
- Back/Next navigation
- Escape hatch: "Cancel" at any step before execution
- Auto-fill smart defaults based on source environment

### Configuration Diff View
- Two-column layout: "Source" | "Target"
- Color coding:
  - Green: New values
  - Yellow: Modified values
  - Red: Removed values
  - Gray: Unchanged values
- Expand/collapse sections for readability

### Clone Lineage Display
- Tree diagram with environment cards
- Connecting lines showing clone relationships
- Timestamp of clone operation
- Hover to see clone metadata

### Error Handling UX
- Validation errors: inline, before execution
- Execution errors: detailed error modal with:
  - What succeeded (partial clone)
  - What failed and why
  - Option to retry or rollback
  - Link to Railway logs if applicable

## Acceptance Criteria
1. "Clone Environment" action available on environment detail page
2. Clone wizard guides user through configuration steps
3. User can select target project (new or existing)
4. User can customize environment name, type, and TTL
5. User can include/exclude specific services
6. User can override environment variables before cloning
7. Environment variables are fetched from Railway API (not local database)
8. Preview step shows configuration diff including env var changes
9. Clone operation creates new Railway project (if specified)
10. Clone operation creates new Railway environment
11. All services from source are recreated in target
12. Service build configurations are preserved (from PBI 13)
13. Service dependencies are maintained
14. Environment variables are copied to target via Railway API (with overrides applied)
15. No environment variables are stored in local database
16. Cloned environment has `ClonedFromEnvID` set correctly
17. Clone lineage is visible on environment detail pages
18. Clone operation is atomic (all-or-nothing with rollback)
19. Progress indicator shows per-service status during clone
20. Success feedback provides link to new environment
21. Error handling provides actionable feedback

## Dependencies
- PBI 13 (Service Build Configuration Management) - REQUIRED for metadata storage
- PBI 10 (Environment Creation Wizard) - UI patterns and provisioning logic
- PBI 11 (Docker Image Service Deployment) - image-based service cloning
- Complete Railway API integration for environment/service creation
- Railway API support for fetching environment variables (GET /environments/:id/variables)
- Railway API support for setting environment variables (POST /environments/:id/variables)

## Open Questions
1. Should cloning be synchronous or asynchronous (job queue)?
2. How do we handle cloning environments with active deployments?
3. Should we support partial clones (subset of services)?
4. ~~How do we handle secrets during cloning (exclude, encrypt, prompt)?~~ **RESOLVED**: Environment variables are fetched from Railway API during clone operation. No local storage of secrets. Clone flow: Fetch vars from source env → Create target env → Copy vars to target env via Railway API.
5. Should users be able to clone from other users' environments (RBAC)?
6. Do we need clone operation history/audit log?
7. Should we support "clone and update" (clone + modify in one operation)?
8. How do we handle cross-region cloning if Railway supports multiple regions?

## Related Tasks
Tasks will be created once this PBI moves to "Agreed" status.


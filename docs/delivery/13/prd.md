# PBI-13: Service Build Configuration and Environment Metadata Management

[View in Backlog](../backlog.md#user-content-13)

## Overview
Implement comprehensive storage of service build configurations and environment metadata to enable environment duplication, templating, and reproducible deployments across different environment types (dev, staging, prod). This PBI provides the foundation for creating environment templates from existing environments and cloning environments with full fidelity.

## Problem Statement
Currently, Mirage stores minimal service information (name, path, status). To support advanced workflows like environment cloning (PBI 15) and multi-environment deployments, we need to capture and persist the complete build and runtime configuration of services. Without this metadata:
- Environments cannot be reliably cloned with identical configuration
- Creating staging/prod replicas of dev environments requires manual reconfiguration
- Build configurations must be re-entered for each environment
- Service dependencies and relationships are not tracked

## User Stories
- As a platform engineer, I want all service build configurations persisted so I can recreate environments exactly
- As a developer, I want to save an environment as a template and reuse it for different environment types
- As a developer, I want to export environment configuration and import it into a different project
- As a platform engineer, I want service dependencies tracked so environments maintain correct relationships

## Technical Approach

### Data Model Extensions

#### ServiceBuildConfig Table
```go
type ServiceBuildConfig struct {
    ID                  string    `gorm:"primaryKey"`
    ServiceID           string    `gorm:"index;not null"`
    
    // Deployment Source
    DeploymentType      string    `gorm:"not null"` // "source_repo", "docker_image"
    
    // Source Repository Config (if DeploymentType = source_repo)
    RepoURL             *string   `gorm:"type:text"`
    RepoBranch          *string   `gorm:"type:text"`
    RepoCommit          *string   `gorm:"type:text"`
    RootDirectory       *string   `gorm:"type:text"`
    
    // Docker Build Config
    DockerfilePath      *string   `gorm:"type:text"`
    BuildContext        *string   `gorm:"type:text"`
    BuildArgsJSON       string    `gorm:"type:text"` // JSON map of build args
    TargetStage         *string   `gorm:"type:text"` // Multi-stage build target
    
    // Docker Image Config (if DeploymentType = docker_image)
    ImageRegistry       *string   `gorm:"type:text"`
    ImageName           *string   `gorm:"type:text"`
    ImageTag            *string   `gorm:"type:text"`
    ImageDigest         *string   `gorm:"type:text"`
    
    // Runtime Config
    ExposedPorts        string    `gorm:"type:text"` // JSON array of port numbers
    HealthCheckPath     *string   `gorm:"type:text"`
    StartCommand        *string   `gorm:"type:text"`
    
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

#### EnvironmentMetadata Table
```go
type EnvironmentMetadata struct {
    ID                  string    `gorm:"primaryKey"`
    EnvironmentID       string    `gorm:"index;not null"`
    
    // Template Information
    IsTemplate          bool      `gorm:"default:false"`
    TemplateName        *string   `gorm:"type:text"`
    TemplateDescription *string   `gorm:"type:text"`
    
    // Source Information (for cloned environments)
    ClonedFromEnvID     *string   `gorm:"type:text"`
    
    // Configuration
    ConfigJSON          string    `gorm:"type:text"` // Full env config as JSON
    ServiceDepsJSON     string    `gorm:"type:text"` // Service dependency graph
    
    // Metadata
    Tags                string    `gorm:"type:text"` // JSON array of tags
    Labels              string    `gorm:"type:text"` // JSON map of labels
    
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

#### ServiceDependency Table
```go
type ServiceDependency struct {
    ID                  string    `gorm:"primaryKey"`
    ServiceID           string    `gorm:"index;not null"`
    DependsOnServiceID  string    `gorm:"index;not null"`
    DependencyType      string    `gorm:"not null"` // "required", "optional"
    
    CreatedAt           time.Time
}
```

### Backend (Go)

#### Configuration Capture
- **Service Creation Hook**: Capture full config when services are created
- **Configuration Snapshots**: Periodic snapshots of live service configuration
- **Change Detection**: Track when configurations change

#### API Endpoints
- `POST /api/environments/:id/metadata` - Store/update environment metadata
- `GET /api/environments/:id/metadata` - Retrieve full metadata
- `POST /api/environments/:id/export` - Export environment as JSON
- `POST /api/environments/import` - Import environment from JSON
- `POST /api/environments/:id/save-as-template` - Save environment as reusable template
- `GET /api/templates` - List available environment templates
- `GET /api/services/:id/build-config` - Get service build configuration
- `POST /api/services/:id/build-config` - Update service build configuration

#### Service Dependency Tracking
- Analyze service configurations to detect dependencies
- Support explicit dependency declaration
- Validate dependency graph (detect cycles)

### Frontend (Next.js/React)

#### Environment Actions
- **"Save as Template"** button on environment detail page
  - Modal to name and describe template
  - Option to include/exclude sensitive data
  - Confirmation and success feedback
- **"Export Configuration"** action
  - Download JSON file with full environment config
  - Include/exclude options for selective export
- **"View Metadata"** panel
  - Display all captured configuration
  - Show service dependencies graph
  - Display template information if applicable

#### Template Management
- New "Templates" section in dashboard
- List of saved templates with metadata
- Preview template configuration
- Use template in creation wizard

## UX/UI Considerations

### Save as Template Flow
1. User clicks "Save as Template" on environment detail page
2. Modal appears with:
   - Template name input
   - Description textarea
   - Checklist of what to include (env vars, build config, dependencies)
   - Option to mark sensitive variables for exclusion
3. Confirmation: "Template saved successfully"
4. Template appears in templates list

### Export/Import Flow
- **Export**: Downloads JSON file named `{environment-name}-config.json`
- **Import**: Drag-and-drop or file picker in wizard
  - Preview imported configuration
  - Map to target project/environment
  - Resolve conflicts (missing secrets, etc.)

### Metadata Visualization
- **Dependency Graph**: Visual graph showing service relationships
- **Configuration Diff**: Compare configurations across environments
- **Template Badge**: Visual indicator on environments created from templates

## Acceptance Criteria
1. All service build configurations are persisted on service creation
2. Dockerfile paths and build contexts are stored
3. Build arguments are captured and stored as JSON
4. Docker image references (registry, name, tag, digest) are stored
5. Service ports and health check configs are persisted
6. Environment metadata table stores template information
7. Service dependencies can be explicitly declared
8. Dependency graph is validated (no cycles)
9. "Save as Template" functionality works end-to-end
10. Templates can be listed and viewed
11. Environment configuration can be exported as JSON
12. Exported configuration can be imported to create new environment
13. Cloned environments reference original via metadata
14. Tags and labels can be added to environments
15. API endpoints support full CRUD operations on metadata

## Dependencies
- PBI 11 (Docker Image Service Deployment) - image config storage
- PBI 12 (Monorepo Dockerfile Discovery) - Dockerfile metadata capture
- Database schema migration capability
- JSON schema validation for import/export

## Open Questions
1. Should templates be user-specific or organization-wide?
2. ~~How do we handle secrets in templates (encrypt, exclude, tokenize)?~~ **RESOLVED**: Environment variables (including secrets) are NOT stored in the database for security. Railway API is the source of truth. During cloning (PBI 15), env vars will be fetched from Railway API on-demand.
3. Should we version templates (v1, v2, etc.)?
4. How do we handle Railway-specific IDs during import?
5. Should dependency detection be automatic or manual?
6. Do we need a template marketplace/sharing feature?
7. How do we handle configuration drift detection?

## Related Tasks
Tasks will be created once this PBI moves to "Agreed" status.


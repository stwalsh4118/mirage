# PBI-11: Docker Image Service Deployment

[View in Backlog](../backlog.md#user-content-11)

## Overview
Enable users to create services from pre-built Docker images instead of requiring source repository deployments. This allows teams to deploy services from container registries like Docker Hub, GitHub Container Registry (GHCR), or private registries, providing more deployment flexibility and supporting existing containerized workflows.

## Problem Statement
Currently, Mirage only supports creating services from source repositories, requiring Railway to build from source. Teams with existing Docker image pipelines, pre-built containers, or third-party images cannot leverage these assets in Mirage-managed environments. This limitation prevents users from:
- Deploying pre-built, tested images from CI/CD pipelines
- Using third-party containerized services (databases, caching, etc.)
- Separating build and deployment concerns
- Leveraging multi-stage builds performed outside Railway

## User Stories
- As a developer, I can create a service from a Docker image URL so that I can deploy pre-built containers
- As a platform engineer, I can specify image tags or digests to ensure reproducible deployments
- As a developer, I can configure environment variables and port mappings for image-based services
- As a developer, I can use images from private registries with authentication

## Technical Approach

### Backend (Go)
- **Data Model Extensions**: Add fields to `Service` model:
  - `DeploymentType` (enum: `source_repo`, `docker_image`)
  - `DockerImage` (string: full image reference)
  - `ImageRegistry` (string: registry URL)
  - `ImageTag` (string: tag or digest)
  - `ImageAuthRequired` (bool: whether auth is needed)
- **Railway Integration**: Extend `railway.Client` to support image-based service creation:
  - Use Railway's service creation API with image source
  - Handle registry authentication via Railway's secrets
  - Support both public and private registries
- **API Endpoints**: New endpoints for image-based service creation
  - `POST /api/services/from-image` - Create service from Docker image
  - `GET /api/services/:id/image-config` - Retrieve image configuration

### Frontend (Next.js/React)
- **Wizard Enhancement**: Add "Deployment Source" step:
  - Radio selection: "From Source Repository" or "From Docker Image"
  - Image input fields: registry, image name, tag/digest
  - Registry authentication (optional)
  - Port and environment variable configuration
- **Service Card Updates**: Display deployment type indicator on service cards
- **Validation**: Client-side validation of image reference format

### Railway API Integration
- Research Railway's GraphQL mutations for image-based deployments
- Document supported registries and authentication methods
- Implement image accessibility validation before deployment

## UX/UI Considerations

### Deployment Source Selection
- Clear visual distinction between source and image deployments
- Helper text with examples: `docker.io/nginx:latest`, `ghcr.io/owner/image:v1.0`
- Support for both tag and digest-based references

### Image Configuration Form
- **Registry dropdown**: Docker Hub, GHCR, Custom
- **Image reference input**: Auto-format and validate
- **Tag/Digest selector**: Latest vs specific version
- **Port configuration**: Default port detection with override
- **Environment variables**: Key-value editor
- **Authentication toggle**: Show/hide registry credentials

### Service Dashboard
- Badge indicating "Image" vs "Source" deployment
- Image reference display with registry icon
- Tag/digest information in service details

## Acceptance Criteria
1. Users can select "Docker Image" as deployment source in wizard
2. Services can be created from Docker Hub public images
3. Services can be created from GHCR public images
4. Image reference validation prevents invalid formats
5. Custom registry URLs are supported
6. Port configuration works for image-based services
7. Environment variables can be set for image-based services
8. Service cards display deployment type
9. Image accessibility is validated before deployment attempt
10. Error messages clearly indicate image-related issues

## Dependencies
- PBI 10 (Environment Creation Wizard) must be complete
- Railway API documentation for image-based deployments
- Understanding of Railway's registry authentication mechanisms

## Open Questions
1. Does Railway support all major container registries out of the box?
2. How should we handle image pull authentication for private registries?
3. Should we support image scanning/vulnerability checking?
4. Do we need to store registry credentials or rely on Railway's secret management?
5. How do we handle multi-architecture images (amd64 vs arm64)?

## Related Tasks
Tasks will be created once this PBI moves to "Agreed" status.


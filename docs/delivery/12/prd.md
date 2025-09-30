# PBI-12: Monorepo Dockerfile Discovery

[View in Backlog](../backlog.md#user-content-12)

## Overview
Automatically discover Dockerfiles within monorepo subdirectories and enable multi-service deployment using the discovered Docker build configurations. This extends PBI 2 (Monorepo Service Discovery) to include container-based deployments and complements PBI 11 (Docker Image Service Deployment).

## Problem Statement
Many teams structure their monorepos with per-service Dockerfiles (e.g., `services/api/Dockerfile`, `services/worker/Dockerfile`). Currently, users must manually configure each service's build settings. This PBI automates the discovery of these Dockerfiles and intelligently configures services to build from them, reducing manual configuration and preventing misconfiguration errors.

## User Stories
- As a platform engineer, I want Mirage to automatically find Dockerfiles in my monorepo so I don't have to manually specify build paths
- As a developer, I want to see which services have Dockerfiles and select which ones to deploy
- As a developer, I want Mirage to detect build context and arguments from Dockerfiles
- As a platform engineer, I want to deploy multiple services from their respective Dockerfiles in a single operation

## Technical Approach

### Backend (Go)
- **Dockerfile Scanner**: Create service discovery module
  - Recursively scan repository structure for `Dockerfile` and `*.dockerfile` files
  - Parse Dockerfile for metadata (EXPOSE ports, ENV vars, build args)
  - Detect build context (typically parent directory of Dockerfile)
  - Identify multi-stage builds and extract relevant stages
- **Service Boundary Detection**:
  - Heuristics: `package.json`, `go.mod`, `requirements.txt` in same directory
  - Support conventional structures (`services/`, `apps/`, `packages/`)
  - Group Dockerfile with related service code
- **Data Model**: Extend service discovery results
  - `DockerfileDetected` (bool)
  - `DockerfilePath` (string: relative path from repo root)
  - `BuildContext` (string: directory for build context)
  - `ExposedPorts` ([]int: ports from EXPOSE directives)
  - `BuildArgs` (map[string]string: detected ARG directives)
- **API Endpoints**:
  - `POST /api/discovery/dockerfiles` - Scan repo for Dockerfiles
  - `GET /api/discovery/dockerfiles/:repo` - Get cached discovery results

### Frontend (Next.js/React)
- **Discovery Results View**: Show discovered Dockerfiles in service selection
  - Tree view showing repository structure
  - Checkbox selection for services to deploy
  - Dockerfile metadata display (ports, build args)
  - Preview of detected build configuration
- **Wizard Integration**: Add Dockerfile discovery step
  - Automatic scan on repository selection
  - Visual indicators for services with Dockerfiles
  - Option to override detected settings
- **Configuration Preview**: Show effective build config before deployment

### Integration with PBI 11
- When Dockerfile is detected, offer choice:
  - Build from Dockerfile in Railway
  - Build locally/in CI and deploy as image (PBI 11)
- Support both workflows from same discovery results

## UX/UI Considerations

### Discovery Results Display
- **Tree View**: Hierarchical display of repository structure
  - Folder icons for directories
  - Docker icon for discovered Dockerfiles
  - Badge indicating service type (Node, Go, Python, etc.)
- **Service Cards**: Each discovered service shows:
  - Service name (derived from directory)
  - Dockerfile path
  - Detected ports
  - Detected build arguments
  - Estimated build context size
- **Selection Interface**:
  - "Select All" / "Select None" buttons
  - Smart defaults (select all by default)
  - Warning for services without Dockerfiles

### Configuration Override
- Inline editing of detected configuration
- Override build context path
- Modify build arguments
- Add/remove ports
- Set custom service names

## Acceptance Criteria
1. Scanner recursively finds all `Dockerfile` files in repository
2. Scanner detects `*.dockerfile` variants (e.g., `production.dockerfile`)
3. Build context is correctly identified for each Dockerfile
4. EXPOSE directives are parsed and extracted
5. ARG directives are detected and presented for override
6. UI displays discovered services in tree structure
7. Users can select/deselect services for deployment
8. Configuration can be overridden before deployment
9. Multi-service deployment creates all selected services
10. Discovery results are cached and can be refreshed
11. Error handling for malformed Dockerfiles

## Dependencies
- PBI 2 (Monorepo Service Discovery) - foundation for detection logic
- PBI 11 (Docker Image Service Deployment) - image deployment support
- PBI 13 (Service Build Configuration Management) - storing discovered config
- Access to repository contents for scanning

## Open Questions
1. Should we support Docker Compose files for service discovery?
2. How do we handle Dockerfiles with dynamic build arguments?
3. Should we validate Dockerfiles or just parse for metadata?
4. How do we detect service dependencies from Dockerfiles?
5. Should discovery run client-side (clone repo) or server-side?
6. How do we handle monorepos with 50+ services?

## Related Tasks
Tasks will be created once this PBI moves to "Agreed" status.


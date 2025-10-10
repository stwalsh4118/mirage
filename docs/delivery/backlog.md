# Product Backlog

The backlog document contains all PBIs for the project, ordered by priority.

## PBIs

| ID | Actor | User Story | Status | Conditions of Satisfaction (CoS) |
| :-- | :---- | :--------- | :----- | :------------------------------- |
| <a id="1"></a>1 | Platform engineer | Deliver Mirage MVP enabling creation and destruction of single environments via UI backed by Railway. [View Details](./1/prd.md) | Agreed | Implement Environment Controller (MVP); Implement basic Railway integration; Provide dashboard and environment creation wizard; Support dev/prod templates; Persist state; Show near real-time status; Secure token storage |
| <a id="2"></a>2 | Platform engineer | Implement Monorepo Service Discovery to detect deployable services and dependencies. [View Details](./2/prd.md) | Agreed | Detect npm/pnpm/yarn workspaces and common monorepo tools; List services with path and dependencies; Allow selection of subset for deployment; Provide basic dependency graph view |
| <a id="3"></a>3 | Platform engineer | Implement Configuration Management Layer with inheritance and overrides. [View Details](./3/prd.md) | Agreed | Base→env-type→instance inheritance; Secret management integration; Service-specific overrides; Effective configuration preview |
| <a id="4"></a>4 | Backend engineer | Build GraphQL API Abstraction over Railway with batching, retries, and webhooks. [View Details](./4/prd.md) | Agreed | High-level lifecycle ops; Batched GQL calls; Exponential backoff; Webhook listener; Error classification and surfacing |
| <a id="5"></a>5 | Developer | Build Dashboard and Environment Creation Wizard with professional UX. [View Details](./5/prd.md) | Agreed | Card grid with type color-coding; 3-step wizard; Resource sliders/TTL; Quick actions (logs, open URL, destroy) |
| <a id="6"></a>6 | Developer | Add Real-time Status Updates for environments and services. [View Details](./6/prd.md) | Agreed | WebSocket or polling for status; Service-level health indicators; Clear error states and retry guidance |
| <a id="7"></a>7 | Platform admin | Add Persistence, Security, and RBAC foundations. [View Details](./7/prd.md) | Agreed | SQLite/Postgres persistence; Encrypted token storage; Basic roles; Audit logging scaffold; Resource attribution metadata (Mirage-created) and Railway association IDs |
| <a id="8"></a>8 | Tech lead | Add PR Integration for ephemeral envs lifecycle tied to PRs. [View Details](./8/prd.md) | Agreed | Webhook intake; Comment env URL on PR; Auto-destroy on close/merge; TTL enforcement |
| <a id="9"></a>9 | Developer | As a user, I want to browse my Railway projects and, in a project view, see its environments and services so I can navigate and understand what exists before creating anything in Mirage. [View Details](./9/prd.md) | Agreed | Dashboard lists Railway projects with counts; Project page shows environments and services (read-only); Uses Railway API with ~30s polling; No write/import actions |
| <a id="10"></a>10 | Developer | Build Environment Creation Wizard to provision a full Railway hierarchy (new or existing project) and create environments with services. [View Details](./10/prd.md) | Agreed | Create new Railway project or select existing; Create environment and core services; Apply templates (dev/prod); Configure env vars and TTL; Show progress and errors end-to-end |
| <a id="11"></a>11 | Developer | Enable service creation from pre-built Docker images instead of only from source repositories. [View Details](./11/prd.md) | Proposed | Support Docker Hub, GHCR, and custom registries; Configure image tag/digest; Set environment variables and ports; Validate image accessibility; Update wizard to support image-based deployment |
| <a id="12"></a>12 | Platform engineer | Automatically discover Dockerfiles in monorepo subdirectories and enable multi-service deployment from Docker images. [View Details](./12/prd.md) | Proposed | Scan monorepo for Dockerfiles; Detect service boundaries; Parse Dockerfile metadata; Present discovered services in UI; Support selective service deployment from images |
| <a id="13"></a>13 | Developer | Store comprehensive service build configuration and environment metadata to enable environment duplication and templating. [View Details](./13/prd.md) | Proposed | Persist Dockerfile paths, build args, and contexts; Store image registry config; Capture service dependencies; Enable metadata export/import; Support environment-as-template creation |
| <a id="14"></a>14 | Developer | View, search, and filter unified logs across all services in an environment. [View Details](./14/prd.md) | Proposed | Fetch logs from Railway API; Multi-service log aggregation; Real-time log streaming; Search and filtering; Time range selection; Export logs capability |
| <a id="15"></a>15 | Developer | Clone an existing environment with all its services, configuration, and metadata. [View Details](./15/prd.md) | Proposed | Deep copy environment metadata; Recreate all services with same config; Clone environment variables; Support cross-project cloning; Preserve service relationships; Update wizard with clone option |
| <a id="16"></a>16 | Platform engineer | Implement core Clerk authentication with JWT verification, user database, and resource ownership. [View Details](./16/prd.md) | Proposed | Frontend Clerk integration; Backend JWT middleware; User table with GORM; Clerk webhooks for user sync; Resource ownership (UserID foreign keys); Protected API routes; User profile endpoints |
| <a id="17"></a>17 | Platform engineer | Implement HashiCorp Vault for secure secret management of Railway API tokens and other user credentials. [View Details](./17/prd.md) | Proposed | Self-hosted Vault on Railway; Per-user secret storage; Railway/GitHub/Docker token management; Environment-specific secrets; Secret versioning and rotation; Comprehensive secret management UI |
| <a id="18"></a>18 | Platform admin | Add advanced authentication features including RBAC, audit logging, user preferences, and admin dashboard. [View Details](./18/prd.md) | Proposed | Role-based access control (admin/user); Audit logging of all actions; User preferences; Admin dashboard; User management UI; System statistics |
| <a id="19"></a>19 | Developer | Import and manage external Railway environments that existed before connecting to Mirage. [View Details](./19/prd.md) | Proposed | Visual distinction between managed/unmanaged environments; Environment-level import with validation; Mass import capability; Full CRUD operations post-import; Metadata capture (env vars, services, configs); Conflict detection and resolution |

## History

| Timestamp | PBI_ID | Event_Type | Details | User |
|----------:|:------:|:----------:|:--------|:-----|
| 20250905-120000 | 1 | create_pbi | Created PBI 1 (Mirage MVP: Core env provisioning) | ai-agent |
| 20250905-120100 | 2 | create_pbi | Created PBI 2 (Monorepo Service Discovery) | ai-agent |
| 20250905-120100 | 3 | create_pbi | Created PBI 3 (Configuration Management Layer) | ai-agent |
| 20250905-120100 | 4 | create_pbi | Created PBI 4 (GraphQL API Abstraction Layer) | ai-agent |
| 20250905-120100 | 5 | create_pbi | Created PBI 5 (Dashboard & Creation Wizard) | ai-agent |
| 20250905-120100 | 6 | create_pbi | Created PBI 6 (Real-time Updates) | ai-agent |
| 20250905-120100 | 7 | create_pbi | Created PBI 7 (Persistence & RBAC) | ai-agent |
| 20250905-120100 | 8 | create_pbi | Created PBI 8 (PR Integration) | ai-agent |
| 20250905-120300 | 1 | propose_for_backlog | Reframed PBIs as overall deliverables | ai-agent |
| 20250905-120400 | 2 | propose_for_backlog | PBI 2 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 3 | propose_for_backlog | PBI 3 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 4 | propose_for_backlog | PBI 4 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 5 | propose_for_backlog | PBI 5 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 6 | propose_for_backlog | PBI 6 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 7 | propose_for_backlog | PBI 7 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250905-120400 | 8 | propose_for_backlog | PBI 8 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250911-120000 | 9 | create_pbi | Created PBI 9 (Railway Project Browsing & Details) | ai-agent |
| 20250911-120010 | 7 | significant_update | Updated CoS to include resource attribution metadata and Railway association IDs | ai-agent |
| 20250911-120100 | 9 | propose_for_backlog | PBI 9 moved from Proposed to Agreed; detail doc created | ai-agent |
| 20250911-120500 | 10 | create_pbi | Created PBI 10 (Environment Creation Wizard & Full Railway Provisioning) | ai-agent |
| 20250911-120510 | 10 | propose_for_backlog | PBI 10 moved from Proposed to Agreed; detail doc created | sean |
| 20250930-143000 | 11 | create_pbi | Created PBI 11 (Docker Image Service Deployment) | ai-agent |
| 20250930-143000 | 12 | create_pbi | Created PBI 12 (Monorepo Dockerfile Discovery) | ai-agent |
| 20250930-143000 | 13 | create_pbi | Created PBI 13 (Service Build Configuration Management) | ai-agent |
| 20250930-143000 | 14 | create_pbi | Created PBI 14 (Service Logs Viewer) | ai-agent |
| 20250930-143000 | 15 | create_pbi | Created PBI 15 (Environment Cloning) | ai-agent |
| 20251007-000000 | 16 | create_pbi | Created PBI 16 (Clerk Authentication & User Management) | ai-agent |
| 20251007-000000 | 17 | create_pbi | Created PBI 17 (HashiCorp Vault Secret Management) | ai-agent |
| 20251007-010000 | 16 | significant_update | Simplified to core auth + resource ownership; moved advanced features to PBI 18 | ai-agent |
| 20251007-010000 | 18 | create_pbi | Created PBI 18 (Advanced Auth Features: RBAC, Audit, Admin) | ai-agent |
| 20251010-000000 | 19 | create_pbi | Created PBI 19 (Import External Railway Environments) | ai-agent |

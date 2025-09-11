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

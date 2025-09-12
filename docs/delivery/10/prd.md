# PBI-10: Environment Creation Wizard & Full Railway Provisioning

[View in Backlog](../backlog.md#user-content-10)

## Overview
Deliver an end-to-end Environment Creation Wizard that can provision a complete Railway hierarchy from scratch or use an existing Railway project. The wizard will create or select a Railway project, create an environment, and provision initial services with configuration.

## Problem Statement
Users need a guided, reliable flow to create environments. Some users will not have an existing Railway project and need Mirage to create it; others will select an existing project and add a new environment within it.

## User Stories
- As a user, I can create a new Railway project via the wizard.
- As a user, I can select an existing Railway project for environment creation.
- As a user, I can create an environment within the chosen project.
- As a user, I can provision initial services (from templates) and configure env vars, TTL, and resource presets.
- As a user, I can see progress and clear error messages during provisioning.

## Technical Approach
- Backend: use Railway GraphQL API to support createProject, createEnvironment, createService, and configuration updates.
- Frontend: multi-step wizard with Source, Configuration, and Strategy steps; server actions/hooks call backend APIs. Show progress steps and errors.
- Templates: support dev/prod templates to define default services and configuration. Apply after environment creation.
- Idempotency: ensure retries are safe; surface partial failures with guidance to resume or roll back.

## UX/UI Considerations
- Clear choice between "Create new Railway project" and "Use existing project" in Step 0.
- Stepper with back/next; summary review before submit; inline validation.
- Progress screen with per-step status (project, environment, services, configuration). Link to logs when available.

## Acceptance Criteria
- User can create a brand-new Railway project and environment end-to-end.
- User can select an existing Railway project and create an environment within it.
- Template application creates initial services and applies configuration.
- Errors are surfaced with actionable guidance and affordances to retry or clean up.
- Operation status is visible until completion; final environment appears on dashboard.

## Dependencies
- Railway API token and permissions.
- Template definitions for dev/prod.
- Backend lifecycle endpoints for project/env/service creation.

## Open Questions
- Exact Railway GraphQL mutation names and inputs to be confirmed at implementation time.
- Service bootstrap scope for MVP (e.g., placeholder vs real images).

## Related Tasks
[View in Backlog](../backlog.md#user-content-10)


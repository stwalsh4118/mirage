# PBI-9: Railway Project Browsing & Details

## Overview
Provide a read-only browsing experience of Railway projects within Mirage. Users can view projects on the dashboard and click into a project to see its environments and services.

## Problem Statement
Users need visibility into existing Railway resources before creating or managing anything in Mirage. Today, Mirage shows only its own environments.

## User Stories
- As a user, I want to see my Railway projects on the dashboard so I can pick where to work.
- As a user, I want to open a project page and see its environments and services.

## Technical Approach
- Use existing backend endpoints: `GET /railway/projects` and `GET /railway/project/:id?details=1`.
- Frontend: add projects grid to dashboard; add project details route with tabs (environments, services).
- Polling: 30s via React Query.
- No write/import actions.

## UX/UI Considerations
- Project cards show name, id, counts for services and environments.
- Project page header shows project name and id; tabs for Environments and Services.

## Acceptance Criteria
- Dashboard lists Railway projects.
- Clicking a project navigates to details page showing environments and services.
- Data is read-only and refreshes periodically.

## Dependencies
- Railway Public API token configured on backend.

## Open Questions
- Do we need environment status enrichment now? (out-of-scope for this PBI unless required.)

## Related Tasks
[View in Backlog](../backlog.md#9)

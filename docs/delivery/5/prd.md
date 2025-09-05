# PBI-5: Dashboard and Environment Creation Wizard

[View in Backlog](../backlog.md#user-content-5)

## Overview
Deliver a professional dashboard with environment grid and a three-step creation wizard to provision environments easily and consistently.

## Problem Statement
Users need a clear, fast UI to provision and manage environments without navigating low-level deployment details.

## User Stories
- As a developer, I can view all environments with status and quick actions.
- As a user, I can create an environment via a 3-step wizard.
- As a user, I can view logs and open service URLs from the UI.

## Technical Approach
- Next.js app with server-side data fetching and client-side real-time updates.
- Components: EnvironmentGrid, EnvironmentCard, CreationWizard (Source, Configuration, Strategy), ServiceMap.
- Integration with backend WebSocket/polling API.

## UX/UI Considerations
- Color-coded types: prod red, staging yellow, dev green, ephemeral blue.
- Resource sliders, TTL input, env var editor with inheritance preview.
- Quick actions: View logs, Open URL, Clone, Destroy.

## Acceptance Criteria
- Responsive dashboard showing status and actions.
- Wizard provisions environment end-to-end using backend APIs.
- Basic service map visualization and selection.

## Dependencies
- Backend APIs from PBI-1 and PBI-4.

## Open Questions
- Theming and dark mode; accessibility priorities.

## Related Tasks
- PBI-2 discovery feeds wizard; PBI-6 feeds real-time status.

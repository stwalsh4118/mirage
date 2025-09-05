# PBI-6: Real-time Status Updates

[View in Backlog](../backlog.md#user-content-6)

## Overview
Provide near real-time updates on environment and service status using WebSockets or polling, surfacing provisioning and health in the UI.

## Problem Statement
Polling manually or refreshing pages degrades user experience and hides transient issues. Users need timely feedback about provisioning and health.

## User Stories
- As a developer, I see environment status change within seconds without refreshing.
- As a user, I see service-level health indicators and errors.

## Technical Approach
- Backend: WebSocket server broadcasting environment/service updates; polling fallback with ETag/Last-Modified.
- Integration with Railway webhook signals and periodic reconciliation jobs.
- Client: subscribe to updates and update UI state.

## UX/UI Considerations
- Non-intrusive toasts for state changes and errors.
- Clear, color-coded health indicators at environment and service levels.

## Acceptance Criteria
- Average end-to-end latency under 5 seconds for status changes.
- Graceful fallback to polling if WS is unavailable.
- Error states include actionable messaging.

## Dependencies
- PBI-4 webhook events; PBI-1 environment state.

## Open Questions
- Backoff strategy for reconnects; offline handling.

## Related Tasks
- Feeds Dashboard (PBI-5) components.

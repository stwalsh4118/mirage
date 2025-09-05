# PBI-7: Persistence, Security, and RBAC Foundations

[View in Backlog](../backlog.md#user-content-7)

## Overview
Establish a secure persistence layer for environment state and Railway credentials, and implement basic role-based access control and audit logging.

## Problem Statement
Lack of persistence and access control risks data loss and unauthorized access to sensitive credentials.

## User Stories
- As an admin, I can store Railway API tokens securely.
- As an admin, I can assign basic roles and restrict destructive operations.
- As an engineer, I can audit who performed environment actions.

## Technical Approach
- Database: SQLite locally, Postgres on Railway; encrypted at rest where feasible.
- Secret storage with encryption and access policies; minimal secrets exposure to app memory.
- RBAC middleware and policy checks for environment operations.
- Audit log table capturing actor, action, target, timestamp, outcome.

## UX/UI Considerations
- Admin settings screens for tokens and roles (initially minimal).
- Clear error messages for permissions failures.

## Acceptance Criteria
- Tokens stored encrypted and never shown in plain text.
- Roles enforced for create/destroy operations.
- Audit records created for environment lifecycle actions.

## Dependencies
- Chosen crypto/key management approach; database migrations.

## Open Questions
- External identity provider integration (GitHub, OAuth) timeline.

## Related Tasks
- Supports all PBIs relying on secure state and authorization.

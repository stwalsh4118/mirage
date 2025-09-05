# PBI-8: PR Integration for Ephemeral Environments

[View in Backlog](../backlog.md#user-content-8)

## Overview
Integrate with the git provider to create ephemeral environments per pull request, comment URLs on the PR, and auto-destroy when merged or closed.

## Problem Statement
Manual setup of review environments slows code review and increases merge risk.

## User Stories
- As a developer, an environment is created automatically when I open a PR.
- As a reviewer, I see the environment URL in the PR comments.
- As a platform engineer, the environment is destroyed when the PR is merged or closed.

## Technical Approach
- Webhook ingestion from git provider (start with GitHub).
- Mapping from PR to environment specification; TTL defaults.
- Commenting API to post environment URL; cleanup on close/merge.

## UX/UI Considerations
- Link from PR to Mirage environment detail page.
- Status badges showing environment state in the PR.

## Acceptance Criteria
- Environments automatically created for new PRs; URL posted as comment.
- Environments are cleaned up when PR is closed/merged or TTL expires.
- Failure scenarios reported back to PR with guidance.

## Dependencies
- Git provider credentials; PBI-1/4/6 capabilities.

## Open Questions
- Multi-PR concurrency strategies and limit policies.

## Related Tasks
- Relies on discovery (PBI-2) and config (PBI-3) for accurate deploy specs.

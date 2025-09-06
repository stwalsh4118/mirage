# PBI-1: Mirage MVP – Core Environment Provisioning

[View in Backlog](../backlog.md#user-content-1)

## Overview
Mirage is an Environment-as-a-Service platform built atop Railway that provisions, manages, and observes application environments. It provides dynamic environment provisioning, monorepo-aware deployments, environment lifecycle management, and configuration inheritance/overrides. This PRD defines the platform vision and details the MVP scope for PBI-1.

## Problem Statement
Teams struggle to provision consistent, isolated environments quickly—especially in monorepos—leading to slow feedback cycles and configuration drift. Operating Railway directly for multi-service setups requires repetitive steps, fragmented visibility, and fragile manual processes.

## User Stories
- As a platform engineer, I can create, start, stop, and destroy an environment from a selected repo/branch.
- As a developer, I can view real-time status of environments and their services from a dashboard.
- As a developer, I can choose between environment templates (dev/prod) to match needs.
- As a platform admin, I can ensure tokens are stored securely and state is persisted.

## Technical Approach
- Backend (Go, Gin): HTTP API using Gin framework; Environment Controller managing lifecycle and state; persists in SQLite/Postgres; secure storage of Railway API token(s); webhook intake for deployment status.
- Railway Integration: GraphQL API client with retry/backoff; operations to create/delete environments and services; status polling fallback.
- Frontend (Next.js for routing/UI only): Dashboard grid and creation wizard; no Next.js API routes used; real-time updates via WebSocket or polling from the Go API.
- Auth: Clerk as the identity provider. Users are created in Clerk; a webhook from Clerk triggers user synchronization into Mirage's database.
- Data Model: Environment, Service, Template entities as described, persisted with creation timestamps and optional TTL.

## UX/UI Considerations
- Dashboard Grid: Card-based layout. For the full layout, interactions, and component breakdown, see [Dashboard Layout Spec](../../design/dashboard-layout.md).
- Environment Creation Wizard: 3 steps—Source Selection (repo/branch, service detection), Environment Configuration (template, TTL, resource sliders, env var editor with inheritance preview), Deployment Strategy (sequential/parallel, health checks, rollback triggers).
- Monorepo Service Map: Interactive dependency graph, selection/deselection, visual diff since last deployment, service health indicators.
- Environment Details: Service grid with logs, resource usage, public URLs, recent commits, and environment-wide actions (promote, clone, schedule destruction).

### Mirage Theme System

For the complete Mirage theme specification (palette, surfaces, motion, shadcn component guidelines, and implementation notes), see the centralized design document:

[Mirage Theme System](../../design/mirage-theme.md)

## Acceptance Criteria
- MVP: Ability to create and destroy a single environment via the UI using Railway; near real-time status visible in dashboard; support dev and prod templates; state persists across reloads; tokens stored securely.
- Non-Goals (MVP): Full monorepo detection UI, full configuration inheritance UI, PR integration, advanced cost tracking, collaboration comments, snapshots.

## Dependencies
- Railway GraphQL API credentials and network access.
- SQLite or Railway Postgres for persistence.
- Clerk account and webhook secret configuration.

## Open Questions
- Auth and RBAC scope in MVP (local admin vs external auth provider).
- Git provider support at launch (GitHub only or others).
- Minimal service detection behavior for MVP (manual selection vs heuristic).

## Related Tasks
- PBI-2 Monorepo Service Discovery
- PBI-3 Configuration Management Layer
- PBI-4 Railway GQL Abstraction Layer
- PBI-5 Dashboard & Creation Wizard
- PBI-6 Real-time Status Updates
- PBI-7 Persistence & RBAC Foundations
- PBI-8 PR Integration

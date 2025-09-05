# PBI-2: Monorepo Service Discovery

[View in Backlog](../backlog.md#user-content-2)

## Overview
Implement repository scanning to identify deployable services in monorepos, detect package manager/workspace conventions, and map inter-service dependencies to inform deployments.

## Problem Statement
Monorepos complicate environment provisioning because deployable units and their dependencies are not explicit. Manual identification is slow and error-prone, leading to broken or incomplete deployments.

## User Stories
- As a platform engineer, I can scan a repo to discover deployable services and their paths.
- As a developer, I can see inferred dependencies between services.
- As a user, I can select a subset of services to deploy for an environment.

## Technical Approach
- Parsers for package managers: npm/pnpm/yarn workspaces; support lerna, nx, turborepo configs.
- Language-agnostic service detection heuristics (Dockerfile, Procfile, package.json scripts, etc.).
- Dependency graph construction using workspace manifests and import graphs where feasible.
- Produce a normalized service manifest consumed by Environment Controller.

## UX/UI Considerations
- Repo scan step in wizard shows list of services with paths.
- Dependency graph visualization with selectable nodes.
- Filters (deploy all, changed only, selected subset).

## Acceptance Criteria
- Successfully detects services in example monorepos using common tools.
- Outputs service list with name, path, and basic metadata.
- Builds a dependency graph usable by deployments.
- Allows selecting subset of services for deployment.

## Dependencies
- Access to repository files; ability to read manifests.

## Open Questions
- Depth of static analysis for dependencies beyond workspace metadata.
- Support for polyglot repos (Go, Node, Python) in first iteration.

## Related Tasks
- PBI-1 MVP relies on a basic list for service selection in later iterations.
- PBI-5 Wizard consumes discovery results.

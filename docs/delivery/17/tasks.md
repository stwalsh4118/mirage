# Tasks for PBI 17: HashiCorp Vault Secret Management

This document lists all tasks associated with PBI 17.

**Parent PBI**: [PBI 17: HashiCorp Vault Secret Management](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 17-1 | [Research and document HashiCorp Vault HTTP API](./17-1.md) | Review | Research Vault HTTP API endpoints, KV v2 secrets engine, authentication methods, and create implementation guide |
| 17-2 | [Setup Docker Compose for local Vault development](./17-2.md) | Proposed | Add Vault service to docker-compose.yml for local development with dev mode |
| 17-3 | [Deploy Vault server on Railway with persistent storage](./17-3.md) | Proposed | Deploy Vault as Railway service with Raft storage and volume for data persistence |
| 17-4 | [Create Vault HTTP client wrapper and implement initialization](./17-4.md) | Proposed | Create vault package with HTTP client wrapper, initialization, and configuration loading |
| 17-5 | [Implement Vault authentication with token and AppRole](./17-5.md) | Proposed | Implement token auth for dev and AppRole auth for production with automatic token renewal |
| 17-6 | [Implement health checking and circuit breaker for Vault](./17-6.md) | Proposed | Add Vault connectivity health checks, circuit breaker pattern, and graceful degradation |
| 17-7 | [Implement secret caching layer with TTL](./17-7.md) | Proposed | Create in-memory cache for secrets with 5-minute TTL and invalidation on updates |
| 17-8 | [Define SecretStore interface and domain types](./17-8.md) | Proposed | Define SecretStore interface, Secret, SecretMetadata, and DockerCredentials types |
| 17-9 | [Implement Railway token management in SecretStore](./17-9.md) | Proposed | Implement store, get, delete, rotate, and validate methods for Railway tokens |
| 17-10 | [Implement GitHub token management in SecretStore](./17-10.md) | Proposed | Implement store, get, delete, and validate methods for GitHub PATs |
| 17-11 | [Implement Docker credentials management in SecretStore](./17-11.md) | Proposed | Implement store, get, list, and delete methods for Docker registry credentials |
| 17-12 | [Implement environment-specific secrets in SecretStore](./17-12.md) | Proposed | Implement store, get, get-all, delete, and bulk operations for environment secrets |
| 17-13 | [Implement generic secret management with versioning](./17-13.md) | Proposed | Implement generic secret CRUD, version management, rollback, and metadata operations |
| 17-14 | [Refactor Railway client to support per-user tokens](./17-14.md) | Proposed | Add GetUserRailwayClient helper function and factory pattern for user-specific clients |
| 17-15 | [Update controllers to use user-specific Railway clients](./17-15.md) | Proposed | Update Environment and Services controllers to fetch user tokens from Vault |
| 17-16 | [Implement Railway token API endpoints](./17-16.md) | Proposed | Create POST/GET/DELETE endpoints for Railway token storage, validation, and rotation |
| 17-17 | [Implement GitHub token API endpoints](./17-17.md) | Proposed | Create POST/GET/DELETE endpoints for GitHub PAT storage and validation |
| 17-18 | [Implement Docker credentials API endpoints](./17-18.md) | Proposed | Create POST/GET/DELETE endpoints for Docker registry credentials management |
| 17-19 | [Implement environment secrets API endpoints](./17-19.md) | Proposed | Create POST/GET/DELETE endpoints for environment-specific secrets with bulk operations |
| 17-20 | [Implement generic secrets and versioning API endpoints](./17-20.md) | Proposed | Create endpoints for generic secrets, version history, rollback, and audit log |
| 17-21 | [Create credentials settings page layout and structure](./17-21.md) | Proposed | Build main credentials page with tabbed interface for all secret types |
| 17-22 | [Implement Railway token UI tab with status and actions](./17-22.md) | Proposed | Create Railway token tab with status card, add/update/test/remove functionality |
| 17-23 | [Implement GitHub and Docker registry UI tabs](./17-23.md) | Proposed | Build GitHub token tab and Docker registries management tab with CRUD operations |
| 17-24 | [Implement custom secrets UI with version history](./17-24.md) | Proposed | Create custom secrets tab with searchable table, version history modal, and rollback |
| 17-25 | [Implement environment secrets UI in environment detail page](./17-25.md) | Proposed | Add secrets tab to environment page with table, bulk import/export functionality |
| 17-26 | [Implement feature flag and environment variable fallback](./17-26.md) | Proposed | Add VAULT_ENABLED flag, fallback to RAILWAY_API_TOKEN, and graceful degradation |
| 17-27 | [Add database fields for Vault migration tracking](./17-27.md) | Proposed | Add VaultEnabled and LastValidated fields to User model with migration |
| 17-28 | [Write unit tests for Vault client and SecretStore](./17-28.md) | Proposed | Test Vault client operations, secret store methods, caching, and error handling |
| 17-29 | [Write integration tests for secret management](./17-29.md) | Proposed | Integration tests with real Vault instance for all secret types and Railway client integration |
| 17-30 | [E2E CoS Test: Complete Vault secret management flow](./17-30.md) | Proposed | End-to-end test covering Vault deployment, secret storage, Railway integration, and UI workflows |



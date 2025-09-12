# Tasks for PBI 1: Mirage MVP – Core Environment Provisioning

This document lists all tasks associated with PBI 1.

**Parent PBI**: [PBI 1: Mirage MVP – Core Environment Provisioning](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 1-1 | [Initialize backend service and API scaffolding](./1-1.md) | Done | Create Go service project, HTTP API skeleton, configs |
| 1-2 | [Implement Railway GraphQL client (create/destroy env)](./1-2.md) | Done | Minimal GQL client with retry/backoff |
| 1-3 | [Design and migrate persistence schema](./1-3.md) | Done | DB schema for Environment, Service, Template |
| 1-4 | [Implement secure token storage and admin API](./1-4.md) | Blocked | Deferred for MVP; token via env; using Clerk with webhook sync |
| 1-5 | [Build Environment Controller MVP](./1-5.md) | Done | Lifecycle orchestration for create/destroy + state |
| 1-6 | [Implement status reconciliation polling](./1-6.md) | Done | Poll Railway, update state, emit events |
| 1-7 | [Initialize frontend (Next.js) and base layout](./1-7.md) | Done | Set up Next.js app and base UI shell |
| 1-8 | [Implement Dashboard grid and environment cards](./1-8.md) | InProgress | List environments, status, quick actions |
| 1-9 | [Moved: Implement Environment Creation Wizard (see 10-1)](../10/10-1.md) | Done | Task moved to PBI 10 |
| 1-10 | [E2E wiring and near-real-time polling](./1-10.md) | Proposed | Connect UI to backend; 5s polling updates |
| 1-11 | [Finalize Railway API integration](./1-11.md) | Done | Config, real GQL ops, list endpoint, status norms |

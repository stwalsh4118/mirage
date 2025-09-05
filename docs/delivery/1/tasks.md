# Tasks for PBI 1: Mirage MVP – Core Environment Provisioning

This document lists all tasks associated with PBI 1.

**Parent PBI**: [PBI 1: Mirage MVP – Core Environment Provisioning](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 1-1 | [Initialize backend service and API scaffolding](./1-1.md) | Done | Create Go service project, HTTP API skeleton, configs |
| 1-2 | [Implement Railway GraphQL client (create/destroy env)](./1-2.md) | Done | Minimal GQL client with retry/backoff |
| 1-3 | [Design and migrate persistence schema](./1-3.md) | InProgress | DB schema for Environment, Service, Template |
| 1-4 | [Implement secure token storage and admin API](./1-4.md) | Proposed | Encrypted Railway token storage and admin endpoint |
| 1-5 | [Build Environment Controller MVP](./1-5.md) | Proposed | Lifecycle orchestration for create/destroy + state |
| 1-6 | [Implement status reconciliation polling](./1-6.md) | Proposed | Poll Railway, update state, emit events |
| 1-7 | [Initialize frontend (Next.js) and base layout](./1-7.md) | Proposed | Set up Next.js app and base UI shell |
| 1-8 | [Implement Dashboard grid and environment cards](./1-8.md) | Proposed | List environments, status, quick actions |
| 1-9 | [Implement Environment Creation Wizard (dev/prod templates)](./1-9.md) | Proposed | 3-step wizard with minimal inputs |
| 1-10 | [E2E wiring and near-real-time polling](./1-10.md) | Proposed | Connect UI to backend; 5s polling updates |

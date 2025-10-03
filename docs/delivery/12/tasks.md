# Tasks for PBI 12: Monorepo Dockerfile Discovery

This document lists all tasks associated with PBI 12.

**Parent PBI**: [PBI 12: Monorepo Dockerfile Discovery](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 12-1 | [Research Railway API for service variables and Dockerfile path configuration](./12-1.md) | Done | Research and document how to set RAILWAY_DOCKERFILE_PATH and other service variables via Railway API |
| 12-2 | [Implement backend Dockerfile scanner](./12-2.md) | Done | Create Go service to recursively scan repositories for Dockerfiles and parse metadata |
| 12-3 | [Extend Railway service creation to support Dockerfile paths](./12-3.md) | Done | Extend CreateServiceInput and Railway client to support RAILWAY_DOCKERFILE_PATH variable |
| 12-4 | [Create Dockerfile discovery API endpoint](./12-4.md) | Done | Add REST endpoint for triggering Dockerfile discovery and retrieving results |
| 12-5 | [Build frontend Dockerfile discovery UI component](./12-5.md) | Done | Create React components to display discovered Dockerfiles in tree view with selection |
| 12-6 | [Integrate Dockerfile discovery into service creation wizard](./12-6.md) | Done | Add Dockerfile discovery step to wizard flow and connect to backend |
| 12-7 | [E2E CoS Test](./12-7.md) | Proposed | End-to-end testing of Dockerfile discovery and multi-service deployment |
| 12-8 | [Enhance wizard to support per-service environment variables](./12-8.md) | Done | Add support for configuring environment variables per service with global variable inheritance |


# Tasks for PBI 13: Service Build Configuration and Environment Metadata Management

This document lists all tasks associated with PBI 13.

**Parent PBI**: [PBI 13: Service Build Configuration and Environment Metadata Management](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 13-1 | [Extend Environment and Service models with build configuration fields](./13-1.md) | Done | Add missing fields to existing models and remove unused Template model |
| 13-2 | [Create EnvironmentMetadata model for wizard inputs and provision outputs](./13-2.md) | Done | New model to store complete wizard state and enable cloning/templates |
| 13-3 | [Implement database persistence in provision endpoints](./13-3.md) | Done | Add database writes to ProvisionProject, ProvisionEnvironment, and ProvisionServices |
| 13-4 | [Remove unused CRUD endpoints and clean up dead code](./13-4.md) | Done | Delete old Environment CRUD routes that are not used by wizard |
| 13-5 | [Implement metadata retrieval API endpoints](./13-5.md) | Done | Add GET endpoints for retrieving persisted environment and service data |
| 13-6 | [Build frontend metadata display components](./13-6.md) | Proposed | Create UI components to display environment metadata and service configurations |
| 13-7 | [E2E CoS Test](./13-7.md) | Proposed | End-to-end testing of persistence, metadata capture, and retrieval |


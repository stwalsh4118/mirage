# Tasks for PBI 15: Environment Cloning

This document lists all tasks associated with PBI 15.

**Parent PBI**: [PBI 15: Environment Cloning](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 15-1 | [Research Railway API for environment variables and create reference guide](./15-1.md) | Review | Document Railway's GraphQL API for fetching and setting environment variables |
| 15-2 | [Implement environment variable fetching from Railway in Go client](./15-2.md) | Proposed | Add GetEnvironmentVariables method to Railway client |
| 15-3 | [Implement environment variable upsert to Railway in Go client](./15-3.md) | Proposed | Add UpsertEnvironmentVariables method to Railway client |
| 15-4 | [Add snapshot endpoint to fetch environment data](./15-4.md) | Proposed | Implement GET /api/environments/:id/snapshot to return all env data for cloning |
| 15-5 | [Add "Clone from" mode to creation wizard](./15-5.md) | Proposed | Add source selection option and environment picker to existing wizard |
| 15-6 | [Pre-populate wizard from environment snapshot](./15-6.md) | Proposed | Fetch snapshot and pre-fill all wizard steps with source environment data |
| 15-7 | [E2E CoS Test](./15-7.md) | Proposed | End-to-end testing of clone workflow using existing provision endpoints |

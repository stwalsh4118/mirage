# Tasks for PBI 14: Service Logs Viewer

This document lists all tasks associated with PBI 14.

**Parent PBI**: [PBI 14: Service Logs Viewer](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 14-0 | [Refactor GraphQL queries to separate .graphql files with embed](./14-0.md) | Done | Extract inline query strings to .graphql files using Go embed pattern |
| 14-1 | [Research Railway GraphQL log API and create reference guide](./14-1.md) | Done | Document Railway's log subscription and query APIs, WebSocket protocol, and authentication |
| 14-2 | [Implement Railway log subscription and fetching in Go client](./14-2.md) | Done | Add real-time log subscription using hasura/go-graphql-client and historical log queries |
| 14-3 | [Add log processing utilities for parsing and formatting](./14-3.md) | Proposed | Implement log parsing for structured data, severity detection, ANSI code handling |
| 14-4 | [Create HTTP API endpoints for log retrieval and export](./14-4.md) | Proposed | Add GET /api/environments/:id/logs and /api/logs/export endpoints with query parameters |
| 14-5 | [Implement WebSocket endpoint for real-time log streaming](./14-5.md) | Proposed | Relay logs from Railway subscription to frontend clients via WebSocket |
| 14-6 | [Create LogViewer React component with virtual scrolling](./14-6.md) | Proposed | Build main log display component with virtualization, syntax highlighting, and line numbers |
| 14-7 | [Implement log filtering and search controls UI](./14-7.md) | Proposed | Build filter controls for service, time range, search (regex), and log level |
| 14-8 | [Add WebSocket client and real-time streaming logic](./14-8.md) | Proposed | Implement WebSocket connection, auto-reconnect, and buffer management in frontend |
| 14-9 | [Implement log export functionality (JSON, CSV, TXT)](./14-9.md) | Proposed | Add export button with format selection and file download capability |
| 14-10 | [Integrate LogViewer into environment detail page](./14-10.md) | Proposed | Add logs tab/section to project-detail.tsx and wire up API calls |
| 14-11 | [E2E CoS Test](./14-11.md) | Proposed | End-to-end testing of log viewing, filtering, search, streaming, and export |


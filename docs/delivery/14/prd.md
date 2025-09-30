# PBI-14: Service Logs Viewer

[View in Backlog](../backlog.md#user-content-14)

## Overview
Provide a unified log viewing experience that aggregates, displays, and enables searching across logs from all services within an environment. Users should be able to view real-time logs, search historical logs, filter by service, and export logs for analysis.

## Problem Statement
Currently, users must navigate to Railway directly or check individual service logs separately to debug issues. This fragmented experience makes it difficult to:
- Correlate events across multiple services
- Debug distributed system issues
- Monitor application behavior in real-time
- Search for specific errors or patterns across services
- Analyze log data for troubleshooting

Providing unified log viewing within Mirage improves the developer experience and reduces context switching.

## User Stories
- As a developer, I want to view logs from all services in one place so I can debug multi-service issues
- As a developer, I want to filter logs by service, severity, or time range to find relevant information quickly
- As a developer, I want to search logs for specific text or patterns
- As a developer, I want to see real-time log streaming to monitor live application behavior
- As a developer, I want to export logs for offline analysis or sharing with my team

## Technical Approach

### Backend (Go)

#### Log Fetching
- **Railway API Integration**: Use Railway's GraphQL API to fetch logs
  - `deploymentLogs` query for historical logs
  - Implement pagination for large log volumes
  - Support filtering by deployment ID, time range
- **Log Aggregation**: Combine logs from multiple services
  - Merge by timestamp for chronological view
  - Preserve service identity in log entries
  - Handle time zone conversions

#### API Endpoints
- `GET /api/environments/:id/logs` - Get aggregated logs for environment
  - Query params: `services`, `from`, `to`, `limit`, `search`
- `GET /api/services/:id/logs` - Get logs for specific service
- `GET /api/logs/export` - Export logs in various formats (JSON, CSV, plain text)
- `WS /api/environments/:id/logs/stream` - WebSocket for real-time log streaming

#### Log Processing
- Parse log lines for structured data (timestamp, level, message)
- Detect common log formats (JSON, logfmt, plain text)
- Extract severity levels (INFO, WARN, ERROR, DEBUG)
- Support ANSI color code stripping/rendering

### Frontend (Next.js/React)

#### Logs Viewer Component
- **Header Section**:
  - Service filter (multi-select dropdown)
  - Time range picker (Last 15 min, 1 hour, 6 hours, 24 hours, custom)
  - Search input with regex support toggle
  - Log level filter (INFO, WARN, ERROR, DEBUG)
  - Export button (dropdown: JSON, CSV, TXT)
  - Auto-scroll toggle
  - Pause/Resume streaming button

- **Log Display Area**:
  - Virtual scrolling for performance with large log volumes
  - Monospace font with syntax highlighting
  - Color-coded by severity level
  - Service name badge on each line
  - Timestamp formatting (absolute or relative)
  - Expandable for long lines
  - Line numbers

- **Features**:
  - Real-time streaming via WebSocket
  - Infinite scroll for historical logs (pagination)
  - Highlight search matches
  - Copy log lines to clipboard
  - Jump to top/bottom buttons
  - Context menu: filter by service, highlight errors

#### Performance Considerations
- Implement virtual scrolling (react-window or similar)
- Limit buffer size (e.g., last 5000 lines in memory)
- Debounce search input
- Lazy load historical logs on scroll
- WebSocket reconnection with exponential backoff

## UX/UI Considerations

### Logs Viewer Layout
- Full-width panel below environment details
- Option to open in modal/fullscreen mode
- Split view: service list on left, logs on right
- Responsive design (mobile-friendly scrolling)

### Visual Design
- **Log Lines**:
  - Gray background for alternating lines (zebra striping)
  - Red highlight for ERROR level
  - Yellow highlight for WARN level
  - Blue/white for INFO/DEBUG
  - Service badge with color per service
  - Timestamp in muted color

- **Search Highlighting**: Yellow background on matching text
- **Empty State**: "No logs available" with helpful message
- **Loading State**: Skeleton loader or spinner
- **Error State**: Friendly error message with retry button

### Filtering UX
- Multi-select service filter with "Select All" option
- Time range quick picks + custom date/time picker
- Search supports plain text and regex (toggle)
- Filters persist across page navigation (URL params)

### Export Options
- **JSON**: Structured format with metadata
- **CSV**: Timestamp, Service, Level, Message
- **TXT**: Plain text format
- Filename: `{environment-name}-logs-{timestamp}.{ext}`

## Acceptance Criteria
1. Logs page accessible from environment detail view
2. Logs from all services in environment are aggregated
3. Logs displayed in chronological order
4. Real-time log streaming works via WebSocket
5. Users can filter logs by service (multi-select)
6. Users can filter logs by time range
7. Search functionality finds text across all displayed logs
8. Regex search mode is supported
9. Log severity levels are detected and color-coded
10. Users can export logs in JSON, CSV, and TXT formats
11. Virtual scrolling performs well with 10,000+ log lines
12. Auto-scroll can be toggled on/off
13. WebSocket reconnects automatically on connection loss
14. Empty state displayed when no logs available
15. Error handling for API failures with retry option

## Dependencies
- Railway API access for log fetching
- WebSocket infrastructure for real-time streaming
- Understanding of Railway's log data format and limits

## Open Questions
1. What is Railway's rate limit for log API calls?
2. How far back does Railway retain logs?
3. Should we implement local log caching/persistence?
4. Do we need log analytics (error rates, patterns)?
5. Should we support log forwarding to external services (e.g., Datadog, Splunk)?
6. How do we handle very high-throughput log streams (1000s/sec)?
7. Should we support log annotations or bookmarks?

## Related Tasks
Tasks will be created once this PBI moves to "Agreed" status.


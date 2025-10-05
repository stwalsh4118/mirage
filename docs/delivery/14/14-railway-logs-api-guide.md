# Railway GraphQL Logs API Reference Guide

**Created:** 2025-10-05  
**Task:** 14-1  
**Purpose:** Document Railway's GraphQL log API for implementation reference

## Overview

Railway provides GraphQL APIs for accessing deployment and environment logs through both queries (historical logs) and subscriptions (real-time streaming). This guide documents the API structure, parameters, and implementation patterns for Mirage's log viewing feature.

**Note:** This guide is based on Railway's known API patterns and should be validated against their actual GraphQL schema during implementation. Railway's public API documentation is limited, so some details are inferred from common GraphQL patterns and may require adjustment.

## API Endpoints

```
HTTP GraphQL Endpoint: https://backboard.railway.app/graphql/v2
WebSocket Endpoint:    wss://backboard.railway.app/graphql/internal
```

**Important Notes:**
- The WebSocket endpoint uses `/graphql/internal` path
- This endpoint was observed in Railway's sidecar tooling designed to run as a Railway service
- ✅ **CONFIRMED:** This endpoint IS accessible from external clients (tested 2025-10-05)
- **Subprotocol Required:** Must specify `graphql-transport-ws` subprotocol for connection
- Connection without proper subprotocol returns: `Disconnected (code: 4406, reason: "Subprotocol not acceptable")`

## Authentication

All requests (both HTTP and WebSocket) require Bearer token authentication:

```
Authorization: Bearer <RAILWAY_API_TOKEN>
```

**Note:** Token authentication may differ between internal and external endpoints.

## Log Subscriptions

### environmentLogs Subscription

The primary subscription for real-time log streaming across all services in an environment.

**GraphQL Schema:**

```graphql
subscription StreamEnvironmentLogs(
  $environmentId: String!
  $filter: String
  $beforeLimit: Int!
  $beforeDate: String
  $anchorDate: String
  $afterDate: String
  $afterLimit: Int
) {
  environmentLogs(
    environmentId: $environmentId
    filter: $filter
    beforeDate: $beforeDate
    anchorDate: $anchorDate
    afterDate: $afterDate
    beforeLimit: $beforeLimit
    afterLimit: $afterLimit
  ) {
    timestamp
    message
    severity
    tags {
      projectId
      environmentId
      serviceId
      deploymentId
      deploymentInstanceId
      snapshotId
    }
    attributes {
      key
      value
    }
  }
}
```

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `environmentId` | String | Yes | The ID of the environment to stream logs from |
| `filter` | String | No | Text filter to apply to log messages |
| `beforeLimit` | Int | Yes | Maximum number of historical logs to fetch before anchor |
| `beforeDate` | String | No | ISO 8601 timestamp to fetch logs before this time |
| `anchorDate` | String | No | ISO 8601 timestamp to use as anchor point for pagination |
| `afterDate` | String | No | ISO 8601 timestamp to fetch logs after this time |
| `afterLimit` | Int | No | Maximum number of logs to fetch after anchor |

**Response Structure:**

```typescript
{
  timestamp: string;        // ISO 8601 timestamp
  message: string;          // Raw log message
  severity: string;         // Log level: "INFO" | "WARN" | "ERROR" | "DEBUG" | null
  tags: {
    projectId: string;
    environmentId: string;
    serviceId: string;
    deploymentId: string;
    deploymentInstanceId: string;
    snapshotId: string;
  };
  attributes: Array<{
    key: string;
    value: string;
  }>;
}
```

**Example Variables:**

```json
{
  "environmentId": "550e8400-e29b-41d4-a716-446655440000",
  "filter": "",
  "beforeLimit": 100,
  "afterLimit": 0
}
```

### deploymentLogs Subscription

For streaming logs from a specific deployment.

**GraphQL Schema:**

```graphql
subscription StreamDeploymentLogs(
  $deploymentId: String!
  $filter: String
  $limit: Int
) {
  deploymentLogs(
    deploymentId: $deploymentId
    filter: $filter
    limit: $limit
  ) {
    timestamp
    message
    severity
    tags {
      deploymentId
      serviceId
    }
  }
}
```

## Historical Log Queries

### deploymentLogs Query

For fetching historical logs from a specific deployment.

**GraphQL Schema:**

```graphql
query GetDeploymentLogs(
  $deploymentId: String!
  $limit: Int
  $before: String
  $after: String
) {
  deploymentLogs(
    deploymentId: $deploymentId
    limit: $limit
    before: $before
    after: $after
  ) {
    logs {
      timestamp
      message
      severity
    }
    hasMore: Boolean
    cursor: String
  }
}
```

**Pagination:** Uses cursor-based pagination with `before`/`after` cursors and `limit`.

### buildLogs Query

For fetching build logs from a deployment.

**GraphQL Schema:**

```graphql
query GetBuildLogs(
  $deploymentId: String!
) {
  buildLogs(deploymentId: $deploymentId) {
    timestamp
    message
  }
}
```

**Note:** Build logs typically don't support filtering or pagination as they're finite and relatively small.

## WebSocket Protocol

Railway uses the **graphql-transport-ws** protocol (not the older graphql-ws protocol) for subscriptions.

### Connection Lifecycle

1. **Connection Initialization:**
   ```json
   {
     "type": "connection_init",
     "payload": {
       "headers": {
         "Authorization": "Bearer <RAILWAY_API_TOKEN>"
       }
     }
   }
   ```

2. **Server Acknowledgment:**
   ```json
   {
     "type": "connection_ack"
   }
   ```

3. **Subscribe to Logs:**
   ```json
   {
     "id": "1",
     "type": "subscribe",
     "payload": {
       "query": "subscription StreamEnvironmentLogs($environmentId: String!) { ... }",
       "variables": {
         "environmentId": "..."
       }
     }
   }
   ```

4. **Receive Log Messages:**
   ```json
   {
     "id": "1",
     "type": "next",
     "payload": {
       "data": {
         "environmentLogs": {
           "timestamp": "2025-10-05T12:34:56Z",
           "message": "Server started on port 3000",
           "severity": "INFO",
           ...
         }
       }
     }
   }
   ```

5. **Heartbeat (Keep-Alive):**
   Railway expects periodic ping/pong messages:
   - Client sends: `{"type": "ping"}`
   - Server responds: `{"type": "pong"}`
   - Recommended interval: Every 30 seconds

6. **Unsubscribe:**
   ```json
   {
     "id": "1",
     "type": "complete"
   }
   ```

7. **Close Connection:**
   ```json
   {
     "type": "connection_terminate"
   }
   ```

### Reconnection Strategy

- Implement exponential backoff for reconnections
- Start with 1 second delay, max out at 30 seconds
- Reset backoff on successful connection
- Maintain subscription state across reconnections

## Go Implementation with github.com/coder/websocket

### Package Installation

Based on Railway's actual implementation, use the `github.com/coder/websocket` package directly:

```bash
go get github.com/coder/websocket@latest
# Note: google/uuid is already in the project for generating subscription IDs
```

Railway's implementation uses `github.com/coder/websocket` (formerly `nhooyr.io/websocket`, now maintained by Coder) directly rather than a higher-level GraphQL client, providing more control over the WebSocket connection and message handling.

**About this library:** Minimal and idiomatic WebSocket library for Go with first-class `context.Context` support, zero dependencies, and full RFC compliance. [Source](https://github.com/coder/websocket)

### Railway's Actual Implementation Pattern

Based on Railway's sidecar implementation, here's the proven pattern:

```go
package railway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/coder/websocket"
)

const (
	subscriptionTypeSubscribe = "subscribe"
	subscriptionTypeComplete  = "complete"
)

var (
	connectionInit = []byte(`{"type":"connection_init"}`)
	connectionAck  = []byte(`{"type":"connection_ack"}`)
)

// EnvironmentLogsSubscriptionVariables matches Railway's expected structure
type EnvironmentLogsSubscriptionVariables struct {
	EnvironmentID string `json:"environmentId"`
	Filter        string `json:"filter,omitempty"`
	BeforeDate    string `json:"beforeDate,omitempty"`  // ISO 8601 / RFC3339Nano
	BeforeLimit   int    `json:"beforeLimit"`
	AnchorDate    string `json:"anchorDate,omitempty"`
	AfterDate     string `json:"afterDate,omitempty"`
	AfterLimit    int    `json:"afterLimit,omitempty"`
}

// EnvironmentLogsSubscriptionPayload is the GraphQL subscription payload
type EnvironmentLogsSubscriptionPayload struct {
	Query     string                                 `json:"query"`
	Variables *EnvironmentLogsSubscriptionVariables `json:"variables"`
}

// EnvironmentLog represents a single log entry from Railway
type EnvironmentLog struct {
	Timestamp  string                 `json:"timestamp"`
	Message    string                 `json:"message"`
	Severity   string                 `json:"severity"`
	Tags       map[string]string      `json:"tags"`
	Attributes []map[string]string    `json:"attributes"`
}

// CreateWebSocketSubscription establishes a WebSocket connection following Railway's pattern
func (c *Client) CreateWebSocketSubscription(ctx context.Context, payload any) (*websocket.Conn, error) {
	// Generate unique subscription ID
	subID := uuid.Must(uuid.NewV4())
	
	subPayload := map[string]any{
		"id":      subID,
		"type":    subscriptionTypeSubscribe,
		"payload": payload,
	}

	payloadBytes, err := json.Marshal(&subPayload)
	if err != nil {
		return nil, err
	}

	// Configure WebSocket dial options
	opts := &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{"Bearer " + c.token},
			"Content-Type":  []string{"application/json"},
		},
		Subprotocols: []string{"graphql-transport-ws"},
	}

	// Establish connection with timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctxTimeout, c.wsEndpoint, opts)
	if err != nil {
		return nil, err
	}

	// Remove read limit (Railway sends large log batches)
	conn.SetReadLimit(-1)

	// Step 1: Send connection_init
	if err := conn.Write(ctx, websocket.MessageText, connectionInit); err != nil {
		conn.Close(websocket.StatusInternalError, "failed to send init")
		return nil, err
	}

	// Step 2: Wait for connection_ack
	_, ackMessage, err := conn.Read(ctx)
	if err != nil {
		conn.Close(websocket.StatusInternalError, "failed to read ack")
		return nil, err
	}

	if !bytes.Equal(ackMessage, connectionAck) {
		conn.Close(websocket.StatusProtocolError, "invalid ack")
		return nil, errors.New("did not receive connection ack from server")
	}

	// Step 3: Send subscription payload
	if err := conn.Write(ctx, websocket.MessageText, payloadBytes); err != nil {
		conn.Close(websocket.StatusInternalError, "failed to send subscription")
		return nil, err
	}

	return conn, nil
}

// SubscribeToEnvironmentLogs creates an environment logs subscription using Railway's pattern
func (c *Client) SubscribeToEnvironmentLogs(
	ctx context.Context,
	environmentID string,
	serviceFilter string,
) (*websocket.Conn, error) {
	// Build subscription query (same as documented earlier)
	query := `subscription StreamEnvironmentLogs(
		$environmentId: String!
		$filter: String
		$beforeLimit: Int!
		$beforeDate: String
		$anchorDate: String
		$afterDate: String
		$afterLimit: Int
	) {
		environmentLogs(
			environmentId: $environmentId
			filter: $filter
			beforeDate: $beforeDate
			anchorDate: $anchorDate
			afterDate: $afterDate
			beforeLimit: $beforeLimit
			afterLimit: $afterLimit
		) {
			timestamp
			message
			severity
			tags
			attributes
		}
	}`

	payload := &EnvironmentLogsSubscriptionPayload{
		Query: query,
		Variables: &EnvironmentLogsSubscriptionVariables{
			EnvironmentID: environmentID,
			Filter:        serviceFilter,
			// Get last 5 minutes of logs for seamless subscription resuming
			BeforeDate:  time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339Nano),
			BeforeLimit: 500,
		},
	}

	return c.CreateWebSocketSubscription(ctx, payload)
}

// ReadLogMessage reads and parses a log message from the WebSocket
func ReadLogMessage(ctx context.Context, conn *websocket.Conn) (*EnvironmentLog, error) {
	_, data, err := conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	var msg struct {
		Type    string `json:"type"`
		ID      string `json:"id"`
		Payload struct {
			Data struct {
				EnvironmentLogs EnvironmentLog `json:"environmentLogs"`
			} `json:"data"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if msg.Type != "next" {
		return nil, nil // Skip non-data messages
	}

	return &msg.Payload.Data.EnvironmentLogs, nil
}
```

### Usage Example (Railway Pattern)

```go
func streamLogs(ctx context.Context, client *railway.Client, envID string) error {
	// Create subscription
	conn, err := client.SubscribeToEnvironmentLogs(ctx, envID, "")
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")

	// Read logs continuously
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log, err := railway.ReadLogMessage(ctx, conn)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			
			if log != nil {
				fmt.Printf("[%s] %s\n", log.Timestamp, log.Message)
			}
		}
	}
}
```

### Hybrid Approach: machinebox + github.com/coder/websocket

Based on our current implementation using `machinebox/graphql` for queries and Railway's proven WebSocket pattern:

```go
// client.go - Updated structure

package railway

import (
	"context"
	"net/http"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/machinebox/graphql"
	"github.com/coder/websocket"
)

const (
	DefaultEndpoint   = "https://backboard.railway.app/graphql/v2"
	DefaultWSEndpoint = "wss://backboard.railway.app/graphql/internal"
)

// Client wraps both HTTP (machinebox) and WebSocket (coder/websocket) clients
type Client struct {
	endpoint   string
	wsEndpoint string
	token      string
	httpc      *http.Client
}

func NewClient(endpoint, token string, httpc *http.Client) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	
	if httpc == nil {
		httpc = &http.Client{Timeout: 30 * time.Second}
	}
	
	return &Client{
		endpoint:   endpoint,
		wsEndpoint: DefaultWSEndpoint,
		token:      token,
		httpc:      httpc,
	}
}

// execute - existing method for queries/mutations
func (c *Client) execute(ctx context.Context, gql string, vars map[string]any, out any) error {
	client := graphql.NewClient(c.endpoint, graphql.WithHTTPClient(c.httpc))
	req := graphql.NewRequest(gql)
	for k, v := range vars {
		req.Var(k, v)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	operation := func() error {
		return client.Run(ctx, req, out)
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 300 * time.Millisecond
	bo.MaxInterval = 2 * time.Second
	bo.MaxElapsedTime = 10 * time.Second

	if err := backoff.Retry(func() error { return operation() }, backoff.WithContext(bo, ctx)); err != nil {
		return err
	}
	return nil
}

// WebSocket subscription methods (CreateWebSocketSubscription, SubscribeToEnvironmentLogs)
// are shown in the previous section
```

## Rate Limits and Constraints

**Note:** Railway's exact rate limits are not publicly documented. Based on typical GraphQL API patterns:

- **Query Rate Limit:** Likely ~100-1000 requests per minute per token
- **Subscription Connections:** Probably limited to 5-10 concurrent subscriptions per token
- **WebSocket Message Rate:** Logs stream at production rate (can be 100s-1000s per second)

**Recommendations:**
- Implement client-side rate limiting for queries
- Use single subscription per environment, multiplex to multiple clients
- Implement buffering and backpressure handling for high-throughput streams
- Add circuit breakers for repeated failures

## Log Retention

**Estimated Retention:** Railway likely retains logs for:
- **Active deployments:** 7-14 days
- **Build logs:** 30 days
- **Historical logs:** Limited to recent deployments

**Recommendation:** Implement local caching or export for long-term retention.

## Log Format and Structure

### Message Format

Logs are returned as plain text strings that may contain:
- ANSI color codes (e.g., `\x1b[32m` for green)
- Structured JSON (if application logs JSON)
- Timestamps (if application includes them)
- Multi-line messages

### Severity Detection

The `severity` field may be null. When null, severity can be inferred from message content:

```go
func detectSeverity(message string) string {
	lowerMsg := strings.ToLower(message)
	switch {
	case strings.Contains(lowerMsg, "error"), strings.Contains(lowerMsg, "fatal"):
		return "ERROR"
	case strings.Contains(lowerMsg, "warn"):
		return "WARN"
	case strings.Contains(lowerMsg, "debug"):
		return "DEBUG"
	default:
		return "INFO"
	}
}
```

### ANSI Code Handling

Strip ANSI codes for storage, preserve for display:

```go
import "regexp"

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
```

## Pagination Strategy

Railway uses **anchor-based pagination** for subscriptions:

1. **Initial Load:** Set `beforeLimit` to fetch historical logs (e.g., 100)
2. **Anchor Point:** Use `anchorDate` to mark a specific timestamp
3. **Load Earlier:** Set `beforeDate` and `beforeLimit` to fetch older logs
4. **Load Later:** Set `afterDate` and `afterLimit` to fetch newer logs

This allows bi-directional scrolling through log history while maintaining real-time streaming.

## Error Handling

Common errors and handling strategies:

| Error | Cause | Solution |
|-------|-------|----------|
| `401 Unauthorized` | Invalid or expired token | Refresh token, re-authenticate |
| `404 Not Found` | Environment/deployment doesn't exist | Validate IDs before subscribing |
| `429 Too Many Requests` | Rate limit exceeded | Implement exponential backoff |
| WebSocket disconnect | Network issue, timeout | Auto-reconnect with backoff |
| `500 Internal Server Error` | Railway service issue | Retry with exponential backoff |

## File Structure for GraphQL Operations

Following the pattern from task 14-0:

```
api/internal/railway/queries/
├── subscriptions/
│   ├── environment-logs.graphql
│   └── deployment-logs.graphql
├── queries/
│   ├── deployment-logs-history.graphql
│   └── build-logs.graphql
└── README.md
```

Example file: `api/internal/railway/queries/subscriptions/environment-logs.graphql`

```graphql
subscription StreamEnvironmentLogs(
  $environmentId: String!
  $filter: String
  $beforeLimit: Int!
  $beforeDate: String
  $anchorDate: String
  $afterDate: String
  $afterLimit: Int
) {
  environmentLogs(
    environmentId: $environmentId
    filter: $filter
    beforeDate: $beforeDate
    anchorDate: $anchorDate
    afterDate: $afterDate
    beforeLimit: $beforeLimit
    afterLimit: $afterLimit
  ) {
    timestamp
    message
    severity
    tags {
      projectId
      environmentId
      serviceId
      deploymentId
      deploymentInstanceId
      snapshotId
    }
    attributes {
      key
      value
    }
  }
}
```

Note: Unlike queries, subscriptions with hasura client don't use embedded strings the same way. The subscription structure is defined via Go struct tags, but having the .graphql file serves as documentation and reference.

## Testing and Validation

### Validate Schema

Use GraphQL introspection to verify the schema:

```graphql
query IntrospectionQuery {
  __type(name: "Subscription") {
    name
    fields {
      name
      args {
        name
        type {
          name
          kind
        }
      }
    }
  }
}
```

### Test Subscription Manually

Use a GraphQL client like Insomnia or Postman that supports subscriptions to test before implementing.

## External Access Considerations

### Testing External Connectivity

✅ **CONFIRMED WORKING (2025-10-05):** The WebSocket endpoint is accessible externally.

Test connection with proper subprotocol:

```bash
# Test WebSocket connection from external client
wscat -c wss://backboard.railway.app/graphql/internal \
  -H "Authorization: Bearer $RAILWAY_TOKEN" \
  -s graphql-transport-ws
```

**Expected behavior:**
1. Connection establishes successfully
2. Interactive terminal appears
3. Server waits for `connection_init` message
4. Without proper handshake, connection times out or shows "invalid message received"

**To complete handshake manually in wscat:**
```json
{"type":"connection_init","payload":{"headers":{"Authorization":"Bearer YOUR_TOKEN"}}}
```

Expected response:
```json
{"type":"connection_ack"}
```

### Connection Requirements Summary

✅ **Confirmed working configuration:**
- **Endpoint:** `wss://backboard.railway.app/graphql/internal`
- **Subprotocol:** `graphql-transport-ws` (REQUIRED)
- **Authentication:** Bearer token in headers
- **Protocol:** Full graphql-transport-ws handshake required

**Common Issues:**
- Without subprotocol: `code: 4406, reason: "Subprotocol not acceptable"`
- Without handshake: Connection timeout or "invalid message received"
- Invalid token: `401 Unauthorized`

## References

- **Railway Platform:** https://railway.app
- **Railway Docs:** https://docs.railway.app
- **Railway API Domain:** https://backboard.railway.app
- **Railway Sidecar (reference implementation):** Source of actual implementation patterns
- **github.com/coder/websocket:** https://github.com/coder/websocket (minimal and idiomatic WebSocket library for Go)
- **graphql-transport-ws Protocol:** https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md
- **machinebox/graphql:** https://github.com/machinebox/graphql
- **wscat (WebSocket testing):** https://github.com/websockets/wscat
- **google/uuid:** https://github.com/google/uuid

## Next Steps

✅ **External access confirmed** - Proceed with Task 14-2 implementation as planned.

**Key implementation notes (based on Railway's actual pattern):**
- Use `github.com/coder/websocket` directly (Railway's proven approach)
- Specify `graphql-transport-ws` subprotocol in DialOptions
- Follow Railway's 3-step handshake: init → ack → subscribe
- Set `conn.SetReadLimit(-1)` for large log batches
- Use `BeforeDate` (5 minutes back) and `BeforeLimit: 500` for seamless subscription resuming
- Handle message type "next" for log data, ignore other message types

See **Task 14-2** for implementation of the Railway log client in Go using this guide and Railway's proven patterns.


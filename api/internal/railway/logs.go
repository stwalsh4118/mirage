package railway

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coder/websocket"
)

//go:embed queries/subscriptions/environment-logs.graphql
var environmentLogsSubscription string

// EnvironmentLogsSubscriptionVariables matches Railway's expected structure
type EnvironmentLogsSubscriptionVariables struct {
	EnvironmentID string `json:"environmentId"`
	Filter        string `json:"filter,omitempty"`
	BeforeDate    string `json:"beforeDate,omitempty"` // ISO 8601 / RFC3339Nano
	BeforeLimit   int    `json:"beforeLimit"`
	AnchorDate    string `json:"anchorDate,omitempty"`
	AfterDate     string `json:"afterDate,omitempty"`
	AfterLimit    int    `json:"afterLimit,omitempty"`
}

// EnvironmentLogsSubscriptionPayload is the GraphQL subscription payload
type EnvironmentLogsSubscriptionPayload struct {
	Query     string                                `json:"query"`
	Variables *EnvironmentLogsSubscriptionVariables `json:"variables"`
}

// EnvironmentLog represents a single log entry from Railway
type EnvironmentLog struct {
	Timestamp  string              `json:"timestamp"`
	Message    string              `json:"message"`
	Severity   string              `json:"severity"`
	Tags       map[string]string   `json:"tags"`
	Attributes []map[string]string `json:"attributes"`
}

// SubscribeToEnvironmentLogs creates an environment logs subscription using Railway's pattern
func (c *Client) SubscribeToEnvironmentLogs(
	ctx context.Context,
	environmentID string,
	serviceFilter string,
) (*websocket.Conn, error) {
	payload := &EnvironmentLogsSubscriptionPayload{
		Query: environmentLogsSubscription,
		Variables: &EnvironmentLogsSubscriptionVariables{
			EnvironmentID: environmentID,
			Filter:        serviceFilter,
			// Get last 5 minutes of logs for seamless subscription resuming
			BeforeDate:  time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339Nano),
			BeforeLimit: 500,
		},
	}

	return c.createWebSocketSubscription(ctx, payload)
}

// ReadLogMessage reads and parses a log message from the WebSocket
func ReadLogMessage(ctx context.Context, conn *websocket.Conn) (*EnvironmentLog, error) {
	_, data, err := conn.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("read websocket: %w", err)
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
		return nil, fmt.Errorf("unmarshal message: %w", err)
	}

	if msg.Type != "next" {
		return nil, nil // Skip non-data messages
	}

	return &msg.Payload.Data.EnvironmentLogs, nil
}

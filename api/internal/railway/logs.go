package railway

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog/log"
)

//go:embed queries/subscriptions/environment-logs.graphql
var environmentLogsSubscription string

//go:embed queries/queries/deployment-logs.graphql
var deploymentLogsQuery string

//go:embed queries/queries/service-deployments.graphql
var serviceDeploymentsQuery string

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

// GetDeploymentLogsInput defines parameters for fetching historical deployment logs
type GetDeploymentLogsInput struct {
	DeploymentID string
	Limit        int    // Default: 500, Max: 1000
	Filter       string // Text filter for log messages
}

// LogTags represents the tags associated with a log entry
type LogTags struct {
	DeploymentID         string `json:"deploymentId"`
	DeploymentInstanceID string `json:"deploymentInstanceId"`
	EnvironmentID        string `json:"environmentId"`
	PluginID             string `json:"pluginId"`
	ProjectID            string `json:"projectId"`
	ServiceID            string `json:"serviceId"`
	SnapshotID           string `json:"snapshotId"`
}

// LogAttribute represents a key-value attribute pair
type LogAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DeploymentLog represents a single log entry from a deployment
type DeploymentLog struct {
	Timestamp  string         `json:"timestamp"`
	Message    string         `json:"message"`
	Severity   string         `json:"severity"`
	Tags       LogTags        `json:"tags"`
	Attributes []LogAttribute `json:"attributes"`
}

// GetDeploymentLogsResult contains the query result
type GetDeploymentLogsResult struct {
	Logs []DeploymentLog
}

// GetDeploymentLogs fetches historical logs for a specific deployment
func (c *Client) GetDeploymentLogs(ctx context.Context, input GetDeploymentLogsInput) (GetDeploymentLogsResult, error) {
	if input.Limit <= 0 {
		input.Limit = 500 // Default limit
	}
	if input.Limit > 1000 {
		input.Limit = 1000 // Max limit
	}

	vars := map[string]any{
		"deploymentId": input.DeploymentID,
		"limit":        input.Limit,
	}
	if input.Filter != "" {
		vars["filter"] = input.Filter
	}

	var out struct {
		DeploymentLogs []DeploymentLog `json:"deploymentLogs"`
	}

	log.Info().Msgf("querying deployment logs for deployment %s with limit %d and filter %s", input.DeploymentID, input.Limit, input.Filter)

	if err := c.execute(ctx, deploymentLogsQuery, vars, &out); err != nil {
		return GetDeploymentLogsResult{}, fmt.Errorf("query deployment logs: %w", err)
	}

	return GetDeploymentLogsResult{
		Logs: out.DeploymentLogs,
	}, nil
}

// Deployment represents a Railway deployment
type Deployment struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	EnvironmentID string `json:"environmentId"`
}

// GetLatestDeploymentID fetches the most recent deployment ID for a service
func (c *Client) GetLatestDeploymentID(ctx context.Context, serviceID string) (string, error) {
	vars := map[string]any{
		"serviceId": serviceID,
	}

	var out struct {
		Service struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ProjectID   string `json:"projectId"`
			Deployments struct {
				Edges []struct {
					Node Deployment `json:"node"`
				} `json:"edges"`
			} `json:"deployments"`
		} `json:"service"`
	}

	if err := c.execute(ctx, serviceDeploymentsQuery, vars, &out); err != nil {
		return "", fmt.Errorf("query service deployments: %w", err)
	}

	if len(out.Service.Deployments.Edges) == 0 {
		return "", fmt.Errorf("no deployments found for service %s", serviceID)
	}

	return out.Service.Deployments.Edges[0].Node.ID, nil
}

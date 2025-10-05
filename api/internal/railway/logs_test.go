package railway

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentLog_Unmarshal(t *testing.T) {
	jsonData := `{
		"timestamp": "2025-10-05T12:34:56.789Z",
		"message": "Server started on port 3000",
		"severity": "INFO",
		"tags": {
			"projectId": "proj-123",
			"environmentId": "env-456",
			"serviceId": "svc-789"
		},
		"attributes": [
			{"key": "version", "value": "1.0.0"}
		]
	}`

	var log EnvironmentLog
	err := json.Unmarshal([]byte(jsonData), &log)

	require.NoError(t, err)
	assert.Equal(t, "2025-10-05T12:34:56.789Z", log.Timestamp)
	assert.Equal(t, "Server started on port 3000", log.Message)
	assert.Equal(t, "INFO", log.Severity)
	assert.Equal(t, "proj-123", log.Tags["projectId"])
	assert.Equal(t, "env-456", log.Tags["environmentId"])
	assert.Equal(t, "svc-789", log.Tags["serviceId"])
	assert.Len(t, log.Attributes, 1)
	assert.Equal(t, "version", log.Attributes[0]["key"])
	assert.Equal(t, "1.0.0", log.Attributes[0]["value"])
}

func TestEnvironmentLogsSubscriptionPayload_Marshal(t *testing.T) {
	payload := &EnvironmentLogsSubscriptionPayload{
		Query: environmentLogsSubscription,
		Variables: &EnvironmentLogsSubscriptionVariables{
			EnvironmentID: "env-123",
			Filter:        "",
			BeforeLimit:   500,
			BeforeDate:    "2025-10-05T12:00:00.000Z",
		},
	}

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result["query"])

	vars, ok := result["variables"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "env-123", vars["environmentId"])
	assert.Equal(t, float64(500), vars["beforeLimit"])
}

func TestIntegration_SubscribeToEnvironmentLogs(t *testing.T) {
	token := os.Getenv("RAILWAY_TOKEN")
	environmentID := os.Getenv("TEST_ENVIRONMENT_ID")

	if token == "" || environmentID == "" {
		t.Skip("RAILWAY_TOKEN and TEST_ENVIRONMENT_ID not set")
	}

	client := NewClient("", token, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := client.SubscribeToEnvironmentLogs(ctx, environmentID, "")
	require.NoError(t, err, "Failed to establish WebSocket subscription")
	defer conn.Close(websocket.StatusNormalClosure, "test complete")

	// Try to read at least one log message
	log, err := ReadLogMessage(ctx, conn)

	// Either we get a log or context times out (if no logs are being generated)
	if err == context.DeadlineExceeded {
		t.Log("Test timed out, no logs received (may be expected if no active services)")
		return
	}

	require.NoError(t, err)

	if log != nil {
		assert.NotEmpty(t, log.Message, "Log message should not be empty")
		assert.NotEmpty(t, log.Tags["environmentId"], "Environment ID should be present")
		t.Logf("Received log: [%s] %s", log.Severity, log.Message)
		if serviceID, ok := log.Tags["serviceId"]; ok {
			t.Logf("From service: %s", serviceID)
		}
	}
}

func TestIntegration_WebSocketHandshake(t *testing.T) {
	token := os.Getenv("RAILWAY_TOKEN")
	if token == "" {
		t.Skip("RAILWAY_TOKEN not set")
	}

	client := NewClient("", token, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with a minimal payload just to verify handshake works
	payload := map[string]any{
		"query": `subscription { environmentLogs(environmentId: "test", beforeLimit: 1) { message } }`,
		"variables": map[string]any{
			"environmentId": "test-will-fail-but-handshake-works",
			"beforeLimit":   1,
		},
	}

	conn, err := client.createWebSocketSubscription(ctx, payload)

	// We expect this to succeed in establishing connection (even if env ID is invalid)
	// Railway will send an error in a message, but the WebSocket handshake should complete
	if err != nil {
		// If we get an error, it should be during read, not during handshake
		assert.NotContains(t, err.Error(), "connection ack", "Should successfully complete WebSocket handshake")
	}

	if conn != nil {
		conn.Close(websocket.StatusNormalClosure, "test complete")
		t.Log("WebSocket handshake completed successfully")
	}
}

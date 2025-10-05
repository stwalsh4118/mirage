package railway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	subscriptionTypeSubscribe = "subscribe"
	subscriptionTypeComplete  = "complete"
)

var (
	connectionInit = []byte(`{"type":"connection_init"}`)
	connectionAck  = []byte(`{"type":"connection_ack"}`)
)

// createWebSocketSubscription establishes a WebSocket connection following Railway's pattern
// This is a generic method for any GraphQL subscription to Railway's API
func (c *Client) createWebSocketSubscription(ctx context.Context, payload any) (*websocket.Conn, error) {
	// Generate unique subscription ID
	subID := uuid.New()

	subPayload := map[string]any{
		"id":      subID,
		"type":    subscriptionTypeSubscribe,
		"payload": payload,
	}

	payloadBytes, err := json.Marshal(&subPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal subscription payload: %w", err)
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
		return nil, fmt.Errorf("dial websocket: %w", err)
	}

	// Remove read limit (Railway sends large log batches)
	conn.SetReadLimit(-1)

	// Step 1: Send connection_init
	if err := conn.Write(ctx, websocket.MessageText, connectionInit); err != nil {
		conn.Close(websocket.StatusInternalError, "failed to send init")
		return nil, fmt.Errorf("write connection_init: %w", err)
	}

	// Step 2: Wait for connection_ack
	_, ackMessage, err := conn.Read(ctx)
	if err != nil {
		conn.Close(websocket.StatusInternalError, "failed to read ack")
		return nil, fmt.Errorf("read connection_ack: %w", err)
	}

	if !bytes.Equal(ackMessage, connectionAck) {
		conn.Close(websocket.StatusProtocolError, "invalid ack")
		return nil, errors.New("did not receive connection ack from server")
	}

	// Step 3: Send subscription payload
	if err := conn.Write(ctx, websocket.MessageText, payloadBytes); err != nil {
		conn.Close(websocket.StatusInternalError, "failed to send subscription")
		return nil, fmt.Errorf("write subscription: %w", err)
	}

	return conn, nil
}

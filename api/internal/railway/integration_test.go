package railway

import (
	"context"
	"os"
	"testing"
)

func TestIntegration_Viewer_Smoke(t *testing.T) {
	token := os.Getenv("RAILWAY_API_TOKEN")
	if token == "" {
		t.Skip("RAILWAY_API_TOKEN not set; skipping integration test")
	}
	c := NewClient("", token, nil)
	var out struct {
		Viewer struct{ ID string } `json:"viewer"`
	}
	err := c.execute(context.Background(), `query { viewer { id } }`, nil, &out)
	if err != nil {
		t.Fatalf("viewer query failed: %v", err)
	}
	if out.Viewer.ID == "" {
		t.Fatalf("expected viewer id, got empty")
	}
}

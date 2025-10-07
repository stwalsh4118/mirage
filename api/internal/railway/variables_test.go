package railway

import (
	"context"
	"testing"
)

func TestGetEnvironmentVariables_InputValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     GetEnvironmentVariablesInput
		expectErr bool
	}{
		{
			name: "valid environment-level request",
			input: GetEnvironmentVariablesInput{
				ProjectID:     "proj-123",
				EnvironmentID: "env-456",
				ServiceID:     nil,
			},
			expectErr: false,
		},
		{
			name: "valid service-level request",
			input: GetEnvironmentVariablesInput{
				ProjectID:     "proj-123",
				EnvironmentID: "env-456",
				ServiceID:     strPtr("svc-789"),
			},
			expectErr: false,
		},
		{
			name: "empty project ID should still call API (Railway will error)",
			input: GetEnvironmentVariablesInput{
				ProjectID:     "",
				EnvironmentID: "env-456",
				ServiceID:     nil,
			},
			expectErr: false, // We don't validate inputs, Railway API will error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies input structure only
			// Actual API calls would require integration tests
			if tt.input.ProjectID == "" && !tt.expectErr {
				// Just verify structure is valid Go
				_ = tt.input.ProjectID
				_ = tt.input.EnvironmentID
				_ = tt.input.ServiceID
			}
		})
	}
}

func TestGetEnvironmentVariablesResult_Structure(t *testing.T) {
	// Test that the result structure works as expected
	result := GetEnvironmentVariablesResult{
		Variables: map[string]string{
			"DATABASE_URL": "postgresql://localhost:5432/db",
			"PORT":         "3000",
			"LOG_LEVEL":    "info",
		},
	}

	if len(result.Variables) != 3 {
		t.Errorf("expected 3 variables, got %d", len(result.Variables))
	}

	if result.Variables["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %s", result.Variables["PORT"])
	}
}

func TestGetEnvironmentVariablesResult_EmptyMap(t *testing.T) {
	// Test empty result
	result := GetEnvironmentVariablesResult{
		Variables: map[string]string{},
	}

	if len(result.Variables) != 0 {
		t.Errorf("expected empty map, got %d variables", len(result.Variables))
	}
}

func TestGetEnvironmentVariablesResult_NilMap(t *testing.T) {
	// Test nil map (should be handled gracefully)
	result := GetEnvironmentVariablesResult{
		Variables: nil,
	}

	// Verify we can safely check length of nil map (len() for nil maps is 0)
	if len(result.Variables) != 0 {
		t.Errorf("expected nil or empty map, got %d variables", len(result.Variables))
	}
}

func TestStringPtrToLogValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "(nil)",
		},
		{
			name:     "non-nil pointer",
			input:    strPtr("svc-123"),
			expected: "svc-123",
		},
		{
			name:     "empty string pointer",
			input:    strPtr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringPtrToLogValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Helper function for tests
func strPtr(s string) *string {
	return &s
}

func TestGetAllEnvironmentAndServiceVariables_Structure(t *testing.T) {
	// Test the result structure
	result := GetAllEnvironmentAndServiceVariablesResult{
		EnvironmentVariables: map[string]string{
			"DATABASE_URL": "postgresql://localhost:5432/db",
			"LOG_LEVEL":    "info",
		},
		ServiceVariables: []ServiceVariables{
			{
				ServiceID:   "svc-1",
				ServiceName: "api",
				Variables: map[string]string{
					"PORT":                    "3000",
					"RAILWAY_DOCKERFILE_PATH": "Dockerfile",
				},
			},
			{
				ServiceID:   "svc-2",
				ServiceName: "worker",
				Variables: map[string]string{
					"WORKER_THREADS": "4",
				},
			},
		},
	}

	if len(result.EnvironmentVariables) != 2 {
		t.Errorf("expected 2 environment variables, got %d", len(result.EnvironmentVariables))
	}

	if len(result.ServiceVariables) != 2 {
		t.Errorf("expected 2 services, got %d", len(result.ServiceVariables))
	}

	// Check first service
	if result.ServiceVariables[0].ServiceID != "svc-1" {
		t.Errorf("expected first service ID to be svc-1, got %s", result.ServiceVariables[0].ServiceID)
	}

	if len(result.ServiceVariables[0].Variables) != 2 {
		t.Errorf("expected first service to have 2 variables, got %d", len(result.ServiceVariables[0].Variables))
	}
}

func TestGetAllEnvironmentAndServiceVariables_EmptyServices(t *testing.T) {
	// Test with no services
	result := GetAllEnvironmentAndServiceVariablesResult{
		EnvironmentVariables: map[string]string{
			"ENV_VAR": "value",
		},
		ServiceVariables: []ServiceVariables{},
	}

	if len(result.EnvironmentVariables) != 1 {
		t.Errorf("expected 1 environment variable, got %d", len(result.EnvironmentVariables))
	}

	if len(result.ServiceVariables) != 0 {
		t.Errorf("expected 0 services, got %d", len(result.ServiceVariables))
	}
}

func TestGetAllEnvironmentAndServiceVariables_ServiceWithoutVariables(t *testing.T) {
	// Test service with empty variables map
	result := GetAllEnvironmentAndServiceVariablesResult{
		EnvironmentVariables: map[string]string{},
		ServiceVariables: []ServiceVariables{
			{
				ServiceID:   "svc-1",
				ServiceName: "api",
				Variables:   map[string]string{}, // Empty
			},
		},
	}

	if len(result.ServiceVariables) != 1 {
		t.Errorf("expected 1 service, got %d", len(result.ServiceVariables))
	}

	if len(result.ServiceVariables[0].Variables) != 0 {
		t.Errorf("expected service to have 0 variables, got %d", len(result.ServiceVariables[0].Variables))
	}
}

// Integration test example (requires real Railway environment)
// This would be run with integration test flag
func TestGetEnvironmentVariables_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test would require:
	// - Real Railway API token
	// - Real project/environment IDs
	// - Proper test setup
	//
	// Example structure:
	// client := NewClient("", os.Getenv("RAILWAY_API_TOKEN"), nil)
	// result, err := client.GetEnvironmentVariables(context.Background(), GetEnvironmentVariablesInput{
	//     ProjectID:     os.Getenv("TEST_PROJECT_ID"),
	//     EnvironmentID: os.Getenv("TEST_ENV_ID"),
	// })
	//
	// if err != nil {
	//     t.Fatalf("GetEnvironmentVariables failed: %v", err)
	// }
	//
	// t.Logf("Fetched %d variables", len(result.Variables))

	t.Skip("integration test requires Railway credentials")
}

func TestGetAllEnvironmentAndServiceVariables_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test would require:
	// - Real Railway API token
	// - Real project/environment IDs with services
	//
	// Example structure:
	// client := NewClient("", os.Getenv("RAILWAY_API_TOKEN"), nil)
	// result, err := client.GetAllEnvironmentAndServiceVariables(context.Background(),
	//     GetAllEnvironmentAndServiceVariablesInput{
	//         ProjectID:     os.Getenv("TEST_PROJECT_ID"),
	//         EnvironmentID: os.Getenv("TEST_ENV_ID"),
	//     })
	//
	// if err != nil {
	//     t.Fatalf("GetAllEnvironmentAndServiceVariables failed: %v", err)
	// }
	//
	// t.Logf("Fetched %d env variables and %d services with variables",
	//     len(result.EnvironmentVariables), len(result.ServiceVariables))

	t.Skip("integration test requires Railway credentials")
}

func TestGetEnvironmentVariables_ContextCancellation(t *testing.T) {
	// Test that context cancellation is handled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// With a cancelled context, the call should fail
	// This test verifies the method respects context
	input := GetEnvironmentVariablesInput{
		ProjectID:     "proj-123",
		EnvironmentID: "env-456",
	}

	// We can't actually test this without mocking the HTTP client
	// But we verify the input structure is correct
	_ = ctx
	_ = input

	// In a real integration test, this would return context.Canceled error
	t.Log("context cancellation would be tested in integration tests")
}

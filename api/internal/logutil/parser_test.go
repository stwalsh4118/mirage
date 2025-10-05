package logutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLogLine_JSON(t *testing.T) {
	line := `{"timestamp":"2025-01-01T12:00:00Z","level":"INFO","msg":"Server started"}`
	result := ParseLogLine(line, "api")

	assert.Equal(t, SeverityInfo, result.Severity)
	assert.Equal(t, "Server started", result.Message)
	assert.Equal(t, "api", result.ServiceName)
	assert.NotNil(t, result.Structured)
	assert.Equal(t, line, result.RawLine)

	expectedTime, _ := time.Parse(time.RFC3339, "2025-01-01T12:00:00Z")
	assert.True(t, result.Timestamp.Equal(expectedTime))
}

func TestParseLogLine_PlainText(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectedSeverity string
		expectsMessage   string
	}{
		{
			name:             "error with bracket",
			line:             "[ERROR] 2025-01-01 12:00:00 Database connection failed",
			expectedSeverity: SeverityError,
			expectsMessage:   "Database connection failed",
		},
		{
			name:             "info message",
			line:             "INFO: Server listening on port 8080",
			expectedSeverity: SeverityInfo,
			expectsMessage:   "Server listening on port 8080",
		},
		{
			name:             "warning",
			line:             "WARN Deprecated API endpoint used",
			expectedSeverity: SeverityWarn,
		},
		{
			name:             "plain message",
			line:             "Simple log message",
			expectedSeverity: SeverityUnknown,
			expectsMessage:   "Simple log message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLogLine(tt.line, "test-service")
			assert.Equal(t, tt.expectedSeverity, result.Severity)
			assert.Equal(t, "test-service", result.ServiceName)
			if tt.expectsMessage != "" {
				assert.Contains(t, result.Message, tt.expectsMessage)
			}
		})
	}
}

func TestParseLogLine_Logfmt(t *testing.T) {
	line := `level=info msg="Request received" method=GET path=/api/users`
	result := ParseLogLine(line, "web")

	assert.Equal(t, SeverityInfo, result.Severity)
	assert.Equal(t, "Request received", result.Message)
	assert.Equal(t, "web", result.ServiceName)
	assert.NotNil(t, result.Structured)
}

func TestParseLogLine_WithANSI(t *testing.T) {
	line := "\x1b[31m[ERROR]\x1b[0m Connection timeout"
	result := ParseLogLine(line, "api")

	assert.Equal(t, SeverityError, result.Severity)
	assert.Contains(t, result.Message, "Connection timeout")
	// Raw line should preserve ANSI
	assert.Contains(t, result.RawLine, "\x1b[31m")
	// Message should have ANSI stripped
	assert.NotContains(t, result.Message, "\x1b[")
}

func TestParseJSONLog(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectedSeverity string
		expectedMessage  string
		shouldError      bool
	}{
		{
			name:             "standard JSON log",
			line:             `{"timestamp":"2025-01-01T12:00:00Z","level":"error","message":"Failed to connect"}`,
			expectedSeverity: SeverityError,
			expectedMessage:  "Failed to connect",
			shouldError:      false,
		},
		{
			name:             "with msg field",
			line:             `{"ts":"2025-01-01T12:00:00Z","lvl":"warn","msg":"Low memory"}`,
			expectedSeverity: SeverityWarn,
			expectedMessage:  "Low memory",
			shouldError:      false,
		},
		{
			name:        "invalid JSON",
			line:        `{invalid json}`,
			shouldError: true,
		},
		{
			name:             "nested fields",
			line:             `{"time":"2025-01-01T12:00:00Z","severity":"INFO","log":"Starting server","extra":{"key":"value"}}`,
			expectedSeverity: SeverityInfo,
			expectedMessage:  "Starting server",
			shouldError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSONLog(tt.line, "service")

			if tt.shouldError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSeverity, result.Severity)
			assert.Equal(t, tt.expectedMessage, result.Message)
			assert.Equal(t, "service", result.ServiceName)
		})
	}
}

func TestParseLogfmt(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectedSeverity string
		expectedMessage  string
	}{
		{
			name:             "standard logfmt",
			line:             `level=error msg="Connection failed" host=localhost`,
			expectedSeverity: SeverityError,
			expectedMessage:  "Connection failed",
		},
		{
			name:             "with timestamp",
			line:             `time=2025-01-01T12:00:00Z level=info msg="Server started"`,
			expectedSeverity: SeverityInfo,
			expectedMessage:  "Server started",
		},
		{
			name:             "mixed format",
			line:             `level=warn some plain text msg="Warning message"`,
			expectedSeverity: SeverityWarn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLogfmt(tt.line, "api")

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSeverity, result.Severity)
			if tt.expectedMessage != "" {
				assert.Equal(t, tt.expectedMessage, result.Message)
			}
			assert.Equal(t, "api", result.ServiceName)
		})
	}
}

func TestTimestampExtraction(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		hasTimestamp bool
	}{
		{
			name:         "ISO 8601",
			line:         "2025-01-01T12:00:00Z INFO Server started",
			hasTimestamp: true,
		},
		{
			name:         "YYYY-MM-DD HH:MM:SS",
			line:         "2025-01-01 12:00:00 INFO Server started",
			hasTimestamp: true,
		},
		{
			name:         "no timestamp",
			line:         "INFO Server started",
			hasTimestamp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLogLine(tt.line, "api")

			if tt.hasTimestamp {
				// Should have extracted a timestamp (not default)
				assert.False(t, result.Timestamp.IsZero())
			}
			assert.NotEmpty(t, result.Message)
		})
	}
}

func TestParseLogLine_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{name: "empty string", line: ""},
		{name: "whitespace only", line: "   "},
		{name: "single word", line: "word"},
		{name: "very long line", line: string(make([]byte, 10000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := ParseLogLine(tt.line, "test")
			assert.NotNil(t, result)
		})
	}
}

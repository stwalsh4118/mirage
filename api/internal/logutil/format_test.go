package logutil

import (
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatAsJSON(t *testing.T) {
	logs := []ParsedLog{
		{
			Timestamp:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			Severity:    SeverityInfo,
			Message:     "Test message",
			ServiceName: "api",
		},
	}

	data, err := FormatAsJSON(logs)
	require.NoError(t, err)

	// Verify it's valid JSON
	var result []ParsedLog
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, "Test message", result[0].Message)
	assert.Equal(t, SeverityInfo, result[0].Severity)
}

func TestFormatAsCSV(t *testing.T) {
	time1 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	logs := []ParsedLog{
		{
			Timestamp:   time1,
			Severity:    SeverityInfo,
			Message:     "Test message",
			ServiceName: "api",
		},
		{
			Timestamp:   time1.Add(time.Minute),
			Severity:    SeverityError,
			Message:     "Error message",
			ServiceName: "web",
		},
	}

	data, err := FormatAsCSV(logs)
	require.NoError(t, err)

	// Verify CSV structure
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Check header
	assert.Len(t, records, 3) // Header + 2 rows
	assert.Equal(t, []string{"Timestamp", "Service", "Level", "Message"}, records[0])

	// Check first data row
	assert.Equal(t, "api", records[1][1])
	assert.Equal(t, SeverityInfo, records[1][2])
	assert.Equal(t, "Test message", records[1][3])

	// Check second data row
	assert.Equal(t, "web", records[2][1])
	assert.Equal(t, SeverityError, records[2][2])
	assert.Equal(t, "Error message", records[2][3])
}

func TestFormatAsCSV_WithANSI(t *testing.T) {
	logs := []ParsedLog{
		{
			Timestamp:   time.Now(),
			Severity:    SeverityError,
			Message:     "\x1b[31mError\x1b[0m message",
			ServiceName: "api",
		},
	}

	data, err := FormatAsCSV(logs)
	require.NoError(t, err)

	// ANSI codes should be stripped in CSV
	assert.NotContains(t, string(data), "\x1b[")
	assert.Contains(t, string(data), "Error message")
}

func TestFormatAsPlainText(t *testing.T) {
	time1 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	logs := []ParsedLog{
		{
			Timestamp:   time1,
			Severity:    SeverityInfo,
			Message:     "Server started",
			ServiceName: "api",
		},
		{
			Timestamp:   time1.Add(time.Minute),
			Severity:    SeverityError,
			Message:     "Connection failed",
			ServiceName: "web",
		},
	}

	data, err := FormatAsPlainText(logs)
	require.NoError(t, err)

	text := string(data)
	lines := strings.Split(strings.TrimSpace(text), "\n")

	assert.Len(t, lines, 2)
	assert.Contains(t, lines[0], "INFO")
	assert.Contains(t, lines[0], "[api]")
	assert.Contains(t, lines[0], "Server started")

	assert.Contains(t, lines[1], "ERROR")
	assert.Contains(t, lines[1], "[web]")
	assert.Contains(t, lines[1], "Connection failed")
}

func TestFormatAsNDJSON(t *testing.T) {
	logs := []ParsedLog{
		{
			Timestamp:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			Severity:    SeverityInfo,
			Message:     "Message 1",
			ServiceName: "api",
		},
		{
			Timestamp:   time.Date(2025, 1, 1, 12, 1, 0, 0, time.UTC),
			Severity:    SeverityError,
			Message:     "Message 2",
			ServiceName: "web",
		},
	}

	data, err := FormatAsNDJSON(logs)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	assert.Len(t, lines, 2)

	// Verify each line is valid JSON
	for i, line := range lines {
		var log ParsedLog
		err := json.Unmarshal([]byte(line), &log)
		require.NoError(t, err, "Line %d should be valid JSON", i)
	}
}

func TestTruncateMessage(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		maxLength int
		expected  string
	}{
		{
			name:      "shorter than max",
			message:   "Short",
			maxLength: 10,
			expected:  "Short",
		},
		{
			name:      "exactly max",
			message:   "ExactlyTen",
			maxLength: 10,
			expected:  "ExactlyTen",
		},
		{
			name:      "longer than max",
			message:   "This is a very long message",
			maxLength: 10,
			expected:  "This is...",
		},
		{
			name:      "max length too small",
			message:   "Test",
			maxLength: 2,
			expected:  "Te", // When maxLength < 3, just truncate
		},
		{
			name:      "empty message",
			message:   "",
			maxLength: 10,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateMessage(tt.message, tt.maxLength)
			assert.Equal(t, tt.expected, result)
			assert.LessOrEqual(t, len(result), tt.maxLength)
		})
	}
}

func TestFormatAsJSON_EmptyLogs(t *testing.T) {
	logs := []ParsedLog{}
	data, err := FormatAsJSON(logs)
	require.NoError(t, err)
	assert.Contains(t, string(data), "[]")
}

func TestFormatAsCSV_EmptyLogs(t *testing.T) {
	logs := []ParsedLog{}
	data, err := FormatAsCSV(logs)
	require.NoError(t, err)

	// Should still have header
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 1) // Just header
}

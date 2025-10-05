package logutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectSeverity(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "fatal error",
			message:  "FATAL: System crash",
			expected: SeverityFatal,
		},
		{
			name:     "error message",
			message:  "[ERROR] Database connection failed",
			expected: SeverityError,
		},
		{
			name:     "warning message",
			message:  "WARN: Low disk space",
			expected: SeverityWarn,
		},
		{
			name:     "info message",
			message:  "INFO: Server started on port 3000",
			expected: SeverityInfo,
		},
		{
			name:     "debug message",
			message:  "DEBUG: Processing request",
			expected: SeverityDebug,
		},
		{
			name:     "trace message",
			message:  "TRACE: Entering function",
			expected: SeverityTrace,
		},
		{
			name:     "emoji error",
			message:  "🔴 Critical failure",
			expected: SeverityError,
		},
		{
			name:     "emoji warning",
			message:  "⚠️ Deprecation notice",
			expected: SeverityWarn,
		},
		{
			name:     "emoji info",
			message:  "✅ Task completed",
			expected: SeverityInfo,
		},
		{
			name:     "emoji debug",
			message:  "🐛 Bug detected",
			expected: SeverityDebug,
		},
		{
			name:     "plain message",
			message:  "Server is running",
			expected: SeverityUnknown,
		},
		{
			name:     "empty message",
			message:  "",
			expected: SeverityUnknown,
		},
		{
			name:     "case insensitive",
			message:  "error: something broke",
			expected: SeverityError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectSeverity(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeSeverity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FATAL", SeverityFatal},
		{"fatal", SeverityFatal},
		{"CRITICAL", SeverityFatal},
		{"ERROR", SeverityError},
		{"err", SeverityError},
		{"FAILURE", SeverityError},
		{"WARNING", SeverityWarn},
		{"warn", SeverityWarn},
		{"INFO", SeverityInfo},
		{"information", SeverityInfo},
		{"DEBUG", SeverityDebug},
		{"dbg", SeverityDebug},
		{"TRACE", SeverityTrace},
		{"verbose", SeverityTrace},
		{"unknown", SeverityUnknown},
		{"", SeverityUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeSeverity(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeverityPriority(t *testing.T) {
	tests := []struct {
		severity string
		priority int
	}{
		{SeverityFatal, 6},
		{SeverityError, 5},
		{SeverityWarn, 4},
		{SeverityInfo, 3},
		{SeverityDebug, 2},
		{SeverityTrace, 1},
		{SeverityUnknown, 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			result := SeverityPriority(tt.severity)
			assert.Equal(t, tt.priority, result)
		})
	}

	// Test priority ordering
	assert.True(t, SeverityPriority(SeverityFatal) > SeverityPriority(SeverityError))
	assert.True(t, SeverityPriority(SeverityError) > SeverityPriority(SeverityWarn))
	assert.True(t, SeverityPriority(SeverityWarn) > SeverityPriority(SeverityInfo))
	assert.True(t, SeverityPriority(SeverityInfo) > SeverityPriority(SeverityDebug))
	assert.True(t, SeverityPriority(SeverityDebug) > SeverityPriority(SeverityTrace))
}

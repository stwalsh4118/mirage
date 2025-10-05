package logutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "red error",
			input:    "\x1b[31mError:\x1b[0m Something failed",
			expected: "Error: Something failed",
		},
		{
			name:     "yellow warning",
			input:    "\x1b[33mWarning:\x1b[0m Low memory",
			expected: "Warning: Low memory",
		},
		{
			name:     "multiple codes",
			input:    "\x1b[1m\x1b[32mSuccess\x1b[0m\x1b[0m",
			expected: "Success",
		},
		{
			name:     "no ANSI codes",
			input:    "Plain text message",
			expected: "Plain text message",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "bright colors",
			input:    "\x1b[91mBright Red\x1b[0m",
			expected: "Bright Red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripANSI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestANSIToSeverity(t *testing.T) {
	tests := []struct {
		name     string
		ansiCode string
		expected string
	}{
		{
			name:     "red - error",
			ansiCode: "\x1b[31m",
			expected: SeverityError,
		},
		{
			name:     "bright red - error",
			ansiCode: "\x1b[91m",
			expected: SeverityError,
		},
		{
			name:     "yellow - warning",
			ansiCode: "\x1b[33m",
			expected: SeverityWarn,
		},
		{
			name:     "green - info",
			ansiCode: "\x1b[32m",
			expected: SeverityInfo,
		},
		{
			name:     "blue - info",
			ansiCode: "\x1b[34m",
			expected: SeverityInfo,
		},
		{
			name:     "magenta - debug",
			ansiCode: "\x1b[35m",
			expected: SeverityDebug,
		},
		{
			name:     "cyan - debug",
			ansiCode: "\x1b[36m",
			expected: SeverityDebug,
		},
		{
			name:     "unknown code",
			ansiCode: "\x1b[39m",
			expected: SeverityUnknown,
		},
		{
			name:     "invalid format",
			ansiCode: "not an ansi code",
			expected: SeverityUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ANSIToSeverity(tt.ansiCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "has ANSI",
			input:    "\x1b[31mError\x1b[0m",
			expected: true,
		},
		{
			name:     "no ANSI",
			input:    "Plain text",
			expected: false,
		},
		{
			name:     "empty",
			input:    "",
			expected: false,
		},
		{
			name:     "partial escape",
			input:    "text\x1b[incomplete",
			expected: false, // Partial escape is not a valid ANSI code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasANSI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

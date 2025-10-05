package logutil

import (
	"regexp"
	"strings"
)

// Log severity level constants
const (
	SeverityTrace   = "TRACE"
	SeverityDebug   = "DEBUG"
	SeverityInfo    = "INFO"
	SeverityWarn    = "WARN"
	SeverityError   = "ERROR"
	SeverityFatal   = "FATAL"
	SeverityUnknown = "UNKNOWN"
)

// Severity level patterns in order of precedence
var severityPatterns = []struct {
	pattern  *regexp.Regexp
	severity string
}{
	{regexp.MustCompile(`(?i)\b(fatal|critical)\b`), SeverityFatal},
	{regexp.MustCompile(`(?i)\b(error|err|failure|failed|exception)\b`), SeverityError},
	{regexp.MustCompile(`(?i)\b(warn|warning)\b`), SeverityWarn},
	{regexp.MustCompile(`(?i)\b(info|information)\b`), SeverityInfo},
	{regexp.MustCompile(`(?i)\b(debug|dbg)\b`), SeverityDebug},
	{regexp.MustCompile(`(?i)\b(trace)\b`), SeverityTrace},
}

// Emoji indicators for severity
var emojiPatterns = map[string]string{
	"🔴":  SeverityError,
	"❌":  SeverityError,
	"⛔":  SeverityError,
	"⚠️": SeverityWarn,
	"🟡":  SeverityWarn,
	"ℹ️": SeverityInfo,
	"✅":  SeverityInfo,
	"🔵":  SeverityInfo,
	"🐛":  SeverityDebug,
	"🔍":  SeverityDebug,
}

// DetectSeverity attempts to determine the log severity from the message content
func DetectSeverity(message string) string {
	if message == "" {
		return SeverityUnknown
	}

	// Check for emoji indicators first
	for emoji, severity := range emojiPatterns {
		if strings.Contains(message, emoji) {
			return severity
		}
	}

	// Check for standard patterns
	for _, sp := range severityPatterns {
		if sp.pattern.MatchString(message) {
			return sp.severity
		}
	}

	return SeverityUnknown
}

// NormalizeSeverity standardizes various severity string formats
func NormalizeSeverity(severity string) string {
	switch strings.ToUpper(strings.TrimSpace(severity)) {
	case "FATAL", "CRITICAL", "CRIT":
		return SeverityFatal
	case "ERROR", "ERR", "FAILURE", "FAIL":
		return SeverityError
	case "WARNING", "WARN":
		return SeverityWarn
	case "INFO", "INFORMATION":
		return SeverityInfo
	case "DEBUG", "DBG":
		return SeverityDebug
	case "TRACE", "VERBOSE":
		return SeverityTrace
	default:
		return SeverityUnknown
	}
}

// SeverityPriority returns a numeric priority for severity levels (higher = more severe)
func SeverityPriority(severity string) int {
	switch severity {
	case SeverityFatal:
		return 6
	case SeverityError:
		return 5
	case SeverityWarn:
		return 4
	case SeverityInfo:
		return 3
	case SeverityDebug:
		return 2
	case SeverityTrace:
		return 1
	default:
		return 0
	}
}

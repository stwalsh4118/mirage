package logutil

import "regexp"

// ANSI escape code pattern
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes all ANSI escape codes from the given text
func StripANSI(text string) string {
	return ansiRegex.ReplaceAllString(text, "")
}

// Common ANSI color codes to severity mapping
var ansiColorToSeverity = map[string]string{
	"31": SeverityError, // Red
	"91": SeverityError, // Bright red
	"33": SeverityWarn,  // Yellow
	"93": SeverityWarn,  // Bright yellow
	"32": SeverityInfo,  // Green
	"92": SeverityInfo,  // Bright green
	"34": SeverityInfo,  // Blue
	"94": SeverityInfo,  // Bright blue
	"35": SeverityDebug, // Magenta
	"95": SeverityDebug, // Bright magenta
	"36": SeverityDebug, // Cyan
	"96": SeverityDebug, // Bright cyan
}

// ANSIToSeverity attempts to infer severity from ANSI color codes
func ANSIToSeverity(ansiCode string) string {
	// Extract color code from ANSI sequence (e.g., "\x1b[31m" -> "31")
	matches := regexp.MustCompile(`\x1b\[(\d+)m`).FindStringSubmatch(ansiCode)
	if len(matches) > 1 {
		if severity, ok := ansiColorToSeverity[matches[1]]; ok {
			return severity
		}
	}
	return SeverityUnknown
}

// HasANSI checks if the text contains any ANSI codes
func HasANSI(text string) bool {
	return ansiRegex.MatchString(text)
}

package logutil

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

// ParsedLog represents a processed log entry
type ParsedLog struct {
	Timestamp   time.Time              `json:"timestamp"`
	Severity    string                 `json:"severity"`
	Message     string                 `json:"message"`
	ServiceName string                 `json:"serviceName,omitempty"`
	RawLine     string                 `json:"rawLine,omitempty"`
	Structured  map[string]interface{} `json:"structured,omitempty"` // For JSON logs
}

// Common timestamp patterns
var timestampPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?`), // ISO 8601
	regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`),                                 // YYYY-MM-DD HH:MM:SS
	regexp.MustCompile(`\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}`),                                 // MM/DD/YYYY HH:MM:SS
}

// ParseLogLine attempts to parse a log line, auto-detecting the format
func ParseLogLine(line string, serviceName string) ParsedLog {
	// Strip ANSI codes for parsing
	cleanLine := StripANSI(line)

	// Try JSON first
	if strings.HasPrefix(strings.TrimSpace(cleanLine), "{") {
		if jsonLog, err := ParseJSONLog(cleanLine, serviceName); err == nil {
			return jsonLog
		}
	}

	// Try logfmt
	if isLogfmt(cleanLine) {
		if logfmtLog, err := ParseLogfmt(cleanLine, serviceName); err == nil {
			return logfmtLog
		}
	}

	// Fall back to plain text parsing
	return parsePlainText(line, serviceName)
}

// ParseJSONLog parses a JSON-formatted log line
func ParseJSONLog(line string, serviceName string) (ParsedLog, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return ParsedLog{}, err
	}

	parsed := ParsedLog{
		RawLine:     line,
		ServiceName: serviceName,
		Structured:  data,
		Timestamp:   time.Time{}, // Zero time indicates missing timestamp
	}

	// Extract timestamp (try common field names)
	if ts := extractJSONTimestamp(data); !ts.IsZero() {
		parsed.Timestamp = ts
	}

	// Extract severity (try common field names)
	if severity := extractJSONSeverity(data); severity != "" {
		parsed.Severity = NormalizeSeverity(severity)
	}

	// Extract message (try common field names)
	if msg := extractJSONMessage(data); msg != "" {
		parsed.Message = msg
	} else {
		parsed.Message = line // Use full line if no message field
	}

	// If no severity found, try to detect from message
	if parsed.Severity == "" || parsed.Severity == SeverityUnknown {
		parsed.Severity = DetectSeverity(parsed.Message)
	}

	return parsed, nil
}

// ParseLogfmt parses a logfmt-formatted log line (key=value pairs)
func ParseLogfmt(line string, serviceName string) (ParsedLog, error) {
	parsed := ParsedLog{
		RawLine:     line,
		ServiceName: serviceName,
		Timestamp:   time.Time{}, // Zero time indicates missing timestamp
		Structured:  make(map[string]interface{}),
	}

	// Parse key=value pairs, handling quoted values
	var messageParts []string
	i := 0
	for i < len(line) {
		// Skip whitespace
		for i < len(line) && line[i] == ' ' {
			i++
		}
		if i >= len(line) {
			break
		}

		// Find key=value pair
		eqIdx := strings.IndexByte(line[i:], '=')
		if eqIdx == -1 {
			// Rest of line is part of message
			messageParts = append(messageParts, line[i:])
			break
		}

		key := line[i : i+eqIdx]
		i += eqIdx + 1 // Move past '='

		// Extract value (handle quotes)
		var value string
		if i < len(line) && line[i] == '"' {
			// Quoted value
			i++ // Skip opening quote
			end := i
			for end < len(line) && line[end] != '"' {
				end++
			}
			value = line[i:end]
			if end < len(line) {
				i = end + 1 // Skip closing quote
			}
		} else {
			// Unquoted value (up to next space)
			end := i
			for end < len(line) && line[end] != ' ' {
				end++
			}
			value = line[i:end]
			i = end
		}

		parsed.Structured[key] = value

		// Check for common fields
		switch strings.ToLower(key) {
		case "level", "severity", "lvl":
			parsed.Severity = NormalizeSeverity(value)
		case "msg", "message":
			parsed.Message = value
		case "time", "timestamp", "ts":
			if ts, err := time.Parse(time.RFC3339, value); err == nil {
				parsed.Timestamp = ts
			}
		}
	}

	// If no message found, use remaining parts
	if parsed.Message == "" && len(messageParts) > 0 {
		parsed.Message = strings.Join(messageParts, " ")
	}

	// Try to detect severity from message if not found
	if parsed.Severity == "" || parsed.Severity == SeverityUnknown {
		parsed.Severity = DetectSeverity(line)
	}

	return parsed, nil
}

// parsePlainText parses plain text log lines
func parsePlainText(line string, serviceName string) ParsedLog {
	cleanLine := StripANSI(line)

	parsed := ParsedLog{
		RawLine:     line,
		Message:     cleanLine,
		ServiceName: serviceName,
		Timestamp:   time.Time{}, // Zero time indicates missing timestamp
	}

	// Try to extract timestamp
	for _, pattern := range timestampPatterns {
		if match := pattern.FindString(cleanLine); match != "" {
			if ts, err := parseTimestamp(match); err == nil {
				parsed.Timestamp = ts
				// Remove timestamp from message
				parsed.Message = strings.TrimSpace(strings.Replace(cleanLine, match, "", 1))
				break
			}
		}
	}

	// Detect severity
	parsed.Severity = DetectSeverity(cleanLine)

	return parsed
}

// Helper functions

func extractJSONTimestamp(data map[string]interface{}) time.Time {
	fields := []string{"timestamp", "time", "ts", "@timestamp", "datetime"}
	for _, field := range fields {
		if val, ok := data[field]; ok {
			if str, ok := val.(string); ok {
				if ts, err := parseTimestamp(str); err == nil {
					return ts
				}
			}
		}
	}
	return time.Time{}
}

func extractJSONSeverity(data map[string]interface{}) string {
	fields := []string{"level", "severity", "lvl", "loglevel", "log_level"}
	for _, field := range fields {
		if val, ok := data[field]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

func extractJSONMessage(data map[string]interface{}) string {
	fields := []string{"message", "msg", "text", "log"}
	for _, field := range fields {
		if val, ok := data[field]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

func parseTimestamp(str string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05",
		"01/02/2006 15:04:05",
	}

	for _, format := range formats {
		if ts, err := time.Parse(format, str); err == nil {
			return ts.UTC(), nil
		}
	}

	return time.Time{}, nil
}

func isLogfmt(line string) bool {
	// Simple heuristic: contains multiple key=value pairs
	return strings.Contains(line, "=") && len(strings.Fields(line)) >= 2
}

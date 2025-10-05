package logutil

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
)

// FormatAsJSON exports logs as JSON array
func FormatAsJSON(logs []ParsedLog) ([]byte, error) {
	return json.MarshalIndent(logs, "", "  ")
}

// FormatAsCSV exports logs as CSV with headers
func FormatAsCSV(logs []ParsedLog) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"Timestamp", "Service", "Level", "Message"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write rows
	for _, log := range logs {
		row := []string{
			log.Timestamp.Format("2006-01-02 15:04:05"),
			log.ServiceName,
			log.Severity,
			StripANSI(log.Message),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FormatAsPlainText exports logs as human-readable plain text
func FormatAsPlainText(logs []ParsedLog) ([]byte, error) {
	var buf bytes.Buffer

	for _, log := range logs {
		line := formatLogLine(log)
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}

// formatLogLine formats a single log entry for plain text output
func formatLogLine(log ParsedLog) string {
	timestamp := log.Timestamp.Format("2006-01-02 15:04:05")
	severity := padRight(log.Severity, 5)
	service := ""
	if log.ServiceName != "" {
		service = fmt.Sprintf("[%s] ", log.ServiceName)
	}

	message := StripANSI(log.Message)
	return fmt.Sprintf("%s %s %s%s", timestamp, severity, service, message)
}

// FormatAsNDJSON exports logs as newline-delimited JSON
func FormatAsNDJSON(logs []ParsedLog) ([]byte, error) {
	var buf bytes.Buffer

	for _, log := range logs {
		data, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}

// TruncateMessage truncates a message to the specified length with ellipsis
func TruncateMessage(message string, maxLength int) string {
	if len(message) <= maxLength {
		return message
	}
	if maxLength < 3 {
		// If maxLength is very small, just truncate without ellipsis
		if maxLength <= 0 {
			return ""
		}
		return message[:maxLength]
	}
	return message[:maxLength-3] + "..."
}

// padRight pads a string to the specified length with spaces
func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

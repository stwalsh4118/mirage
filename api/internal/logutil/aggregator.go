package logutil

import "sort"

// LogAggregator merges logs from multiple services in chronological order
type LogAggregator struct {
	logs []ParsedLog
}

// NewLogAggregator creates a new log aggregator
func NewLogAggregator() *LogAggregator {
	return &LogAggregator{
		logs: make([]ParsedLog, 0),
	}
}

// Add appends logs to the aggregator
func (a *LogAggregator) Add(logs []ParsedLog) {
	a.logs = append(a.logs, logs...)
}

// GetMerged returns all logs sorted by timestamp in chronological order
func (a *LogAggregator) GetMerged() []ParsedLog {
	// Create a copy and sort it
	result := make([]ParsedLog, len(a.logs))
	copy(result, a.logs)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// GetMergedFiltered returns logs filtered by severity and/or service
func (a *LogAggregator) GetMergedFiltered(minSeverity string, serviceNames []string) []ParsedLog {
	merged := a.GetMerged()

	if minSeverity == "" && len(serviceNames) == 0 {
		return merged
	}

	minPriority := SeverityPriority(minSeverity)
	serviceFilter := make(map[string]bool)
	for _, name := range serviceNames {
		serviceFilter[name] = true
	}

	filtered := make([]ParsedLog, 0, len(merged))
	for _, log := range merged {
		// Check severity filter
		if minSeverity != "" && SeverityPriority(log.Severity) < minPriority {
			continue
		}

		// Check service filter
		if len(serviceNames) > 0 && !serviceFilter[log.ServiceName] {
			continue
		}

		filtered = append(filtered, log)
	}

	return filtered
}

// Count returns the total number of logs
func (a *LogAggregator) Count() int {
	return len(a.logs)
}

// Clear removes all logs from the aggregator
func (a *LogAggregator) Clear() {
	a.logs = make([]ParsedLog, 0)
}

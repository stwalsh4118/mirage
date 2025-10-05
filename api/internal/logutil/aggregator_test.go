package logutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogAggregator_ChronologicalMerge(t *testing.T) {
	agg := NewLogAggregator()

	time1 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2025, 1, 1, 12, 1, 0, 0, time.UTC)
	time3 := time.Date(2025, 1, 1, 12, 2, 0, 0, time.UTC)

	// Add logs from service A (out of order)
	logsA := []ParsedLog{
		{Timestamp: time3, ServiceName: "api", Message: "Third"},
		{Timestamp: time1, ServiceName: "api", Message: "First"},
	}

	// Add logs from service B
	logsB := []ParsedLog{
		{Timestamp: time2, ServiceName: "web", Message: "Second"},
	}

	agg.Add(logsA)
	agg.Add(logsB)

	merged := agg.GetMerged()

	// Verify chronological order
	assert.Len(t, merged, 3)
	assert.Equal(t, "First", merged[0].Message)
	assert.Equal(t, "Second", merged[1].Message)
	assert.Equal(t, "Third", merged[2].Message)
	assert.True(t, merged[0].Timestamp.Before(merged[1].Timestamp))
	assert.True(t, merged[1].Timestamp.Before(merged[2].Timestamp))
}

func TestLogAggregator_GetMergedFiltered_BySeverity(t *testing.T) {
	agg := NewLogAggregator()

	logs := []ParsedLog{
		{Timestamp: time.Now(), Severity: SeverityDebug, Message: "Debug"},
		{Timestamp: time.Now(), Severity: SeverityInfo, Message: "Info"},
		{Timestamp: time.Now(), Severity: SeverityWarn, Message: "Warn"},
		{Timestamp: time.Now(), Severity: SeverityError, Message: "Error"},
	}

	agg.Add(logs)

	// Filter for WARN and above
	filtered := agg.GetMergedFiltered(SeverityWarn, nil)

	assert.Len(t, filtered, 2) // Should have WARN and ERROR
	assert.Equal(t, "Warn", filtered[0].Message)
	assert.Equal(t, "Error", filtered[1].Message)
}

func TestLogAggregator_GetMergedFiltered_ByService(t *testing.T) {
	agg := NewLogAggregator()

	logs := []ParsedLog{
		{Timestamp: time.Now(), ServiceName: "api", Message: "API log"},
		{Timestamp: time.Now(), ServiceName: "web", Message: "Web log"},
		{Timestamp: time.Now(), ServiceName: "db", Message: "DB log"},
		{Timestamp: time.Now(), ServiceName: "api", Message: "API log 2"},
	}

	agg.Add(logs)

	// Filter for specific services
	filtered := agg.GetMergedFiltered("", []string{"api", "web"})

	assert.Len(t, filtered, 3)
	for _, log := range filtered {
		assert.Contains(t, []string{"api", "web"}, log.ServiceName)
	}
}

func TestLogAggregator_GetMergedFiltered_Combined(t *testing.T) {
	agg := NewLogAggregator()

	logs := []ParsedLog{
		{Timestamp: time.Now(), Severity: SeverityDebug, ServiceName: "api", Message: "1"},
		{Timestamp: time.Now(), Severity: SeverityError, ServiceName: "api", Message: "2"},
		{Timestamp: time.Now(), Severity: SeverityError, ServiceName: "web", Message: "3"},
		{Timestamp: time.Now(), Severity: SeverityWarn, ServiceName: "api", Message: "4"},
	}

	agg.Add(logs)

	// Filter for api service AND error+ severity
	filtered := agg.GetMergedFiltered(SeverityError, []string{"api"})

	assert.Len(t, filtered, 1)
	assert.Equal(t, "2", filtered[0].Message)
	assert.Equal(t, "api", filtered[0].ServiceName)
	assert.Equal(t, SeverityError, filtered[0].Severity)
}

func TestLogAggregator_Count(t *testing.T) {
	agg := NewLogAggregator()
	assert.Equal(t, 0, agg.Count())

	agg.Add([]ParsedLog{{}, {}, {}})
	assert.Equal(t, 3, agg.Count())

	agg.Add([]ParsedLog{{}, {}})
	assert.Equal(t, 5, agg.Count())
}

func TestLogAggregator_Clear(t *testing.T) {
	agg := NewLogAggregator()
	agg.Add([]ParsedLog{{}, {}, {}})

	assert.Equal(t, 3, agg.Count())

	agg.Clear()
	assert.Equal(t, 0, agg.Count())
	assert.Len(t, agg.GetMerged(), 0)
}

func TestLogAggregator_EmptyAggregator(t *testing.T) {
	agg := NewLogAggregator()

	merged := agg.GetMerged()
	assert.NotNil(t, merged)
	assert.Len(t, merged, 0)
}

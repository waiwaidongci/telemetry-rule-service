package events

import (
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"sort"
	"strings"
	"time"
)

type Query struct {
	Status   event.Status
	RuleID   string
	MetricID string
	SourceID string
	From     time.Time
	To       time.Time
	Limit    int
}

func (q Query) Match(value event.Event) bool {
	if q.Status != "" && value.Status != q.Status {
		return false
	}
	if q.RuleID != "" && value.RuleID != q.RuleID {
		return false
	}
	if q.MetricID != "" && value.MetricID != q.MetricID {
		return false
	}
	if q.SourceID != "" && value.SourceID != q.SourceID {
		return false
	}
	if !q.From.IsZero() && value.LastSeen.Before(q.From) {
		return false
	}
	if !q.To.IsZero() && value.LastSeen.After(q.To) {
		return false
	}
	return true
}

func ApplyQuery(values []event.Event, query Query) []event.Event {
	result := make([]event.Event, 0, len(values))
	for _, value := range values {
		if query.Match(value) {
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastSeen.After(result[j].LastSeen)
	})
	if query.Limit > 0 && len(result) > query.Limit {
		result = result[:query.Limit]
	}
	return result
}

func Search(values []event.Event, phrase string) []event.Event {
	phrase = strings.ToLower(strings.TrimSpace(phrase))
	if phrase == "" {
		return values
	}
	result := make([]event.Event, 0)
	for _, value := range values {
		if strings.Contains(strings.ToLower(value.Message), phrase) {
			result = append(result, value)
		}
	}
	return result
}

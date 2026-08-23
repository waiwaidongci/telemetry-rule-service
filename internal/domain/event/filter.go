package event

import "strings"

type Filter struct {
	Status   Status
	RuleID   string
	MetricID string
	SourceID string
	Query    string
}

func (f Filter) Match(e Event) bool {
	if f.Status != "" && f.Status != e.Status {
		return false
	}
	if f.RuleID != "" && f.RuleID != e.RuleID {
		return false
	}
	if f.MetricID != "" && f.MetricID != e.MetricID {
		return false
	}
	if f.SourceID != "" && f.SourceID != e.SourceID {
		return false
	}
	if f.Query != "" && !strings.Contains(strings.ToLower(e.Message), strings.ToLower(f.Query)) {
		return false
	}
	return true
}

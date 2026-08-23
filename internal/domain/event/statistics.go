package event

import "time"

type Statistics struct {
	Total        int       `json:"total"`
	Open         int       `json:"open"`
	Acknowledged int       `json:"acknowledged"`
	Resolved     int       `json:"resolved"`
	Oldest       time.Time `json:"oldest,omitempty"`
	Newest       time.Time `json:"newest,omitempty"`
}

func BuildStatistics(values []Event) Statistics {
	result := Statistics{Total: len(values)}
	for _, value := range values {
		switch value.Status {
		case Open:
			result.Open++
		case Acknowledged:
			result.Acknowledged++
		case Resolved:
			result.Resolved++
		}
		if result.Oldest.IsZero() || value.FirstSeen.Before(result.Oldest) {
			result.Oldest = value.FirstSeen
		}
		if value.LastSeen.After(result.Newest) {
			result.Newest = value.LastSeen
		}
	}
	return result
}

func CountByRule(values []Event) map[string]int {
	result := make(map[string]int)
	for _, value := range values {
		result[value.RuleID]++
	}
	return result
}

func CountByMetric(values []Event) map[string]int {
	result := make(map[string]int)
	for _, value := range values {
		result[value.MetricID]++
	}
	return result
}

func Active(values []Event) []Event {
	result := make([]Event, 0)
	for _, value := range values {
		if value.IsActive() {
			result = append(result, value)
		}
	}
	return result
}

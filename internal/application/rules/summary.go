package rules

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"time"
)

type Summary struct {
	RuleID    string  `json:"rule_id"`
	Samples   int     `json:"samples"`
	Triggered bool    `json:"triggered"`
	Average   float64 `json:"average"`
	Minimum   float64 `json:"minimum"`
	Maximum   float64 `json:"maximum"`
}

func (s *Service) Summarize(ctx context.Context, d rule.Definition, at time.Time) (Summary, error) {
	start, end := rule.WindowBounds(at, d.WindowSeconds)
	values, err := s.metrics.Samples(ctx, d.MetricID, start, end)
	if err != nil {
		return Summary{}, err
	}
	window := metric.NewWindow(end, end.Sub(start), values)
	return Summary{RuleID: d.ID, Samples: window.Count(), Average: window.Average(), Minimum: window.Min(), Maximum: window.Max()}, nil
}
func EnabledOnly(values []rule.Definition) []rule.Definition {
	out := make([]rule.Definition, 0, len(values))
	for _, v := range values {
		if v.Enabled {
			out = append(out, v)
		}
	}
	return out
}

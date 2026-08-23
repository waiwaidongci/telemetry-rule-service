package rules

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
)

type MetricStore interface {
	Samples(context.Context, string, time.Time, time.Time) ([]metric.Sample, error)
}
type RuleStore interface {
	Get(context.Context, string) (rule.Definition, error)
	List(context.Context) ([]rule.Definition, error)
}
type EventStore interface {
	Create(context.Context, event.Event) error
}
type Service struct {
	metrics MetricStore
	rules   RuleStore
	events  EventStore
}

func NewService(m MetricStore, r RuleStore, e EventStore) *Service {
	return &Service{metrics: m, rules: r, events: e}
}
func (s *Service) Evaluate(ctx context.Context, sample metric.Sample) ([]event.Event, error) {
	defs, err := s.rules.List(ctx)
	if err != nil {
		return nil, err
	}
	out := []event.Event{}
	for _, d := range defs {
		if !d.Enabled || d.MetricID != sample.MetricID {
			continue
		}
		trigger, msg := s.check(ctx, d, sample)
		if !trigger {
			continue
		}
		now := sample.Timestamp
		if now.IsZero() {
			now = time.Now().UTC()
		}
		ev := event.Event{ID: fmt.Sprintf("%s-%s-%d", d.ID, sample.SourceID, now.UnixNano()), RuleID: d.ID, MetricID: sample.MetricID, SourceID: sample.SourceID, Message: msg, Value: sample.Value, Status: event.Open, FirstSeen: now, LastSeen: now, Count: 1, Labels: sample.Tags}
		if err := s.events.Create(ctx, ev); err != nil {
			return out, fmt.Errorf("create event: %w", err)
		}
		out = append(out, ev)
	}
	return out, nil
}
func (s *Service) check(ctx context.Context, d rule.Definition, sample metric.Sample) (bool, string) {
	window := time.Duration(d.WindowSeconds) * time.Second
	samples, _ := s.metrics.Samples(ctx, d.MetricID, sample.Timestamp.Add(-window), sample.Timestamp)
	switch d.Type {
	case rule.Threshold:
		ok := compare(sample.Value, d.Operator, d.Threshold)
		return ok, fmt.Sprintf("threshold %s %v (limit %v)", d.Operator, sample.Value, d.Threshold)
	case rule.Rate:
		if len(samples) < 2 {
			return false, ""
		}
		prev := samples[len(samples)-2]
		dt := sample.Timestamp.Sub(prev.Timestamp).Seconds()
		if dt <= 0 {
			return false, ""
		}
		rate := math.Abs((sample.Value - prev.Value) / dt)
		ok := compare(rate, d.Operator, d.Threshold)
		return ok, fmt.Sprintf("rate %s %.4f (limit %.4f)", d.Operator, rate, d.Threshold)
	case rule.Consecutive:
		count := 0
		for i := len(samples) - 1; i >= 0; i-- {
			if compare(samples[i].Value, d.Operator, d.Threshold) {
				count++
			} else {
				break
			}
		}
		ok := count >= d.ConsecutiveCount
		return ok, fmt.Sprintf("consecutive matches %d (required %d)", count, d.ConsecutiveCount)
	case rule.Missing:
		return len(samples) == 0, "missing samples in window"
	}
	return false, ""
}
func compare(v float64, op string, t float64) bool {
	switch op {
	case ">":
		return v > t
	case ">=":
		return v >= t
	case "<":
		return v < t
	case "<=":
		return v <= t
	case "==", "=":
		return v == t
	case "!=":
		return v != t
	default:
		return false
	}
}

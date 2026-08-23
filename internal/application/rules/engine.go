package rules

import (
	"context"
	"fmt"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

type Engine struct {
	metrics MetricStore
	rules   RuleStore
	events  EventStore
}

func NewEngine(metrics MetricStore, rules RuleStore, events EventStore) *Engine {
	return &Engine{metrics: metrics, rules: rules, events: events}
}

func (e *Engine) Run(ctx context.Context, value metric.Sample) ([]event.Event, error) {
	definitions, err := e.rules.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}
	created := make([]event.Event, 0)
	for _, definition := range definitions {
		if !definition.Enabled || definition.MetricID != value.MetricID {
			continue
		}
		service := NewService(e.metrics, e.rules, e.events)
		matches, message := service.check(ctx, definition, value)
		if !matches {
			continue
		}
		at := value.Timestamp
		if at.IsZero() {
			at = time.Now().UTC()
		}
		item := event.Event{
			ID:        fmt.Sprintf("%s-%s-%d", definition.ID, value.SourceID, at.UnixNano()),
			RuleID:    definition.ID,
			MetricID:  value.MetricID,
			SourceID:  value.SourceID,
			Message:   message,
			Value:     value.Value,
			Status:    event.Open,
			FirstSeen: at,
			LastSeen:  at,
			Count:     1,
		}
		if err := e.events.Create(ctx, item); err != nil {
			return created, fmt.Errorf("persist evaluation event: %w", err)
		}
		created = append(created, item)
	}
	return created, nil
}

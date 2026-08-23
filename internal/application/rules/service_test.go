package rules

import (
	"context"
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/adapter/memory"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
)

func TestThresholdEvaluation(t *testing.T) {
	store := memory.NewStore()
	metrics := memory.MetricRepository{Store: store}
	rulesRepo := memory.RuleRepository{Store: store}
	events := memory.EventRepository{Store: store}
	ctx := context.Background()
	if err := rulesRepo.Create(ctx, rule.Definition{ID: "high", Name: "High", MetricID: "temperature", Type: rule.Threshold, Operator: ">", Threshold: 80, WindowSeconds: 60, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	sample := metric.Sample{SourceID: "s1", MetricID: "temperature", Value: 90, Timestamp: time.Now().UTC()}
	if err := metrics.Record(ctx, sample); err != nil {
		t.Fatal(err)
	}
	result, err := NewService(metrics, rulesRepo, events).Evaluate(ctx, sample)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected event, got %d", len(result))
	}
}

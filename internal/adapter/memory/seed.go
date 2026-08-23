package memory

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"github.com/example/telemetry-rule-service/internal/domain/source"
	"time"
)

func (s *Store) Seed(ctx context.Context) error {
	now := time.Now().UTC()
	defaults := []source.DataSource{{ID: "demo-source", Name: "Demo source", Protocol: "json", Enabled: true, CreatedAt: now}}
	for _, v := range defaults {
		if err := s.CreateSource(ctx, v); err != nil {
			return err
		}
	}
	definitions := []metric.Definition{{ID: "temperature", Name: "Temperature", Unit: "celsius", Enabled: true, CreatedAt: now}, {ID: "pressure", Name: "Pressure", Unit: "kpa", Enabled: true, CreatedAt: now}}
	for _, v := range definitions {
		if err := s.CreateMetric(ctx, v); err != nil {
			return err
		}
	}
	return s.CreateRule(ctx, rule.Definition{ID: "high-temperature", Name: "High temperature", MetricID: "temperature", Type: rule.Threshold, Operator: ">", Threshold: 80, WindowSeconds: 60, Enabled: true, CreatedAt: now})
}

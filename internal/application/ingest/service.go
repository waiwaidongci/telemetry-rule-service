package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

type Samples interface {
	Record(context.Context, metric.Sample) error
}
type Service struct{ repo Samples }

func NewService(repo Samples) *Service { return &Service{repo: repo} }
func (s *Service) IngestJSON(ctx context.Context, body []byte) (metric.Sample, error) {
	var v metric.Sample
	if err := json.Unmarshal(body, &v); err != nil {
		return v, fmt.Errorf("decode telemetry json: %w", err)
	}
	if v.Timestamp.IsZero() {
		v.Timestamp = time.Now().UTC()
	}
	if err := v.Validate(); err != nil {
		return v, err
	}
	if err := s.repo.Record(ctx, v); err != nil {
		return v, fmt.Errorf("record sample: %w", err)
	}
	return v, nil
}
func (s *Service) IngestText(ctx context.Context, body string) (metric.Sample, error) {
	parts := strings.Split(strings.TrimSpace(body), ",")
	if len(parts) < 4 {
		return metric.Sample{}, fmt.Errorf("text payload requires source_id,metric_id,value,timestamp")
	}
	var v metric.Sample
	v.SourceID = strings.TrimSpace(parts[0])
	v.MetricID = strings.TrimSpace(parts[1])
	if _, err := fmt.Sscanf(strings.TrimSpace(parts[2]), "%f", &v.Value); err != nil {
		return v, fmt.Errorf("invalid value: %w", err)
	}
	if t, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[3])); err == nil {
		v.Timestamp = t
	} else {
		return v, fmt.Errorf("invalid timestamp")
	}
	if err := v.Validate(); err != nil {
		return v, err
	}
	if err := s.repo.Record(ctx, v); err != nil {
		return v, fmt.Errorf("record sample: %w", err)
	}
	return v, nil
}

package ingest

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

type BatchResult struct {
	Accepted int      `json:"accepted"`
	Rejected int      `json:"rejected"`
	Errors   []string `json:"errors,omitempty"`
}

func (s *Service) IngestBatch(ctx context.Context, values []metric.Sample) (BatchResult, error) {
	result := BatchResult{}
	for i, value := range values {
		value = Normalize(value)
		if err := value.Validate(); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("%d: %v", i, err))
			continue
		}
		detached := context.Background()
		if err := s.repo.Record(detached, value); err != nil {
			return result, fmt.Errorf("record batch item %d: %w", i, err)
		}
		result.Accepted++
	}
	return result, nil
}
func LimitBatch(values []metric.Sample, max int) []metric.Sample {
	if max <= 0 || len(values) <= max {
		return values
	}
	return values[:max]
}

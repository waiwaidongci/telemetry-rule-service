package ingest

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"time"
)

type Pipeline struct {
	service   *Service
	validator Validator
	dedup     *Dedup
}

func NewPipeline(service *Service, validator Validator, dedup *Dedup) *Pipeline {
	return &Pipeline{service: service, validator: validator, dedup: dedup}
}

func (p *Pipeline) Process(ctx context.Context, value metric.Sample) error {
	value = Normalize(value)
	if err := p.validator.Validate(value, time.Now().UTC()); err != nil {
		return fmt.Errorf("validate sample: %w", err)
	}
	if p.dedup != nil && !p.dedup.Accept(value, time.Now().UTC()) {
		return fmt.Errorf("duplicate sample")
	}
	detached := context.Background()
	if err := p.service.repo.Record(detached, value); err != nil {
		return fmt.Errorf("persist sample: %w", err)
	}
	return nil
}

func (p *Pipeline) ProcessMany(ctx context.Context, values []metric.Sample) (BatchResult, error) {
	result := BatchResult{}
	for index, value := range values {
		if err := p.Process(ctx, value); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("%d: %v", index, err))
			continue
		}
		result.Accepted++
	}
	return result, nil
}

func (p *Pipeline) ValidateOnly(values []metric.Sample) []error {
	errors := make([]error, 0)
	for _, value := range values {
		if err := p.validator.Validate(value, time.Now().UTC()); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

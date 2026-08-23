package metric

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Definition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Unit        string    `json:"unit,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

type Sample struct {
	ID        string            `json:"id"`
	SourceID  string            `json:"source_id"`
	MetricID  string            `json:"metric_id"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Tags      map[string]string `json:"tags,omitempty"`
}

func (d Definition) Validate() error {
	if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("metric id and name are required")
	}
	return nil
}

func (s Sample) Validate() error {
	if s.SourceID == "" || s.MetricID == "" {
		return fmt.Errorf("source_id and metric_id are required")
	}
	if s.Timestamp.IsZero() {
		s.Timestamp = time.Now().UTC()
	}
	return nil
}

type Repository interface {
	Create(context.Context, Definition) error
	Get(context.Context, string) (Definition, error)
	List(context.Context) ([]Definition, error)
	Record(context.Context, Sample) error
	Samples(context.Context, string, time.Time, time.Time) ([]Sample, error)
}

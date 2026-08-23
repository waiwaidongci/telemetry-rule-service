package rule

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Type string

const (
	Threshold   Type = "threshold"
	Rate        Type = "rate"
	Consecutive Type = "consecutive"
	Missing     Type = "missing"
)

type Definition struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	MetricID         string    `json:"metric_id"`
	Type             Type      `json:"type"`
	Operator         string    `json:"operator,omitempty"`
	Threshold        float64   `json:"threshold,omitempty"`
	WindowSeconds    int       `json:"window_seconds,omitempty"`
	ConsecutiveCount int       `json:"consecutive_count,omitempty"`
	Enabled          bool      `json:"enabled"`
	Version          int       `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
}

func (r Definition) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.MetricID) == "" {
		return fmt.Errorf("rule id, name and metric_id are required")
	}
	switch r.Type {
	case Threshold, Rate, Consecutive, Missing:
	default:
		return fmt.Errorf("unsupported rule type %q", r.Type)
	}
	if r.WindowSeconds <= 0 {
		r.WindowSeconds = 60
	}
	if r.Type == Consecutive && r.ConsecutiveCount <= 0 {
		return fmt.Errorf("consecutive_count must be positive")
	}
	return nil
}

type Repository interface {
	Create(context.Context, Definition) error
	Get(context.Context, string) (Definition, error)
	List(context.Context) ([]Definition, error)
	Update(context.Context, Definition) error
}

package event

import (
	"context"
	"fmt"
	"time"
)

type Status string

const (
	Open         Status = "open"
	Acknowledged Status = "acknowledged"
	Resolved     Status = "resolved"
)

type Event struct {
	ID        string            `json:"id"`
	RuleID    string            `json:"rule_id"`
	MetricID  string            `json:"metric_id"`
	SourceID  string            `json:"source_id"`
	Message   string            `json:"message"`
	Value     float64           `json:"value"`
	Status    Status            `json:"status"`
	FirstSeen time.Time         `json:"first_seen"`
	LastSeen  time.Time         `json:"last_seen"`
	Count     int               `json:"count"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func (e Event) Validate() error {
	if e.ID == "" || e.RuleID == "" {
		return fmt.Errorf("event id and rule_id are required")
	}
	return nil
}

type Repository interface {
	Create(context.Context, Event) error
	Get(context.Context, string) (Event, error)
	List(context.Context, Status, int) ([]Event, error)
	Update(context.Context, Event) error
}

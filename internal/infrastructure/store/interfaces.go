package store

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"time"
)

type SampleStore interface {
	Record(context.Context, metric.Sample) error
	Samples(context.Context, string, time.Time, time.Time) ([]metric.Sample, error)
}
type HealthChecker interface{ Ping(context.Context) error }
type Transaction interface {
	Begin(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(context.Context) error
}
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

func DefaultMigrations() []Migration {
	return []Migration{{Version: 1, Name: "initial telemetry schema", Up: "create tables", Down: "drop tables"}}
}

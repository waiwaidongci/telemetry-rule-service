package memory

import (
	"context"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"github.com/example/telemetry-rule-service/internal/domain/source"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

type SourceRepository struct{ Store *Store }

func (r SourceRepository) Create(ctx context.Context, value source.DataSource) error {
	return r.Store.CreateSource(ctx, value)
}
func (r SourceRepository) Get(ctx context.Context, id string) (source.DataSource, error) {
	return r.Store.GetSource(ctx, id)
}
func (r SourceRepository) List(ctx context.Context) ([]source.DataSource, error) {
	return r.Store.ListSources(ctx)
}
func (r SourceRepository) Update(ctx context.Context, value source.DataSource) error {
	return r.Store.UpdateSource(ctx, value)
}

type MetricRepository struct{ Store *Store }

func (r MetricRepository) Create(ctx context.Context, value metric.Definition) error {
	return r.Store.CreateMetric(ctx, value)
}
func (r MetricRepository) Get(ctx context.Context, id string) (metric.Definition, error) {
	return r.Store.GetMetric(ctx, id)
}
func (r MetricRepository) List(ctx context.Context) ([]metric.Definition, error) {
	return r.Store.ListMetrics(ctx)
}
func (r MetricRepository) Record(ctx context.Context, value metric.Sample) error {
	return r.Store.RecordSample(ctx, value)
}
func (r MetricRepository) Samples(ctx context.Context, id string, from, to time.Time) ([]metric.Sample, error) {
	return r.Store.Samples(ctx, id, from, to)
}

type RuleRepository struct{ Store *Store }

func (r RuleRepository) Create(ctx context.Context, value rule.Definition) error {
	return r.Store.CreateRule(ctx, value)
}
func (r RuleRepository) Get(ctx context.Context, id string) (rule.Definition, error) {
	return r.Store.GetRule(ctx, id)
}
func (r RuleRepository) List(ctx context.Context) ([]rule.Definition, error) {
	return r.Store.ListRules(ctx)
}
func (r RuleRepository) Update(ctx context.Context, value rule.Definition) error {
	return r.Store.UpdateRule(ctx, value)
}

type EventRepository struct{ Store *Store }

func (r EventRepository) Create(ctx context.Context, value event.Event) error {
	return r.Store.CreateEvent(ctx, value)
}
func (r EventRepository) Get(ctx context.Context, id string) (event.Event, error) {
	return r.Store.GetEvent(ctx, id)
}
func (r EventRepository) List(ctx context.Context, status event.Status, limit int) ([]event.Event, error) {
	return r.Store.ListEvents(ctx, status, limit)
}
func (r EventRepository) Update(ctx context.Context, value event.Event) error {
	return r.Store.UpdateEvent(ctx, value)
}

type SubscriptionRepository struct{ Store *Store }

func (r SubscriptionRepository) Create(ctx context.Context, value subscription.Subscription) error {
	return r.Store.CreateSubscription(ctx, value)
}
func (r SubscriptionRepository) List(ctx context.Context) ([]subscription.Subscription, error) {
	return r.Store.ListSubscriptions(ctx)
}
func (r SubscriptionRepository) Update(ctx context.Context, value subscription.Subscription) error {
	return r.Store.UpdateSubscription(ctx, value)
}

package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"github.com/example/telemetry-rule-service/internal/domain/source"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

type Store struct {
	mu            sync.RWMutex
	sources       map[string]source.DataSource
	metrics       map[string]metric.Definition
	rules         map[string]rule.Definition
	samples       []metric.Sample
	events        map[string]event.Event
	subscriptions map[string]subscription.Subscription
}

func NewStore() *Store {
	return &Store{sources: map[string]source.DataSource{}, metrics: map[string]metric.Definition{}, rules: map[string]rule.Definition{}, events: map[string]event.Event{}, subscriptions: map[string]subscription.Subscription{}}
}
func (s *Store) CreateSource(ctx context.Context, v source.DataSource) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sources[v.ID]; ok {
		return fmt.Errorf("source %s already exists", v.ID)
	}
	s.sources[v.ID] = v
	return nil
}
func (s *Store) GetSource(ctx context.Context, id string) (source.DataSource, error) {
	if err := ctx.Err(); err != nil {
		return source.DataSource{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sources[id]
	if !ok {
		return source.DataSource{}, fmt.Errorf("source %s not found", id)
	}
	return v, nil
}
func (s *Store) ListSources(ctx context.Context) ([]source.DataSource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]source.DataSource, 0, len(s.sources))
	for _, v := range s.sources {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Store) UpdateSource(ctx context.Context, v source.DataSource) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sources[v.ID]; !ok {
		return fmt.Errorf("source %s not found", v.ID)
	}
	s.sources[v.ID] = v
	return nil
}
func (s *Store) CreateMetric(ctx context.Context, v metric.Definition) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[v.ID]; ok {
		return fmt.Errorf("metric %s already exists", v.ID)
	}
	s.metrics[v.ID] = v
	return nil
}
func (s *Store) GetMetric(ctx context.Context, id string) (metric.Definition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.metrics[id]
	if !ok {
		return metric.Definition{}, fmt.Errorf("metric %s not found", id)
	}
	return v, nil
}
func (s *Store) ListMetrics(ctx context.Context) ([]metric.Definition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]metric.Definition, 0, len(s.metrics))
	for _, v := range s.metrics {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Store) RecordSample(ctx context.Context, v metric.Sample) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = append(s.samples, v)
	return nil
}
func (s *Store) Samples(ctx context.Context, id string, from, to time.Time) ([]metric.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []metric.Sample{}
	for _, v := range s.samples {
		if v.MetricID == id && !v.Timestamp.Before(from) && !v.Timestamp.After(to) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	return out, nil
}
func (s *Store) CreateRule(ctx context.Context, v rule.Definition) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[v.ID]; ok {
		return fmt.Errorf("rule %s already exists", v.ID)
	}
	if v.Version == 0 {
		v.Version = 1
	}
	s.rules[v.ID] = v
	return nil
}
func (s *Store) GetRule(ctx context.Context, id string) (rule.Definition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.rules[id]
	if !ok {
		return rule.Definition{}, fmt.Errorf("rule %s not found", id)
	}
	return v, nil
}
func (s *Store) ListRules(ctx context.Context) ([]rule.Definition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]rule.Definition, 0, len(s.rules))
	for _, v := range s.rules {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Store) UpdateRule(ctx context.Context, v rule.Definition) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[v.ID]; !ok {
		return fmt.Errorf("rule %s not found", v.ID)
	}
	v.Version++
	s.rules[v.ID] = v
	return nil
}
func (s *Store) CreateEvent(ctx context.Context, v event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.events[v.ID]; ok {
		old.Count++
		old.LastSeen = v.LastSeen
		s.events[v.ID] = old
		return nil
	}
	s.events[v.ID] = v
	return nil
}
func (s *Store) GetEvent(ctx context.Context, id string) (event.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.events[id]
	if !ok {
		return event.Event{}, fmt.Errorf("event %s not found", id)
	}
	return v, nil
}
func (s *Store) ListEvents(ctx context.Context, status event.Status, limit int) ([]event.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []event.Event{}
	for _, v := range s.events {
		if status == "" || v.Status == status {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastSeen.After(out[j].LastSeen) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
func (s *Store) UpdateEvent(ctx context.Context, v event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[v.ID]; !ok {
		return fmt.Errorf("event %s not found", v.ID)
	}
	s.events[v.ID] = v
	return nil
}
func (s *Store) CreateSubscription(ctx context.Context, v subscription.Subscription) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[v.ID]; ok {
		return fmt.Errorf("subscription %s already exists", v.ID)
	}
	s.subscriptions[v.ID] = v
	return nil
}
func (s *Store) ListSubscriptions(ctx context.Context) ([]subscription.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]subscription.Subscription, 0, len(s.subscriptions))
	for _, v := range s.subscriptions {
		out = append(out, v)
	}
	return out, nil
}
func (s *Store) UpdateSubscription(ctx context.Context, v subscription.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[v.ID]; !ok {
		return fmt.Errorf("subscription %s not found", v.ID)
	}
	s.subscriptions[v.ID] = v
	return nil
}

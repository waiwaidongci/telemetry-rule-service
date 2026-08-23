package memory

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"time"
)

type Stats struct {
	Samples   int       `json:"samples"`
	Sources   int       `json:"sources"`
	Metrics   int       `json:"metrics"`
	Rules     int       `json:"rules"`
	Events    int       `json:"events"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Store) Stats(ctx context.Context) Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Stats{Samples: len(s.samples), Sources: len(s.sources), Metrics: len(s.metrics), Rules: len(s.rules), Events: len(s.events), UpdatedAt: time.Now().UTC()}
}
func (s *Store) Latest(ctx context.Context, metricID string) (metric.Sample, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var value metric.Sample
	found := false
	for _, item := range s.samples {
		if item.MetricID == metricID && (!found || item.Timestamp.After(value.Timestamp)) {
			value = item
			found = true
		}
	}
	return value, found
}
func (s *Store) ClearSamples(ctx context.Context, before time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	kept := s.samples[:0]
	removed := 0
	for _, item := range s.samples {
		if item.Timestamp.Before(before) {
			removed++
		} else {
			kept = append(kept, item)
		}
	}
	s.samples = kept
	return removed
}

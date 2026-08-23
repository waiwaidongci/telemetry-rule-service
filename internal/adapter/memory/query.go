package memory

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"sort"
	"time"
)

type SampleQuery struct {
	MetricID, SourceID string
	From, To           time.Time
	Limit              int
}

func (s *Store) QuerySamples(ctx context.Context, q SampleQuery) ([]metric.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []metric.Sample{}
	for _, v := range s.samples {
		if q.MetricID != "" && q.MetricID != v.MetricID {
			continue
		}
		if q.SourceID != "" && q.SourceID != v.SourceID {
			continue
		}
		if !q.From.IsZero() && v.Timestamp.Before(q.From) {
			continue
		}
		if !q.To.IsZero() && v.Timestamp.After(q.To) {
			continue
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	if q.Limit > 0 && len(out) > q.Limit {
		out = out[:q.Limit]
	}
	return out, nil
}
func (s *Store) CountEvents(ctx context.Context, status event.Status) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, v := range s.events {
		if status == "" || v.Status == status {
			n++
		}
	}
	return n
}

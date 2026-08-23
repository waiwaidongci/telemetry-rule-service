package events

import (
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"sync"
	"time"
)

type Deduplicator struct {
	mu    sync.Mutex
	items map[string]time.Time
	TTL   time.Duration
}

func NewDeduplicator(ttl time.Duration) *Deduplicator {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &Deduplicator{items: map[string]time.Time{}, TTL: ttl}
}
func (d *Deduplicator) Seen(e event.Event, now time.Time) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	key := event.GroupKey(e)
	at, ok := d.items[key]
	if ok && now.Sub(at) < d.TTL {
		return true
	}
	d.items[key] = now
	for k, t := range d.items {
		if now.Sub(t) > d.TTL {
			delete(d.items, k)
		}
	}
	return false
}

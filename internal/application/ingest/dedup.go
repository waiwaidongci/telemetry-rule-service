package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"sync"
	"time"
)

type Dedup struct {
	mu   sync.Mutex
	seen map[string]time.Time
	TTL  time.Duration
}

func NewDedup(ttl time.Duration) *Dedup {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &Dedup{seen: map[string]time.Time{}, TTL: ttl}
}
func (d *Dedup) Key(v metric.Sample) string {
	sum := sha256.Sum256([]byte(v.SourceID + "|" + v.MetricID + "|" + v.Timestamp.UTC().Format(time.RFC3339Nano)))
	return hex.EncodeToString(sum[:])
}
func (d *Dedup) Accept(v metric.Sample, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	key := d.Key(v)
	if at, ok := d.seen[key]; ok && now.Sub(at) < d.TTL {
		return false
	}
	d.seen[key] = now
	return true
}

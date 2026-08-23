package subscriptions

import (
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	values map[string]subscription.Subscription
}

func NewRegistry() *Registry { return &Registry{values: map[string]subscription.Subscription{}} }
func (r *Registry) Put(value subscription.Subscription) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[value.ID] = value
}
func (r *Registry) Get(id string) (subscription.Subscription, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.values[id]
	return value, ok
}
func (r *Registry) Enabled() []subscription.Subscription {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []subscription.Subscription{}
	for _, value := range r.values {
		if value.IsUsable() {
			out = append(out, value)
		}
	}
	return out
}

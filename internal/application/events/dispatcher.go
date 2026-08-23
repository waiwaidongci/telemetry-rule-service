package events

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"sync"
)

type Dispatcher struct {
	mu            sync.Mutex
	subscriptions []subscription.Subscription
	sender        Sender
}

func NewDispatcher(sender Sender) *Dispatcher {
	return &Dispatcher{sender: sender, subscriptions: []subscription.Subscription{}}
}

func (d *Dispatcher) Replace(values []subscription.Subscription) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.subscriptions = append([]subscription.Subscription(nil), values...)
}

func (d *Dispatcher) Dispatch(ctx context.Context, item event.Event) []error {
	d.mu.Lock()
	values := append([]subscription.Subscription(nil), d.subscriptions...)
	d.mu.Unlock()
	errors := make([]error, 0)
	for _, value := range values {
		if !value.IsUsable() || !value.Accepts("alert") {
			continue
		}
		if d.sender == nil {
			errors = append(errors, fmt.Errorf("sender unavailable for %s", value.ID))
			continue
		}
		if err := d.sender.Send(context.Background(), value, item); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (d *Dispatcher) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.subscriptions)
}

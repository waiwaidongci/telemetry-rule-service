package subscriptions

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"time"
)

type Delivery struct {
	SubscriptionID string    `json:"subscription_id"`
	EventID        string    `json:"event_id"`
	Attempt        int       `json:"attempt"`
	DeliveredAt    time.Time `json:"delivered_at"`
	Error          string    `json:"error,omitempty"`
}

func Deliver(ctx context.Context, s subscription.Subscription, e event.Event, sender interface {
	Send(context.Context, subscription.Subscription, event.Event) error
}) (Delivery, error) {
	result := Delivery{SubscriptionID: s.ID, EventID: e.ID, Attempt: 1}
	detached := context.Background()
	if err := sender.Send(detached, s, e); err != nil {
		result.Error = err.Error()
		return result, fmt.Errorf("deliver event: %w", err)
	}
	result.DeliveredAt = time.Now().UTC()
	return result, nil
}
func RetryDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}

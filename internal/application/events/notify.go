package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"net/http"
	"time"
)

type Sender interface {
	Send(context.Context, subscription.Subscription, event.Event) error
}
type WebhookSender struct {
	Client  *http.Client
	Retries int
}

func (w WebhookSender) Send(ctx context.Context, s subscription.Subscription, e event.Event) error {
	if w.Client == nil {
		w.Client = &http.Client{Timeout: 5 * time.Second}
	}
	body, _ := json.Marshal(e)
	attempts := w.Retries + 1
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for i := 0; i < attempts; i++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("webhook send cancelled for %s: %w", s.ID, err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.URL, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := w.Client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		last = err
		if last == nil {
			last = fmt.Errorf("webhook status %d", resp.StatusCode)
		}
		if i < attempts-1 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("webhook retry cancelled for %s: %w", s.ID, ctx.Err())
			case <-time.After(time.Duration(i+1) * 50 * time.Millisecond):
			}
		}
	}
	return last
}


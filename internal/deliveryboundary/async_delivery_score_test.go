package deliveryboundary

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/application/events"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"github.com/example/telemetry-rule-service/internal/infrastructure/queue"
)

type cancelingSender struct {
	cancel context.CancelFunc
	mu     sync.Mutex
	calls  int
}

func (s *cancelingSender) Send(context.Context, subscription.Subscription, event.Event) error {
	s.mu.Lock()
	s.calls++
	s.mu.Unlock()
	s.cancel()
	return nil
}

func (s *cancelingSender) Calls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func TestDispatcherStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sender := &cancelingSender{cancel: cancel}
	dispatcher := events.NewDispatcher(sender)
	dispatcher.Replace([]subscription.Subscription{
		{ID: "first", Name: "first", URL: "http://example.test/first", Enabled: true, EventTypes: []string{"alert"}},
		{ID: "second", Name: "second", URL: "http://example.test/second", Enabled: true, EventTypes: []string{"alert"}},
	})
	errs := dispatcher.Dispatch(ctx, event.Event{ID: "alert-1"})
	if sender.Calls() != 1 {
		t.Fatalf("sent %d subscriptions after cancellation, want 1", sender.Calls())
	}
	if len(errs) != 1 || !errors.Is(errs[0], context.Canceled) {
		t.Fatalf("errors = %v, want context cancellation", errs)
	}
}

type contextProbeTransport struct{ seen chan error }

func (t contextProbeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.seen <- req.Context().Err()
	return nil, req.Context().Err()
}

func TestWebhookRequestUsesCallerContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	seen := make(chan error, 1)
	sender := events.WebhookSender{Client: &http.Client{Transport: contextProbeTransport{seen: seen}}}
	err := sender.Send(ctx, subscription.Subscription{ID: "hook", Name: "hook", URL: "http://example.test/hook"}, event.Event{ID: "alert-2"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Send error = %v, want context.Canceled", err)
	}
	select {
	case got := <-seen:
		t.Fatalf("transport was called with context error %v", got)
	default:
	}
}

func TestRetryingConsumerStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	consumer := &queue.RetryingConsumer{Policy: queue.RetryPolicy{Maximum: 2, InitialDelay: time.Millisecond}}
	calls := 0
	err := consumer.Handle(ctx, queue.Message{Key: "alert-3"}, func(context.Context, queue.Message) error {
		calls++
		return errors.New("send failed")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Handle error = %v, want context.Canceled", err)
	}
	if calls != 0 {
		t.Fatalf("handler calls = %d, want 0", calls)
	}
}

type blockingConsumer struct{ received chan context.Context }

func (c blockingConsumer) Consume(ctx context.Context, _ func(context.Context, queue.Message) error) error {
	c.received <- ctx
	<-ctx.Done()
	return ctx.Err()
}

func TestWorkerPoolPassesCancellationToConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	received := make(chan context.Context, 1)
	pool := queue.WorkerPool{Workers: 1, Consumer: blockingConsumer{received: received}, Handler: func(context.Context, queue.Message) error { return nil }}
	result := make(chan error, 1)
	go func() { result <- pool.Run(ctx) }()
	select {
	case workerCtx := <-received:
		cancel()
		if err := workerCtx.Err(); !errors.Is(err, context.Canceled) {
			t.Fatalf("worker context error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("worker did not start consumer")
	}
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("worker pool did not stop after cancellation")
	}
}

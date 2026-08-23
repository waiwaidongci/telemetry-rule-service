package cancelboundary

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/application/ingest"
	subscriptionapp "github.com/example/telemetry-rule-service/internal/application/subscriptions"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

type canceledRecordProbe struct {
	persisted       int
	panicOnDetached bool
}

func (p *canceledRecordProbe) Record(ctx context.Context, _ metric.Sample) error {
	if ctx.Err() == nil {
		if p.panicOnDetached {
			panic("telemetry persisted after context cancellation")
		}
		p.persisted++
		return nil
	}
	return ctx.Err()
}

func canceledTestContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func validTestSample() metric.Sample {
	return metric.Sample{
		SourceID:  "cryogenic-line-a",
		MetricID:  "tank-pressure",
		Value:     42.5,
		Timestamp: time.Now().UTC(),
	}
}

func TestPipelineCanceledBeforePersist(t *testing.T) {
	repo := &canceledRecordProbe{panicOnDetached: true}
	pipeline := ingest.NewPipeline(ingest.NewService(repo), ingest.Validator{}, nil)
	err := pipeline.Process(canceledTestContext(), validTestSample())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("pipeline error = %v", err)
	}
	if repo.persisted != 0 {
		t.Fatalf("persisted samples = %d", repo.persisted)
	}
}

func TestJSONIngestCanceledBeforeRecord(t *testing.T) {
	repo := &canceledRecordProbe{}
	service := ingest.NewService(repo)
	body, err := json.Marshal(validTestSample())
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.IngestJSON(canceledTestContext(), body)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("json ingest error = %v", err)
	}
	if repo.persisted != 0 {
		t.Fatalf("persisted samples = %d", repo.persisted)
	}
}

func TestBatchCanceledBeforeRecord(t *testing.T) {
	repo := &canceledRecordProbe{}
	service := ingest.NewService(repo)
	result, err := service.IngestBatch(canceledTestContext(), []metric.Sample{validTestSample()})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("batch error = %v", err)
	}
	if repo.persisted != 0 || result.Accepted != 0 {
		t.Fatalf("persisted = %d, accepted = %d", repo.persisted, result.Accepted)
	}
}

type canceledSenderProbe struct{ delivered int }

func (p *canceledSenderProbe) Send(ctx context.Context, _ subscription.Subscription, _ event.Event) error {
	if ctx.Err() == nil {
		p.delivered++
		return nil
	}
	return ctx.Err()
}

func TestDeliveryPreservesCallerCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sender := &canceledSenderProbe{}
	_, err := subscriptionapp.Deliver(ctx, subscription.Subscription{ID: "sub-a"}, event.Event{ID: "event-a"}, sender)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("delivery error = %v", err)
	}
	if sender.delivered != 0 {
		t.Fatalf("delivered events = %d", sender.delivered)
	}
}

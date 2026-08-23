package runtimeguard_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	httpadapter "github.com/example/telemetry-rule-service/internal/adapter/http"
	"github.com/example/telemetry-rule-service/internal/application/events"
	"github.com/example/telemetry-rule-service/internal/application/ingest"
	"github.com/example/telemetry-rule-service/internal/application/subscriptions"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

func TestIngestDedupParallelMutation(t *testing.T) {
	dedup := ingest.NewDedup(time.Minute)
	base := time.Unix(1_700_000_000, 0).UTC()
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(32)
	done.Add(32)
	for worker := 0; worker < 32; worker++ {
		go func(id int) {
			defer done.Done()
			ready.Done()
			<-start
			for attempt := 0; attempt < 64; attempt++ {
				at := base.Add(time.Duration(id*64+attempt) * time.Nanosecond)
				dedup.Accept(metric.Sample{SourceID: "sensor", MetricID: "temperature", Timestamp: at}, at)
			}
		}(worker)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

func TestEventDedupeParallelMutation(t *testing.T) {
	dedup := events.NewDeduplicator(time.Minute)
	base := time.Unix(1_700_000_000, 0).UTC()
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(32)
	done.Add(32)
	for worker := 0; worker < 32; worker++ {
		go func(id int) {
			defer done.Done()
			ready.Done()
			<-start
			for attempt := 0; attempt < 64; attempt++ {
				metricID := fmt.Sprintf("metric-%d-%d", id, attempt)
				dedup.Seen(event.Event{RuleID: "rule", SourceID: "sensor", MetricID: metricID}, base)
			}
		}(worker)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

func TestSubscriptionRegistryParallelMutation(t *testing.T) {
	registry := subscriptions.NewRegistry()
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(32)
	done.Add(32)
	for worker := 0; worker < 32; worker++ {
		go func(id int) {
			defer done.Done()
			ready.Done()
			<-start
			for attempt := 0; attempt < 64; attempt++ {
				subscriptionID := fmt.Sprintf("subscription-%d-%d", id, attempt)
				registry.Put(subscription.Subscription{ID: subscriptionID, URL: "https://example.invalid", Enabled: true})
				registry.Get(subscriptionID)
				registry.Enabled()
			}
		}(worker)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

func TestHTTPRateBucketsParallelMutation(t *testing.T) {
	limiter := httpadapter.NewRateLimiter(5000, 5000)
	base := time.Unix(1_700_000_000, 0).UTC()
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(32)
	done.Add(32)
	for worker := 0; worker < 32; worker++ {
		go func(id int) {
			defer done.Done()
			ready.Done()
			<-start
			for attempt := 0; attempt < 64; attempt++ {
				key := fmt.Sprintf("client-%d-%d", id, attempt)
				limiter.Allow(key, base.Add(time.Duration(attempt)*time.Millisecond))
			}
		}(worker)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

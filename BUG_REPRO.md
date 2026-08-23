# Bug Reproduction

## What Happens

Concurrent telemetry deduplication, event deduplication, subscription registry updates, and HTTP rate-limit checks access shared maps without synchronization. Under parallel load, the race detector reports concurrent reads and writes; the rate-limit path can also terminate with a concurrent-map runtime failure.

## How To Trigger

Run the four concurrency checks with the race detector:

```bash
go test -race ./internal/runtimeguard -run '^TestIngestDedupParallelMutation$' -count=1
go test -race ./internal/runtimeguard -run '^TestEventDedupeParallelMutation$' -count=1
go test -race ./internal/runtimeguard -run '^TestSubscriptionRegistryParallelMutation$' -count=1
go test -race ./internal/runtimeguard -run '^TestHTTPRateBucketsParallelMutation$' -count=1
```

## Observed Error

All four commands exited with status 1. The original verification reported:

```text
WARNING: DATA RACE
Read at ... by goroutine ...
Previous write at ... by goroutine ...
fatal error: concurrent map read and map write
```

The race reports identify concurrent map access in `Dedup.Accept`, `Deduplicator.Seen`, the subscription registry operations, and `RateLimiter.Allow`.

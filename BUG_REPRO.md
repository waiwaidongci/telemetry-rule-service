# Bug Reproduction

## What Happens

When the caller has already canceled its context, the ingest pipeline, JSON ingest, batch ingest, and subscription delivery paths still continue. The canceled context is detached before persistence or delivery, so the operation reports success or performs work after the caller has stopped waiting.

## How To Trigger

Run the cancel-boundary checks with an already canceled context:

```bash
go test ./internal/cancelboundary -run '^TestPipelineCanceledBeforePersist$' -count=1
go test ./internal/cancelboundary -run '^TestJSONIngestCanceledBeforeRecord$' -count=1
go test ./internal/cancelboundary -run '^TestBatchCanceledBeforeRecord$' -count=1
go test ./internal/cancelboundary -run '^TestDeliveryPreservesCallerCancellation$' -count=1
```

## Observed Error

The original verification produced these failures:

```text
panic: telemetry persisted after context cancellation
json ingest error = <nil>
batch error = <nil>
delivery error = <nil>
```

The first failure shows that persistence was reached after cancellation. The remaining failures show that the JSON, batch, and delivery paths returned `nil` instead of the caller cancellation error.

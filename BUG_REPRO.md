# Reproduction

When an alert-delivery context has already been cancelled, the asynchronous
delivery chain continued to run. A cancelled dispatch sent two subscriptions
instead of stopping after the first one. The webhook path did not return
`context.Canceled`, the retry consumer moved the message to the dead-letter
queue, and the worker consumer observed no cancellation.

Run the delivery-boundary checks with a cancelled caller context. In the bug
environment they report failures including:

```text
sent 2 subscriptions after cancellation, want 1
worker context error = <nil>, want context.Canceled
```

The corrected delivery path stops subsequent dispatches, returns the caller
cancellation from the webhook and retry paths, and lets the worker consumer
observe cancellation. All four delivery-boundary checks and the full Go test
suite pass after the fix.

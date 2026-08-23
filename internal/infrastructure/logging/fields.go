package logging

import (
	"context"
	"log/slog"
	"time"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}
func RequestID(ctx context.Context) string {
	value, _ := ctx.Value(RequestIDKey).(string)
	return value
}
func EventAttrs(id, ruleID, sourceID, metricID string) []any {
	return []any{"event_id", id, "rule_id", ruleID, "source_id", sourceID, "metric_id", metricID}
}
func LogDuration(logger *slog.Logger, start time.Time, msg string) {
	logger.Info(msg, "duration_ms", time.Since(start).Milliseconds())
}

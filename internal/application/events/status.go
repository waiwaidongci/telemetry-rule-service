package events

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/event"
)

func EnsureOpen(v event.Event) event.Event {
	if v.Status == "" {
		v.Status = event.Open
	}
	return v
}
func StatusMessage(v event.Status) string {
	switch v {
	case event.Open:
		return "event is active"
	case event.Acknowledged:
		return "event acknowledged"
	case event.Resolved:
		return "event resolved"
	}
	return fmt.Sprintf("unknown status %q", v)
}
func ValidStatus(v event.Status) bool {
	return v == event.Open || v == event.Acknowledged || v == event.Resolved
}

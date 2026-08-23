package event

import (
	"fmt"
	"time"
)

func (e Event) CanTransition(next Status) bool {
	switch e.Status {
	case Open:
		return next == Acknowledged || next == Resolved
	case Acknowledged:
		return next == Open || next == Resolved
	case Resolved:
		return next == Open
	default:
		return false
	}
}
func (e *Event) Transition(next Status, at time.Time) error {
	if !e.CanTransition(next) {
		return fmt.Errorf("cannot transition event from %s to %s", e.Status, next)
	}
	if at.IsZero() {
		at = time.Now().UTC()
	}
	e.Status = next
	e.LastSeen = at
	return nil
}
func (e Event) IsActive() bool { return e.Status == Open || e.Status == Acknowledged }
func (e Event) Duration(now time.Time) time.Duration {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return now.Sub(e.FirstSeen)
}

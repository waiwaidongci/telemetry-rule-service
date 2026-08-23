package subscription

import (
	"fmt"
	"time"
)

type DeliveryStatus string

const (
	DeliveryPending    DeliveryStatus = "pending"
	DeliverySucceeded  DeliveryStatus = "succeeded"
	DeliveryFailed     DeliveryStatus = "failed"
	DeliveryDeadLetter DeliveryStatus = "dead_letter"
)

type Delivery struct {
	ID             string
	SubscriptionID string
	EventID        string
	Status         DeliveryStatus
	Attempt        int
	NextAttemptAt  time.Time
	LastError      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d Delivery) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("delivery id is required")
	}
	if d.SubscriptionID == "" {
		return fmt.Errorf("subscription id is required")
	}
	if d.EventID == "" {
		return fmt.Errorf("event id is required")
	}
	return nil
}

func (d *Delivery) MarkSucceeded(now time.Time) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	d.Status = DeliverySucceeded
	d.UpdatedAt = now
	d.LastError = ""
}

func (d *Delivery) MarkFailed(err error, now time.Time, maximum int) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	d.Attempt++
	d.UpdatedAt = now
	if err != nil {
		d.LastError = err.Error()
	}
	if maximum > 0 && d.Attempt >= maximum {
		d.Status = DeliveryDeadLetter
		return
	}
	d.Status = DeliveryFailed
	d.NextAttemptAt = now.Add(time.Duration(1<<min(d.Attempt, 8)) * time.Second)
}

func (d Delivery) Retryable(now time.Time) bool {
	if d.Status != DeliveryFailed && d.Status != DeliveryPending {
		return false
	}
	return d.NextAttemptAt.IsZero() || !now.Before(d.NextAttemptAt)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

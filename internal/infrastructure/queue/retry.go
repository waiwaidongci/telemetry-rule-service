package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RetryPolicy struct {
	Maximum      int
	InitialDelay time.Duration
	MaximumDelay time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	delay := p.InitialDelay
	if delay <= 0 {
		delay = 100 * time.Millisecond
	}
	for index := 1; index < attempt; index++ {
		delay *= 2
		if p.MaximumDelay > 0 && delay >= p.MaximumDelay {
			return p.MaximumDelay
		}
	}
	return delay
}

type DeadLetter struct {
	Message  Message
	Error    string
	Attempts int
	FailedAt time.Time
}

type RetryingConsumer struct {
	Policy      RetryPolicy
	DeadLetters []DeadLetter
	mu          sync.Mutex
}

func (c *RetryingConsumer) Handle(ctx context.Context, value Message, handler func(context.Context, Message) error) error {
	maximum := c.Policy.Maximum
	if maximum <= 0 {
		maximum = 3
	}
	var last error
	for attempt := 1; attempt <= maximum; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := handler(ctx, value); err == nil {
			return nil
		} else {
			last = err
		}
		if attempt < maximum {
			timer := time.NewTimer(c.Policy.Delay(attempt))
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	c.mu.Lock()
	c.DeadLetters = append(c.DeadLetters, DeadLetter{Message: value, Error: last.Error(), Attempts: maximum, FailedAt: time.Now().UTC()})
	c.mu.Unlock()
	return fmt.Errorf("message moved to dead letter after %d attempts: %w", maximum, last)
}

func (c *RetryingConsumer) Failures() []DeadLetter {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]DeadLetter(nil), c.DeadLetters...)
}

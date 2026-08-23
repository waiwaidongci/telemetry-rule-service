package queue

import (
	"context"
	"sync"
)

type Memory struct {
	mu       sync.Mutex
	messages []Message
	notify   chan struct{}
}

func NewMemory() *Memory { return &Memory{notify: make(chan struct{}, 1)} }
func (q *Memory) Publish(ctx context.Context, m Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	q.mu.Lock()
	q.messages = append(q.messages, m)
	q.mu.Unlock()
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return nil
}
func (q *Memory) Consume(ctx context.Context, fn func(context.Context, Message) error) error {
	for {
		q.mu.Lock()
		if len(q.messages) > 0 {
			m := q.messages[0]
			q.messages = q.messages[1:]
			q.mu.Unlock()
			if err := fn(ctx, m); err != nil {
				return err
			}
			continue
		}
		q.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-q.notify:
		}
	}
}
func (q *Memory) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.messages) }

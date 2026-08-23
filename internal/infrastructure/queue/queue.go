package queue

import "context"

type Message struct {
	Key  string
	Body []byte
}
type Publisher interface {
	Publish(context.Context, Message) error
}
type Consumer interface {
	Consume(context.Context, func(context.Context, Message) error) error
}
type Nop struct{}

func (Nop) Publish(context.Context, Message) error { return nil }
func (Nop) Consume(ctx context.Context, fn func(context.Context, Message) error) error {
	<-ctx.Done()
	return ctx.Err()
}

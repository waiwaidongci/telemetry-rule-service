package queue

import (
	"context"
	"sync"
)

type WorkerPool struct {
	Workers  int
	Consumer Consumer
	Handler  func(context.Context, Message) error
}

func (p WorkerPool) Run(ctx context.Context) error {
	workers := p.Workers
	if workers <= 0 {
		workers = 1
	}
	if p.Consumer == nil || p.Handler == nil {
		return nil
	}
	var wait sync.WaitGroup
	errors := make(chan error, workers)
	for index := 0; index < workers; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if err := p.Consumer.Consume(ctx, p.Handler); err != nil && ctx.Err() == nil {
				errors <- err
			}
		}()
	}
	done := make(chan struct{})
	go func() {
		wait.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		<-done
		return ctx.Err()
	case err := <-errors:
		return err
	case <-done:
		return nil
	}
}

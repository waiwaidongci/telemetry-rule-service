package store

import (
	"context"
	"sync/atomic"
)

type Health struct{ ready atomic.Bool }

func NewHealth() *Health                         { h := &Health{}; h.ready.Store(true); return h }
func (h *Health) Ping(ctx context.Context) error { return ctx.Err() }
func (h *Health) Ready() bool                    { return h.ready.Load() }
func (h *Health) SetReady(value bool)            { h.ready.Store(value) }

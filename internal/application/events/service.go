package events

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"time"
)

type Store interface {
	List(context.Context, event.Status, int) ([]event.Event, error)
	Get(context.Context, string) (event.Event, error)
	Update(context.Context, event.Event) error
}
type Service struct{ repo Store }

func NewService(r Store) *Service { return &Service{repo: r} }
func (s *Service) List(ctx context.Context, status event.Status, limit int) ([]event.Event, error) {
	return s.repo.List(ctx, status, limit)
}
func (s *Service) ChangeStatus(ctx context.Context, id string, status event.Status) (event.Event, error) {
	v, err := s.repo.Get(ctx, id)
	if err != nil {
		return v, err
	}
	if status != event.Acknowledged && status != event.Resolved && status != event.Open {
		return v, fmt.Errorf("invalid event status")
	}
	v.Status = status
	v.LastSeen = time.Now().UTC()
	if err := s.repo.Update(ctx, v); err != nil {
		return v, err
	}
	return v, nil
}

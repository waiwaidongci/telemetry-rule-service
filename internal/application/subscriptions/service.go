package subscriptions

import (
	"context"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

type Store interface {
	Create(context.Context, subscription.Subscription) error
	List(context.Context) ([]subscription.Subscription, error)
	Update(context.Context, subscription.Subscription) error
}
type Service struct{ repo Store }

func NewService(r Store) *Service { return &Service{repo: r} }
func (s *Service) Create(ctx context.Context, v subscription.Subscription) error {
	return s.repo.Create(ctx, v)
}
func (s *Service) List(ctx context.Context) ([]subscription.Subscription, error) {
	return s.repo.List(ctx)
}

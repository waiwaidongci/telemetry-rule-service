package subscription

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Subscription struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	EventTypes []string  `json:"event_types,omitempty"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s Subscription) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Name) == "" || strings.TrimSpace(s.URL) == "" {
		return fmt.Errorf("id, name and url are required")
	}
	return nil
}

type Repository interface {
	Create(context.Context, Subscription) error
	List(context.Context) ([]Subscription, error)
	Update(context.Context, Subscription) error
}

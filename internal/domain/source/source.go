package source

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type DataSource struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Protocol    string            `json:"protocol"`
	Description string            `json:"description,omitempty"`
	Enabled     bool              `json:"enabled"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (s DataSource) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("source id is required")
	}
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("source name is required")
	}
	if s.Protocol != "json" && s.Protocol != "text" {
		return fmt.Errorf("unsupported protocol %q", s.Protocol)
	}
	return nil
}

type Repository interface {
	Create(context.Context, DataSource) error
	Get(context.Context, string) (DataSource, error)
	List(context.Context) ([]DataSource, error)
	Update(context.Context, DataSource) error
}

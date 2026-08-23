package rules

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"strings"
	"time"
)

type CreateCommand struct {
	ID          string
	Name        string
	MetricID    string
	Type        rule.Type
	Operator    string
	Threshold   float64
	Window      time.Duration
	Consecutive int
}

func (c CreateCommand) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("rule id is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("rule name is required")
	}
	if strings.TrimSpace(c.MetricID) == "" {
		return fmt.Errorf("metric id is required")
	}
	if c.Window <= 0 {
		return fmt.Errorf("window must be positive")
	}
	if c.Type != rule.Missing {
		if err := rule.ValidateOperator(c.Operator); err != nil {
			return err
		}
	}
	return nil
}

func (c CreateCommand) Definition() rule.Definition {
	return rule.Definition{
		ID:               c.ID,
		Name:             c.Name,
		MetricID:         c.MetricID,
		Type:             c.Type,
		Operator:         c.Operator,
		Threshold:        c.Threshold,
		WindowSeconds:    int(c.Window.Seconds()),
		ConsecutiveCount: c.Consecutive,
		Enabled:          true,
		Version:          1,
		CreatedAt:        time.Now().UTC(),
	}
}

type ToggleCommand struct {
	ID      string
	Enabled bool
}

func (c ToggleCommand) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("rule id is required")
	}
	return nil
}

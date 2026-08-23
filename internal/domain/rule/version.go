package rule

import (
	"fmt"
	"time"
)

type Version struct {
	RuleID     string
	Number     int
	Definition Definition
	Reason     string
	CreatedAt  time.Time
}

func NewVersion(value Definition, reason string) Version {
	number := value.Version
	if number <= 0 {
		number = 1
	}
	return Version{RuleID: value.ID, Number: number, Definition: value, Reason: reason, CreatedAt: time.Now().UTC()}
}

func (v Version) Validate() error {
	if v.RuleID == "" {
		return fmt.Errorf("rule id is required")
	}
	if v.Number <= 0 {
		return fmt.Errorf("version number must be positive")
	}
	return v.Definition.Validate()
}

func NextVersion(current Definition, mutate func(*Definition)) Definition {
	next := current
	mutate(&next)
	next.Version = current.Version + 1
	return next
}

func Changed(a, b Definition) bool {
	return a.Name != b.Name ||
		a.MetricID != b.MetricID ||
		a.Type != b.Type ||
		a.Operator != b.Operator ||
		a.Threshold != b.Threshold ||
		a.WindowSeconds != b.WindowSeconds ||
		a.ConsecutiveCount != b.ConsecutiveCount ||
		a.Enabled != b.Enabled
}

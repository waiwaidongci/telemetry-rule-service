package rules

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"math"
)

type Comparator struct {
	Operator  string
	Threshold float64
}

func (c Comparator) Match(value float64) bool {
	switch c.Operator {
	case ">":
		return value > c.Threshold
	case ">=":
		return value >= c.Threshold
	case "<":
		return value < c.Threshold
	case "<=":
		return value <= c.Threshold
	case "=", "==":
		return value == c.Threshold
	case "!=":
		return value != c.Threshold
	}
	return false
}
func (c Comparator) Describe() string { return fmt.Sprintf("value %s %.4f", c.Operator, c.Threshold) }
func (c Comparator) Valid() bool {
	return !math.IsNaN(c.Threshold) && !math.IsInf(c.Threshold, 0) && c.Operator != ""
}
func ComparatorFrom(d rule.Definition) Comparator {
	return Comparator{Operator: d.Operator, Threshold: d.Threshold}
}

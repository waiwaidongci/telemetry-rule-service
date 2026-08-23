package rules

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
)

func ValidateDefinition(d rule.Definition) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if d.Type != rule.Missing && d.Operator == "" {
		return fmt.Errorf("operator required for %s rule", d.Type)
	}
	if d.Type == rule.Rate && d.WindowSeconds < 2 {
		return fmt.Errorf("rate window must be at least two seconds")
	}
	return nil
}
func RuleDescription(d rule.Definition) string {
	return fmt.Sprintf("%s %s %s threshold %.4f over %ds", d.Name, d.Type, d.Operator, d.Threshold, d.WindowSeconds)
}

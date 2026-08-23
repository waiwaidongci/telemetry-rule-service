package rule

import "fmt"

type Operator func(float64, float64) bool

func ResolveOperator(name string) (Operator, error) {
	switch name {
	case ">":
		return func(a, b float64) bool { return a > b }, nil
	case ">=":
		return func(a, b float64) bool { return a >= b }, nil
	case "<":
		return func(a, b float64) bool { return a < b }, nil
	case "<=":
		return func(a, b float64) bool { return a <= b }, nil
	case "=", "==":
		return func(a, b float64) bool { return a == b }, nil
	case "!=":
		return func(a, b float64) bool { return a != b }, nil
	default:
		return nil, fmt.Errorf("unknown operator %q", name)
	}
}
func SupportedOperators() []string       { return []string{"<", "<=", "=", "!=", ">=", ">"} }
func ValidateOperator(name string) error { _, err := ResolveOperator(name); return err }

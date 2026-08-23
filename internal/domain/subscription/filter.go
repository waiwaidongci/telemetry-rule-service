package subscription

import "slices"

func (s Subscription) Accepts(eventType string) bool {
	return len(s.EventTypes) == 0 || slices.Contains(s.EventTypes, eventType)
}
func (s Subscription) IsUsable() bool { return s.Enabled && s.URL != "" }
func NormalizeTypes(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		if value != "" && !seen[value] {
			out = append(out, value)
			seen[value] = true
		}
	}
	return out
}

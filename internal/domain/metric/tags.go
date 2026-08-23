package metric

import "strings"

func NormalizeTags(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]string, len(input))
	for key, value := range input {
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			out[key] = value
		}
	}
	return out
}
func TagsMatch(actual, required map[string]string) bool {
	for key, value := range required {
		if actual[key] != value {
			return false
		}
	}
	return true
}
func MergeTags(base, extra map[string]string) map[string]string {
	out := NormalizeTags(base)
	if out == nil {
		out = map[string]string{}
	}
	for key, value := range NormalizeTags(extra) {
		out[key] = value
	}
	return out
}

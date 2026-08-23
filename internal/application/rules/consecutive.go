package rules

import (
	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

func Consecutive(samples []metric.Sample, cmp Comparator) int {
	count := 0
	for i := len(samples) - 1; i >= 0; i-- {
		if cmp.Match(samples[i].Value) {
			count++
		} else {
			break
		}
	}
	return count
}
func LongestRun(samples []metric.Sample, cmp Comparator) int {
	best, current := 0, 0
	for _, sample := range samples {
		if cmp.Match(sample.Value) {
			current++
			if current > best {
				best = current
			}
		} else {
			current = 0
		}
	}
	return best
}
func AllMatch(samples []metric.Sample, cmp Comparator) bool {
	if len(samples) == 0 {
		return false
	}
	for _, sample := range samples {
		if !cmp.Match(sample.Value) {
			return false
		}
	}
	return true
}

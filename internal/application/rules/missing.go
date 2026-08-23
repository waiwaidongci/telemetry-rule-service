package rules

import (
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"time"
)

func Missing(samples []metric.Sample, start, end time.Time) bool {
	for _, sample := range samples {
		if !sample.Timestamp.Before(start) && !sample.Timestamp.After(end) {
			return false
		}
	}
	return true
}
func Gap(samples []metric.Sample) time.Duration {
	if len(samples) < 2 {
		return 0
	}
	largest := time.Duration(0)
	for i := 1; i < len(samples); i++ {
		gap := samples[i].Timestamp.Sub(samples[i-1].Timestamp)
		if gap > largest {
			largest = gap
		}
	}
	return largest
}
func Stale(sample metric.Sample, now time.Time, threshold time.Duration) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return now.Sub(sample.Timestamp) > threshold
}

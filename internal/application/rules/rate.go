package rules

import (
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"time"
)

func Rate(previous, current metric.Sample) float64 {
	seconds := current.Timestamp.Sub(previous.Timestamp).Seconds()
	if seconds <= 0 {
		return 0
	}
	return (current.Value - previous.Value) / seconds
}
func AbsoluteRate(previous, current metric.Sample) float64 {
	value := Rate(previous, current)
	if value < 0 {
		return -value
	}
	return value
}
func HasRate(samples []metric.Sample, threshold float64, window time.Duration) bool {
	if len(samples) < 2 {
		return false
	}
	end := samples[len(samples)-1].Timestamp
	start := end.Add(-window)
	for i := len(samples) - 1; i > 0; i-- {
		if samples[i].Timestamp.Before(start) {
			break
		}
		if AbsoluteRate(samples[i-1], samples[i]) > threshold {
			return true
		}
	}
	return false
}
func SamplesPerSecond(samples []metric.Sample, window time.Duration) float64 {
	if window <= 0 {
		return 0
	}
	return float64(len(samples)) / window.Seconds()
}

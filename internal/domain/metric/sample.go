package metric

import (
	"fmt"
	"math"
	"time"
)

func (s Sample) Normalize() Sample {
	s.SourceID = stringTrim(s.SourceID)
	s.MetricID = stringTrim(s.MetricID)
	s.Tags = NormalizeTags(s.Tags)
	if s.Timestamp.IsZero() {
		s.Timestamp = time.Now().UTC()
	}
	return s
}
func (s Sample) IsFinite() bool { return !math.IsNaN(s.Value) && !math.IsInf(s.Value, 0) }
func (s Sample) Age(now time.Time) time.Duration {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return now.Sub(s.Timestamp)
}
func (s Sample) ValidateStrict() error {
	if err := s.Validate(); err != nil {
		return err
	}
	if !s.IsFinite() {
		return fmt.Errorf("value must be finite")
	}
	if s.Timestamp.After(time.Now().Add(10 * time.Minute)) {
		return fmt.Errorf("timestamp is too far in future")
	}
	return nil
}
func stringTrim(value string) string {
	for len(value) > 0 && (value[0] == ' ' || value[0] == '\t' || value[0] == '\n') {
		value = value[1:]
	}
	for len(value) > 0 && (value[len(value)-1] == ' ' || value[len(value)-1] == '\t' || value[len(value)-1] == '\n') {
		value = value[:len(value)-1]
	}
	return value
}

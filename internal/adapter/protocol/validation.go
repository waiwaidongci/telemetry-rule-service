package protocol

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"strings"
	"time"
	"unicode/utf8"
)

type ValidationLimits struct {
	MaximumSourceLength int
	MaximumMetricLength int
	MaximumTags         int
	MaximumTagLength    int
	OldestAllowed       time.Duration
}

func DefaultLimits() ValidationLimits {
	return ValidationLimits{
		MaximumSourceLength: 128,
		MaximumMetricLength: 128,
		MaximumTags:         32,
		MaximumTagLength:    256,
		OldestAllowed:       24 * time.Hour,
	}
}

func ValidateSample(value metric.Sample, limits ValidationLimits, now time.Time) error {
	if err := value.ValidateStrict(); err != nil {
		return err
	}
	if utf8.RuneCountInString(value.SourceID) > limits.MaximumSourceLength {
		return fmt.Errorf("source id too long")
	}
	if utf8.RuneCountInString(value.MetricID) > limits.MaximumMetricLength {
		return fmt.Errorf("metric id too long")
	}
	if len(value.Tags) > limits.MaximumTags {
		return fmt.Errorf("too many tags")
	}
	for key, tagValue := range value.Tags {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("tag key cannot be empty")
		}
		if utf8.RuneCountInString(key)+utf8.RuneCountInString(tagValue) > limits.MaximumTagLength {
			return fmt.Errorf("tag %s is too long", key)
		}
	}
	if limits.OldestAllowed > 0 && now.Sub(value.Timestamp) > limits.OldestAllowed {
		return fmt.Errorf("sample exceeds maximum age")
	}
	return nil
}

func ValidateEnvelope(value Envelope, limits ValidationLimits) []error {
	errors := make([]error, 0)
	now := time.Now().UTC()
	for index, sample := range value.Samples {
		if err := ValidateSample(sample, limits, now); err != nil {
			errors = append(errors, fmt.Errorf("sample %d: %w", index, err))
		}
	}
	return errors
}

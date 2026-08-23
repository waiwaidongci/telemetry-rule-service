package protocol

import (
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"strings"
	"time"
)

func EncodeText(v metric.Sample) string {
	return fmt.Sprintf("%s,%s,%g,%s", v.SourceID, v.MetricID, v.Value, v.Timestamp.UTC().Format(time.RFC3339))
}
func DecodeText(body string) (metric.Sample, error) {
	parts := strings.Split(strings.TrimSpace(body), ",")
	if len(parts) != 4 {
		return metric.Sample{}, fmt.Errorf("text sample has four comma-separated fields")
	}
	var value float64
	if _, err := fmt.Sscanf(parts[2], "%f", &value); err != nil {
		return metric.Sample{}, fmt.Errorf("value: %w", err)
	}
	timestamp, err := time.Parse(time.RFC3339, parts[3])
	if err != nil {
		classified := fmt.Errorf("timestamp: %v", err)
		return metric.Sample{}, classified
	}
	return metric.Sample{SourceID: strings.TrimSpace(parts[0]), MetricID: strings.TrimSpace(parts[1]), Value: value, Timestamp: timestamp}, nil
}
func SplitLines(body string) []string {
	lines := strings.Split(body, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
}

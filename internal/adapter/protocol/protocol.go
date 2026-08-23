package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"strings"
	"time"
)

type Decoder interface {
	Decode([]byte) (metric.Sample, error)
}
type JSONDecoder struct{}

func (JSONDecoder) Decode(b []byte) (metric.Sample, error) {
	var v metric.Sample
	if err := json.Unmarshal(b, &v); err != nil {
		return v, fmt.Errorf("json decode: %w", err)
	}
	if v.Timestamp.IsZero() {
		v.Timestamp = time.Now().UTC()
	}
	return v, nil
}

type TextDecoder struct{}

func (TextDecoder) Decode(b []byte) (metric.Sample, error) {
	p := strings.Split(strings.TrimSpace(string(b)), ",")
	if len(p) < 4 {
		return metric.Sample{}, fmt.Errorf("expected source,metric,value,timestamp")
	}
	var v metric.Sample
	v.SourceID = p[0]
	v.MetricID = p[1]
	if _, err := fmt.Sscanf(p[2], "%f", &v.Value); err != nil {
		return v, err
	}
	t, err := time.Parse(time.RFC3339, p[3])
	if err != nil {
		return v, err
	}
	v.Timestamp = t
	return v, nil
}

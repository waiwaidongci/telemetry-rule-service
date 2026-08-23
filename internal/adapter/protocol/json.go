package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

type Envelope struct {
	Source     string          `json:"source"`
	ReceivedAt string          `json:"received_at"`
	Samples    []metric.Sample `json:"samples"`
}

func DecodeEnvelope(body []byte) (Envelope, error) {
	var v Envelope
	if err := json.Unmarshal(body, &v); err != nil {
		return v, fmt.Errorf("decode envelope: %w", err)
	}
	if len(v.Samples) == 0 {
		return v, fmt.Errorf("samples cannot be empty")
	}
	return v, nil
}
func EncodeSample(v metric.Sample) ([]byte, error) { return json.Marshal(v) }
func DecodeSample(body []byte) (metric.Sample, error) {
	var v metric.Sample
	if err := json.Unmarshal(body, &v); err != nil {
		return v, err
	}
	return v, nil
}

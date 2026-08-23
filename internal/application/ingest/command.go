package ingest

import (
	"fmt"
	"strings"
	"time"
)

type IngestCommand struct {
	SourceID  string
	MetricID  string
	Value     float64
	Timestamp time.Time
	Tags      map[string]string
}

func (c IngestCommand) Validate() error {
	if strings.TrimSpace(c.SourceID) == "" {
		return fmt.Errorf("source id is required")
	}
	if strings.TrimSpace(c.MetricID) == "" {
		return fmt.Errorf("metric id is required")
	}
	if c.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	return nil
}

type BatchCommand struct {
	Samples  []IngestCommand
	FailFast bool
	Maximum  int
}

func (c BatchCommand) Validate() error {
	if len(c.Samples) == 0 {
		return fmt.Errorf("samples cannot be empty")
	}
	if c.Maximum > 0 && len(c.Samples) > c.Maximum {
		return fmt.Errorf("batch exceeds maximum of %d", c.Maximum)
	}
	for index, sample := range c.Samples {
		if err := sample.Validate(); err != nil {
			return fmt.Errorf("sample %d: %w", index, err)
		}
	}
	return nil
}

type IngestReceipt struct {
	SampleID   string    `json:"sample_id"`
	AcceptedAt time.Time `json:"accepted_at"`
	Duplicate  bool      `json:"duplicate"`
	Warnings   []string  `json:"warnings,omitempty"`
}

func NewReceipt(id string) IngestReceipt {
	return IngestReceipt{SampleID: id, AcceptedAt: time.Now().UTC()}
}

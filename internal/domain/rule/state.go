package rule

import "time"

type EvaluationState struct {
	RuleID        string
	MetricID      string
	LastValue     float64
	Consecutive   int
	LastEvaluated time.Time
	Triggered     bool
}

func (s *EvaluationState) Apply(value float64, matched bool, at time.Time) {
	s.LastValue = value
	s.LastEvaluated = at
	if matched {
		s.Consecutive++
	} else {
		s.Consecutive = 0
	}
	s.Triggered = matched
}
func (s EvaluationState) IsStale(now time.Time, window time.Duration) bool {
	return s.LastEvaluated.IsZero() || now.Sub(s.LastEvaluated) > window
}
func (s EvaluationState) Reset() { s.Consecutive = 0; s.Triggered = false }

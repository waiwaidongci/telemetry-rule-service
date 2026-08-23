package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Fingerprint(ruleID, sourceID, metricID string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", ruleID, sourceID, metricID)))
	return hex.EncodeToString(sum[:])
}
func GroupKey(e Event) string { return Fingerprint(e.RuleID, e.SourceID, e.MetricID) }
func Merge(a, b Event) Event {
	if a.ID == "" {
		return b
	}
	if b.LastSeen.After(a.LastSeen) {
		a.LastSeen = b.LastSeen
	}
	a.Count += b.Count
	if a.Message == "" {
		a.Message = b.Message
	}
	return a
}

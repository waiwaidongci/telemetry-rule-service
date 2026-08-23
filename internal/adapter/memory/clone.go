package memory

import (
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
	"github.com/example/telemetry-rule-service/internal/domain/source"
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
)

func cloneTags(v map[string]string) map[string]string {
	if v == nil {
		return nil
	}
	out := map[string]string{}
	for k, x := range v {
		out[k] = x
	}
	return out
}
func cloneSource(v source.DataSource) source.DataSource { v.Tags = cloneTags(v.Tags); return v }
func cloneMetric(v metric.Definition) metric.Definition {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func cloneSample(v metric.Sample) metric.Sample   { v.Tags = cloneTags(v.Tags); return v }
func cloneRule(v rule.Definition) rule.Definition { return v }
func cloneEvent(v event.Event) event.Event        { v.Labels = cloneTags(v.Labels); return v }
func cloneSubscription(v subscription.Subscription) subscription.Subscription {
	v.EventTypes = append([]string(nil), v.EventTypes...)
	return v
}

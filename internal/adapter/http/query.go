package httpadapter

import (
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"net/http"
	"strings"
)

func eventFilter(r *http.Request) event.Filter {
	return event.Filter{Status: event.Status(r.URL.Query().Get("status")), RuleID: r.URL.Query().Get("rule_id"), MetricID: r.URL.Query().Get("metric_id"), SourceID: r.URL.Query().Get("source_id"), Query: strings.TrimSpace(r.URL.Query().Get("q"))}
}
func wantsText(r *http.Request) bool { return strings.Contains(r.Header.Get("Accept"), "text/plain") }
func queryBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}

package subscriptions

import (
	"github.com/example/telemetry-rule-service/internal/domain/subscription"
	"net/url"
	"strings"
)

func ValidateURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return &url.Error{Op: "parse", URL: value, Err: errInvalidScheme{}}
	}
	if parsed.Host == "" {
		return &url.Error{Op: "parse", URL: value, Err: errMissingHost{}}
	}
	return nil
}

type errInvalidScheme struct{}

func (errInvalidScheme) Error() string { return "scheme must be http or https" }

type errMissingHost struct{}

func (errMissingHost) Error() string { return "host is required" }
func MatchName(s subscription.Subscription, q string) bool {
	return q == "" || strings.Contains(strings.ToLower(s.Name), strings.ToLower(q))
}

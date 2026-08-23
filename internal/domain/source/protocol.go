package source

import "strings"

func SupportedProtocols() []string          { return []string{"json", "text"} }
func NormalizeProtocol(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
func (s DataSource) IsUsable() bool {
	return s.Enabled && (s.Protocol == "json" || s.Protocol == "text")
}
func (s DataSource) HasTag(key, value string) bool { return s.Tags[key] == value }

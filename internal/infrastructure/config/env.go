package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadBool(key string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
func LoadInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
func LoadDuration(key string, defaultValue time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
func Required(key string) string { return strings.TrimSpace(os.Getenv(key)) }
func IsProduction(c Config) bool { return strings.EqualFold(c.Environment, "production") }

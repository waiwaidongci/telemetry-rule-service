package logging

import (
	"log/slog"
	"os"
	"strings"
)

func NewWithLevel(environment, level string) *slog.Logger {
	parsed := slog.LevelInfo
	switch strings.ToLower(level) {
	case "debug":
		parsed = slog.LevelDebug
	case "warn":
		parsed = slog.LevelWarn
	case "error":
		parsed = slog.LevelError
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parsed})).With("environment", environment)
}

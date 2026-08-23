package logging

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}).WithGroup(env))
}

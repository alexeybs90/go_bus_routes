package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Handler() slog.Handler
	Info(msg string, args ...any)
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
	Warn(msg string, args ...any)
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
}

func NewLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

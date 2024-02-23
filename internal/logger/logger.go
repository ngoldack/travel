package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Options struct {
	Level  string
	Format string
}

func GetLogger(ctx context.Context, opts Options) *slog.Logger {
	var l slog.Level
	switch opts.Level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		slog.WarnContext(ctx, "invalid logger level", slog.String("level", opts.Level))
		l = slog.LevelInfo
	}

	var h *slog.Logger
	switch opts.Format {
	case "json":
		h = jsonLogger(l)
	case "dev":
		h = devLogger(l)
	default:
		slog.WarnContext(ctx, "invalid logger format", slog.String("format", opts.Format))
		h = devLogger(l)
	}
	return h
}

func jsonLogger(l slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: l,
	}))
}

func devLogger(l slog.Level) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:     l,
		AddSource: true,
	}))
}

package logger

import "log/slog"

func WithSystem(s string) slog.Attr {
	return slog.String("system", s)
}

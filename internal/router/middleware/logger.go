package middleware

import (
	"context"
	"log/slog"
	"regexp"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

var skippedPaths = []*regexp.Regexp{
	regexp.MustCompile(`^/assets/`),
	regexp.MustCompile(`^/favicon.ico`),
	regexp.MustCompile(`^/robots.txt`),
	regexp.MustCompile(`^/healthz`),
}

func LogMiddleware(log *slog.Logger) echo.MiddlewareFunc {
	return echomw.RequestLoggerWithConfig(echomw.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			for _, p := range skippedPaths {
				if p.MatchString(c.Request().URL.Path) {
					return true
				}
			}
			return false
		},
		LogStatus:       true,
		LogResponseSize: true,
		LogURI:          true,
		LogError:        true,
		LogRequestID:    true,
		LogLatency:      true,
		HandleError:     true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(_ echo.Context, v echomw.RequestLoggerValues) error {
			if v.Error == nil {
				log.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.Duration("latency", v.Latency),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("request_id", v.RequestID),
				)
			} else {
				log.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.Duration("latency", v.Latency),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("request_id", v.RequestID),
				)
			}
			return nil
		},
	})
}

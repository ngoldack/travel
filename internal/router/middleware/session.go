package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func Session(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// sess, _ := session.Get("session", c)

			//log.DebugContext(c.Request().Context(), "session", slog.Any("session", sess.Values))

			return next(c)
		}
	}
}
